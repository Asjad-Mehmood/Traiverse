package update

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIsNewer(t *testing.T) {
	cases := []struct {
		remote, local string
		want          bool
	}{
		{"2.0.1", "1.0.1", true}, // Traverse Pro 1.x is offered Traverse & Level
		{"2.0.1", "2.0.1", false},
		{"2.0.1", "2.0.2", false},
		{"v2.1.0", "2.0.9", true},
		{"2.0.10", "2.0.9", true},
		{"2.0.1", "0.0.0", true}, // a build without a stamped version
	}
	for _, c := range cases {
		if got := IsNewer(c.remote, c.local); got != c.want {
			t.Errorf("IsNewer(%q, %q) = %v, want %v", c.remote, c.local, got, c.want)
		}
	}
}

func server(t *testing.T, installer []byte, manifest string) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/latest.json":
			w.Write([]byte(manifest))
		case "/TraverseLevel-Setup.exe":
			w.Write(installer)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(s.Close)
	return s
}

func TestFetchAndDownload(t *testing.T) {
	installer := []byte("MZ pretend installer")
	sum := sha256.Sum256(installer)
	hexSum := hex.EncodeToString(sum[:])
	s := server(t, installer, `{"version":"2.0.1","sha256":"`+hexSum+`"}`)

	u, err := Fetch(s.URL + "/latest.json")
	if err != nil || u.Version != "2.0.1" || u.SHA256 != hexSum {
		t.Fatalf("Fetch = %+v, %v", u, err)
	}
	dest := filepath.Join(t.TempDir(), "setup.exe")
	if err := Download(s.URL+"/TraverseLevel-Setup.exe", dest, u.SHA256); err != nil {
		t.Fatalf("Download: %v", err)
	}
	if b, _ := os.ReadFile(dest); string(b) != string(installer) {
		t.Fatalf("saved file differs")
	}
}

func TestDownloadRejectsWrongChecksum(t *testing.T) {
	s := server(t, []byte("tampered"), "")
	dest := filepath.Join(t.TempDir(), "setup.exe")
	err := Download(s.URL+"/TraverseLevel-Setup.exe", dest, "0000000000000000000000000000000000000000000000000000000000000000")
	if err == nil {
		t.Fatal("a file with the wrong SHA-256 was accepted")
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatal("the rejected file was left on disk")
	}
}

func TestFetchRejectsBadManifest(t *testing.T) {
	for _, m := range []string{`{"version":"2.0.1"}`, `{"sha256":"abc"}`, `not json`} {
		s := server(t, nil, m)
		if _, err := Fetch(s.URL + "/latest.json"); err == nil {
			t.Errorf("manifest %q was accepted", m)
		}
	}
	s := server(t, nil, "")
	if _, err := Fetch(s.URL + "/missing.json"); err == nil {
		t.Error("a missing manifest was accepted")
	}
}
