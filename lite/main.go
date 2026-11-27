// Traverse & Level - lightweight Windows shell.
// Uses the WebView2 runtime that ships with Windows 10 and 11, so the program itself stays
// small. The window opens the web app at www.exceltokml.com/traverse/app (the website, now
// www.surveyxdigital.com, under its old address), so the user signs in
// with their SurveyX Digital account and every job is saved to that account. After the first
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
	// Still the old address on purpose: the window keeps the sign-in, the offline copy of the
	// app and the jobs waiting to upload per address, so moving it to www.surveyxdigital.com
	// would leave them behind. exceltokml.com keeps answering for these paths (docs: the moved
	// site forwards /traverse/app and /api). A later version moves it, with the data.
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

// whereToOpen finds where the window opens, which takes a moment (the internet, the Toolkit's
// sign-in), and gives what to do with the window then
func whereToOpen(dataPath, marker string) func(webview2.WebView) {
	// Once the app has been opened online, its own offline copy takes over when there is no
	// internet, so go straight to it. The very first time, it needs a connection.
	if _, err := os.Stat(marker); err == nil || reachable() {
		os.MkdirAll(dataPath, 0o755)
		os.WriteFile(marker, []byte(time.Now().Format(time.RFC3339)), 0o644)
		// the Toolkit's code needs the website; offline, the app's own copy opens as before
		code := signInCode(os.Args[1:])
		// opened from a shortcut: the Survey Toolkit on this PC, if it is signed in, gives a code for
		// its account, once (the website leaves a sign-in the window already has as it is)
		toolkitMarker := filepath.Join(dataPath, "signed-in-by-toolkit")
		if _, err := os.Stat(toolkitMarker); code == "" && err != nil {
			code = askToolkitForCode("traverse")
		}
		if code != "" && !reachable() {
			code = ""
		}
		if code != "" {
			os.WriteFile(toolkitMarker, []byte(time.Now().Format(time.RFC3339)), 0o644)
		}
		first := firstURL(code)
		return func(w webview2.WebView) { w.Navigate(first) }
	}
	page := strings.ReplaceAll(offlinePage, "{{APP_URL}}", startURL)
	return func(w webview2.WebView) { w.SetHtml(page) }
}

func main() {
	defer stopOnPanic()
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
			// the program's icon in the title bar: IDI_APPLICATION (32512), the icon group goversioninfo
			// (pinned in the release workflow) writes from icon.ico. Not 2: that is one of its images.
			IconId: 32512,
			Center: true,
		},
	})
	if w == nil {
		warn("This computer does not have the Microsoft Edge WebView2 runtime.\n\n" +
			"Install the WebView2 Runtime free from Microsoft, then open Traverse & Level again.")
		return
	}
	defer w.Destroy()

	// Help in the app: the error log, the open-source licences, and pages in the browser;
	// the page's own errors go into the log
	w.Bind("kmlOpenUrl", func(u string) error { return openURL(u) })
	w.Bind("kmlOpenLog", func() error { return openLogFolder() })
	w.Bind("kmlOpenNotices", func() error { return openNotices() })
	w.Bind("kmlLogError", func(t string) { logError("the page: " + t) })
	w.Init(pageErrorScript)

	go func() {
		defer stopOnPanic()
		checkForUpdate(w)
	}()

	// the start screen while it finds out where to open (online, the Toolkit's sign-in)
	w.SetHtml(splashPage(appName))
	go func() {
		defer stopOnPanic()
		open := whereToOpen(dataPath, marker)
		w.Dispatch(func() { open(w) })
	}()
	w.Run()
}
