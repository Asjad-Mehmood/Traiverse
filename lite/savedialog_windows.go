//go:build windows

package main

// The standard Windows "Save As" window, so every file the program writes goes
// exactly where you choose. The page hands over the file name and the bytes; this
// asks for the folder, then writes the file.

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	comdlg32             = syscall.NewLazyDLL("comdlg32.dll")
	procGetSaveFileNameW = comdlg32.NewProc("GetSaveFileNameW")
)

// OPENFILENAMEW, laid out exactly as Windows expects it on 64-bit.
type openFileName struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        uintptr
	dwReserved        uint32
	flagsEx           uint32
}

const (
	ofnOverwritePrompt = 0x00000002
	ofnNoChangeDir     = 0x00000008
	ofnPathMustExist   = 0x00000800
	ofnExplorer        = 0x00080000
)

var fileKinds = map[string]string{
	".json": "Traverse Pro job",
	".csv":  "CSV sheet",
	".txt":  "Text file",
	".html": "Web page",
	".pdf":  "PDF document",
}

// the folder used last time, so the next Save opens there
var lastDir string

// utf16z turns text into the NUL-separated form the dialog's filter wants.
func utf16z(s string) []uint16 {
	out, _ := syscall.UTF16FromString(s)
	return out
}

// askSavePath shows the Save As window. It returns the chosen path, or "" if the
// user pressed Cancel.
func askSavePath(owner uintptr, suggested string) string {
	ext := strings.ToLower(filepath.Ext(suggested))
	kind := fileKinds[ext]
	if kind == "" {
		kind = "File"
	}
	// "CSV sheet (*.csv)\0*.csv\0All files (*.*)\0*.*\0\0"
	filter := kind + " (*" + ext + ")\x00*" + ext + "\x00All files (*.*)\x00*.*\x00\x00"
	f16 := make([]uint16, 0, len(filter))
	for _, r := range filter {
		f16 = append(f16, uint16(r))
	}

	buf := make([]uint16, 1024)
	copy(buf, utf16z(suggested))

	var dir *uint16
	if lastDir != "" {
		d := utf16z(lastDir)
		dir = &d[0]
	}
	title := utf16z("Save " + suggested)
	var defExt *uint16
	if len(ext) > 1 {
		e := utf16z(ext[1:])
		defExt = &e[0]
	}

	ofn := openFileName{
		hwndOwner:       owner,
		lpstrFilter:     &f16[0],
		nFilterIndex:    1,
		lpstrFile:       &buf[0],
		nMaxFile:        uint32(len(buf)),
		lpstrInitialDir: dir,
		lpstrTitle:      &title[0],
		flags:           ofnOverwritePrompt | ofnNoChangeDir | ofnPathMustExist | ofnExplorer,
		lpstrDefExt:     defExt,
	}
	ofn.lStructSize = uint32(unsafe.Sizeof(ofn))

	ok, _, _ := procGetSaveFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if ok == 0 {
		return "" // Cancel, or the window could not open
	}
	path := syscall.UTF16ToString(buf)
	lastDir = filepath.Dir(path)
	return path
}

// saveWithDialog is what the page calls: name suggested, content as base64.
// It answers with the full path written, or "" when cancelled.
func saveWithDialog(owner uintptr, name string, b64 string) (string, error) {
	path := askSavePath(owner, name)
	if path == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
