//go:build windows

package main

// Self-update for Traverse & Level.
// On start-up the program reads latest.json from your newest GitHub release.
// If the version there is higher, it asks the user, downloads the installer,
// checks its SHA-256, runs it silently and the installer reopens the program.
// Any network problem is ignored: the program simply starts as normal.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/jchv/go-webview2"
)

const (
	// The releases of this repository. It is public, so its releases are too. Every installed
	// copy reads this address, so keep it working: if the code ever moves to a private
	// repository, publish the releases to a public one and keep this name (see README).
	updateBase = "https://github.com/Asjad-Mehmood/Traiverse/releases/latest/download/"
	// Traverse Pro 1.x copies look for TraversePro-Setup.exe; each release carries the same
	// installer under that name too, so they update into Traverse & Level.
	installerName = "TraverseLevel-Setup.exe"
	manifestName  = "latest.json"
)

// Overwritten at build time by the workflow (-X main.appVersion=...).
var appVersion = "0.0.0"

type updateInfo struct {
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
}

func askYesNo(msg string) bool {
	mb := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW")
	t, _ := syscall.UTF16PtrFromString(appName)
	m, _ := syscall.UTF16PtrFromString(msg)
	// 0x24 = Yes/No + question icon, 0x40000 = stay on top. 6 = Yes.
	r, _, _ := mb.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x40024)
	return r == 6
}

func parseVer(v string) [3]int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	var out [3]int
	for i, p := range strings.SplitN(v, ".", 3) {
		if j := strings.IndexAny(p, "-+"); j >= 0 {
			p = p[:j]
		}
		n, _ := strconv.Atoi(p)
		out[i] = n
	}
	return out
}

func isNewer(remote, local string) bool {
	r, l := parseVer(remote), parseVer(local)
	for i := 0; i < 3; i++ {
		if r[i] != l[i] {
			return r[i] > l[i]
		}
	}
	return false
}

func fetchUpdateInfo() (*updateInfo, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Get(updateBase + manifestName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("update manifest: " + resp.Status)
	}
	var u updateInfo
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&u); err != nil {
		return nil, err
	}
	if u.Version == "" || len(u.SHA256) != 64 {
		return nil, errors.New("bad update manifest")
	}
	return &u, nil
}

func downloadInstaller(dest, wantSum string) error {
	c := &http.Client{Timeout: 10 * time.Minute}
	resp, err := c.Get(updateBase + installerName)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("download: " + resp.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	h := sha256.New()
	_, err = io.Copy(io.MultiWriter(f, h), resp.Body)
	cerr := f.Close()
	if err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(dest)
		return err
	}
	if !strings.EqualFold(hex.EncodeToString(h.Sum(nil)), wantSum) {
		os.Remove(dest)
		return errors.New("the downloaded file failed its safety check")
	}
	return nil
}

func checkForUpdate(w webview2.WebView) {
	time.Sleep(4 * time.Second) // let the window appear first
	u, err := fetchUpdateInfo()
	if err != nil || !isNewer(u.Version, appVersion) {
		return
	}
	if !askYesNo("A new version of Traverse & Level is available.\n\n" +
		"Installed: " + appVersion + "\nNew: " + u.Version + "\n\n" +
		"Update now? The program will close, update itself and reopen in a few seconds.") {
		return
	}

	// Notify the Webview UI to display the downloading progress screen immediately
	if w != nil {
		w.Dispatch(func() {
			w.Eval("window.__showUpdateProgress && window.__showUpdateProgress('downloading', '" + u.Version + "');")
		})
	}

	dest := filepath.Join(os.TempDir(), "TraverseLevel-Update.exe")
	if err := downloadInstaller(dest, u.SHA256); err != nil {
		if w != nil {
			w.Dispatch(func() {
				w.Eval("window.__showUpdateProgress && window.__showUpdateProgress('failed', 'Download could not be completed.');")
			})
		}
		warn("The update could not be downloaded. Nothing was changed.\n\n" + err.Error())
		return
	}

	// Notify the Webview UI that download completed and installation/restart is beginning
	if w != nil {
		w.Dispatch(func() {
			w.Eval("window.__showUpdateProgress && window.__showUpdateProgress('installing', '" + u.Version + "');")
		})
	}
	time.Sleep(600 * time.Millisecond)

	cmd := exec.Command(dest, "/S")
	// DETACHED_PROCESS | CREATE_NEW_PROCESS_GROUP so the installer outlives us
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000008 | 0x00000200}
	if err := cmd.Start(); err != nil {
		if w != nil {
			w.Dispatch(func() {
				w.Eval("window.__showUpdateProgress && window.__showUpdateProgress('failed', 'Update installer could not start.');")
			})
		}
		warn("The update could not be started.\n\n" + err.Error())
		return
	}
	os.Exit(0)
}
