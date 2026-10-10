// Traverse Pro - lightweight Windows shell.
// Uses the WebView2 runtime that ships with Windows 10 and 11, so the program
// itself stays small. The whole program is app/index.html, embedded at build time.
package main

import (
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/jchv/go-webview2"
)

//go:embed index.html
var page []byte

func warn(msg string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	mb := user32.NewProc("MessageBoxW")
	t, _ := syscall.UTF16PtrFromString("Traverse Pro")
	m, _ := syscall.UTF16PtrFromString(msg)
	mb.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}

func main() {
	// Serve the page from loopback only. Nothing is exposed to the network.
	// Browser storage is tied to the address, port included, so a fixed port keeps the
	// traverse you were working on between launches. The fallbacks only matter if
	// something else already holds it. (Excel to KML uses 47831-47833.)
	var ln net.Listener
	var err error
	for _, port := range []string{"47851", "47852", "47853", "0"} {
		ln, err = net.Listen("tcp", "127.0.0.1:"+port)
		if err == nil {
			break
		}
	}
	if err != nil {
		warn("Could not start the local viewer.\n\n" + err.Error())
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(page)
	})
	go http.Serve(ln, mux)

	url := fmt.Sprintf("http://127.0.0.1:%d/", ln.Addr().(*net.TCPAddr).Port)

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		// a fixed data folder, so the work kept in the window survives updates
		DataPath:  filepath.Join(os.Getenv("AppData"), "TraversePro"),
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Traverse Pro",
			Width:  1240,
			Height: 840,
			IconId: 2,
			Center: true,
		},
	})
	if w == nil {
		warn("This computer does not have the Microsoft Edge WebView2 runtime.\n\n" +
			"Install the WebView2 Runtime free from Microsoft, then open Traverse Pro again.")
		return
	}
	defer w.Destroy()

	// the page calls this to save any file through the Windows Save As window
	owner := uintptr(w.Window())
	w.Bind("traverseSave", func(name string, b64 string) (string, error) {
		return saveWithDialog(owner, name, b64)
	})

	w.Init("window.__APP_VERSION__ = " + strconv.Quote(appVersion) + ";")

	go checkForUpdate(w)

	w.Navigate(url)
	w.Run()
}
