// Package update holds the parts of the self-updater that do not need Windows: reading
// latest.json, comparing versions and downloading the installer with its SHA-256 check.
// The program (update_windows.go) asks the user and runs the installer.
package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Info is latest.json: the newest version and its installer's SHA-256.
type Info struct {
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
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

// IsNewer reports whether version remote is higher than version local.
func IsNewer(remote, local string) bool {
	r, l := parseVer(remote), parseVer(local)
	for i := 0; i < 3; i++ {
		if r[i] != l[i] {
			return r[i] > l[i]
		}
	}
	return false
}

// Fetch reads latest.json from url.
func Fetch(url string) (*Info, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("update manifest: " + resp.Status)
	}
	var u Info
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&u); err != nil {
		return nil, err
	}
	if u.Version == "" || len(u.SHA256) != 64 {
		return nil, errors.New("bad update manifest")
	}
	return &u, nil
}

// Download saves url to dest and keeps it only if its SHA-256 is wantSum.
func Download(url, dest, wantSum string) error {
	c := &http.Client{Timeout: 10 * time.Minute}
	resp, err := c.Get(url)
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
