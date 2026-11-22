package main

// Opened some other way than from the Survey Toolkit (a desktop shortcut, the Start menu) and not
// signed in yet: the Toolkit on this PC, when there is one and it is signed in, gives a one-use
// sign-in code for its account (SurveyToolkit.exe --app-code=<id> --result=<file>), so the user
// signs in once, in the Toolkit.

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

var codePattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// toolkitExe is where the Survey Toolkit is: as its installer recorded it, or its usual folder; "" when it isn't installed
func toolkitExe() string {
	var places []string
	if k, err := registry.OpenKey(registry.CURRENT_USER, `Software\SurveyToolkit`, registry.QUERY_VALUE); err == nil {
		if dir, _, err := k.GetStringValue("InstallDir"); err == nil && dir != "" {
			places = append(places, filepath.Join(dir, "SurveyToolkit.exe"))
		}
		k.Close()
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		places = append(places, filepath.Join(local, "Survey Toolkit", "Program", "SurveyToolkit.exe"))
	}
	for _, p := range places {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

// askToolkitForCode returns a one-use sign-in code for app from the Survey Toolkit on this PC, or
// "" when there is no Toolkit, it isn't signed in, or it doesn't answer within 20 seconds.
func askToolkitForCode(app string) string {
	exe := toolkitExe()
	if exe == "" {
		return ""
	}
	dir, err := os.MkdirTemp("", "toolkit-signin")
	if err != nil {
		return ""
	}
	defer os.RemoveAll(dir)
	out := filepath.Join(dir, "code.json")
	cmd := exec.Command(exe, "--app-code="+app, "--result="+out)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		return ""
	}
	done := make(chan struct{})
	go func() { _ = cmd.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(20 * time.Second):
		_ = cmd.Process.Kill()
		return ""
	}
	b, err := os.ReadFile(out)
	if err != nil {
		return ""
	}
	var r struct {
		Code *string `json:"code"`
	}
	if json.Unmarshal(b, &r) != nil || r.Code == nil || !codePattern.MatchString(*r.Code) {
		return ""
	}
	return *r.Code
}
