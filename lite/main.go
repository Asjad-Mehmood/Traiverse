// Traverse & Level - lightweight Windows shell.
// Uses the WebView2 runtime that ships with Windows 10 and 11, so the program itself stays
// small. The window opens the web app at www.exceltokml.com/traverse/app, so the user signs in
// with their Excel to KML account and every job is saved to that account. After the first
// visit the web app keeps its files on this computer, so the window also opens offline.
package main

import (
	_ "embed"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/jchv/go-webview2"
)

const (
	appName = "Traverse & Level"
	appURL  = "https://www.exceltokml.com/traverse/app"
	// ?app=1 opens it as an app: no front page, straight to log in (or the dashboard when
	// already signed in)
	startURL = appURL + "?app=1"
)

// signInCode is the one-use code the Survey Toolkit passed (--signin-code=…), or "".
func signInCode(args []string) string {
	ok := regexp.MustCompile(`^[0-9a-f]{64}$`)
	for _, a := range args {
		if c, found := strings.CutPrefix(a, "--signin-code="); found && ok.MatchString(c) {
			return c
		}
	}
	return ""
}

// firstURL is where the window opens: opened from the Survey Toolkit, through the website's
// sign-in with the Toolkit's code (it goes on to the app, signed in to the same account).
func firstURL(code string) string {
	if code == "" {
		return startURL
	}
	return "https://www.exceltokml.com/api/toolkit/app-code/redeem?app=traverse&code=" + url.QueryEscape(code)
}

// shown only when the web app has never been opened on this computer and there is no internet
//
//go:embed offline.html
var offlinePage string

func warn(msg string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	mb := user32.NewProc("MessageBoxW")
	t, _ := syscall.UTF16PtrFromString(appName)
	m, _ := syscall.UTF16PtrFromString(msg)
	mb.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}

// reachable reports whether the web app answers within a few seconds.
func reachable() bool {
	c := &http.Client{Timeout: 8 * time.Second}
	resp, err := c.Head(appURL)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

func main() {
	// a fixed data folder: the sign-in, the offline copy of the app and the jobs waiting to
	// upload are kept here, so they survive restarts and updates
	dataPath := filepath.Join(os.Getenv("AppData"), "TraverseLevel")
	marker := filepath.Join(dataPath, "opened-online")

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		DataPath:  dataPath,
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Professional Traverse and Level Adjustment Software",
			Width:  1440,
			Height: 920,
			IconId: 2,
			Center: true,
		},
	})
	if w == nil {
		warn("This computer does not have the Microsoft Edge WebView2 runtime.\n\n" +
			"Install the WebView2 Runtime free from Microsoft, then open Traverse & Level again.")
		return
	}
	defer w.Destroy()

	go checkForUpdate(w)

	// Once the app has been opened online, its own offline copy takes over when there is no
	// internet, so go straight to it. The very first time, it needs a connection.
	if _, err := os.Stat(marker); err == nil || reachable() {
		os.MkdirAll(dataPath, 0o755)
		os.WriteFile(marker, []byte(time.Now().Format(time.RFC3339)), 0o644)
		// the Toolkit's code needs the website; offline, the app's own copy opens as before
		code := signInCode(os.Args[1:])
		if code != "" && !reachable() {
			code = ""
		}
		w.Navigate(firstURL(code))
	} else {
		w.SetHtml(strings.ReplaceAll(offlinePage, "{{APP_URL}}", startURL))
	}
	w.Run()
}
