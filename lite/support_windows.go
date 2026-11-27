//go:build windows

package main

// Help and support: the error log, links opened in the browser, and the open-source notices.
// Errors (the program's own and the page's) go into logs\errors.log in the program's folder
// (%LOCALAPPDATA%\TraversePro), 1 MB at most; the older half is kept as errors.log.1.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var logMu sync.Mutex

// programDir is the folder TraverseLevel.exe is in (the install folder)
func programDir() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Dir(exe)
	}
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "TraversePro")
}

func logDir() string { return filepath.Join(programDir(), "logs") }

// logError adds one error to the log. A log that can't be written never stops the program.
func logError(text string) {
	logMu.Lock()
	defer logMu.Unlock()
	dir := logDir()
	if os.MkdirAll(dir, 0o755) != nil {
		return
	}
	p := filepath.Join(dir, "errors.log")
	if fi, err := os.Stat(p); err == nil && fi.Size() > 1<<20 {
		_ = os.Remove(p + ".1")
		_ = os.Rename(p, p+".1")
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	if len(text) > 8000 {
		text = text[:8000] + "…"
	}
	fmt.Fprintf(f, "[%s] Traverse & Level %s (Windows)\n%s\n\n", time.Now().Format(time.RFC3339), appVersion, strings.TrimSpace(text))
}

// stopOnPanic is deferred in main and in every goroutine: a crash is written to the log and
// explained, instead of the window just closing
func stopOnPanic() {
	r := recover()
	if r == nil {
		return
	}
	logError(fmt.Sprintf("the program stopped: %v\n%s", r, debug.Stack()))
	warn("Traverse & Level stopped because of an error. We're sorry.\n\n" +
		"The details are in " + filepath.Join(logDir(), "errors.log") + ".\n\n" +
		"Please tell us at www.surveyxdigital.com/contact (Help → Report a problem) and attach that file. " +
		"Your jobs are kept.")
	os.Exit(1)
}

// openURL opens one of the website's pages in the default browser. The window shows the
// website, so only its own pages open: nothing else can be started from it.
func openURL(u string) error {
	site := strings.HasPrefix(u, "https://www.surveyxdigital.com/") || strings.HasPrefix(u, "https://www.exceltokml.com/")
	if !site || strings.ContainsAny(u, "\"\r\n ") {
		return errors.New("only the website's pages open")
	}
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
}

// openLogFolder shows the folder with errors.log in File Explorer
func openLogFolder() error {
	dir := logDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return exec.Command("explorer", dir).Start()
}

// openNotices opens the open-source notices installed next to the program
func openNotices() error {
	p := filepath.Join(programDir(), "THIRD-PARTY-NOTICES.txt")
	if _, err := os.Stat(p); err != nil {
		return err
	}
	return exec.Command("notepad", p).Start()
}

// pageErrorScript sends the page's own errors to the log
const pageErrorScript = `(function () {
  function send(t) { try { if (window.kmlLogError) window.kmlLogError(String(t).slice(0, 4000)); } catch (e) {} }
  window.addEventListener("error", function (e) {
    send((e.message || "error") + " at " + (e.filename || "") + ":" + (e.lineno || 0) + "\n" + (e.error && e.error.stack || ""));
  });
  window.addEventListener("unhandledrejection", function (e) {
    var r = e.reason; send("unhandled: " + (r && (r.stack || r.message) || r));
  });
})();`
