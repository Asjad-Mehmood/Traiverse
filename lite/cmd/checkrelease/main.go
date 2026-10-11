// checkrelease downloads a published release the way an installed copy does: it reads
// latest.json, then downloads each installer name and checks it against the SHA-256 there.
// The release workflow runs it after publishing; run it by hand with
//
//	go run ./cmd/checkrelease -want 2.0.1
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"traverselevel/update"
)

func main() {
	base := flag.String("base", "https://github.com/Asjad-Mehmood/Traiverse/releases/latest/download/", "where installed copies read releases")
	want := flag.String("want", "", "the version latest.json must name (empty: any)")
	flag.Parse()

	u, err := update.Fetch(*base + "latest.json")
	if err != nil {
		fail("latest.json: %v", err)
	}
	fmt.Printf("latest.json: version %s, sha256 %s\n", u.Version, u.SHA256)
	if *want != "" && u.Version != *want {
		fail("latest.json names %s, expected %s", u.Version, *want)
	}
	if !update.IsNewer(u.Version, "1.0.1") {
		fail("%s is not newer than Traverse Pro 1.0.1, so Traverse Pro copies would not update", u.Version)
	}
	dir, _ := os.MkdirTemp("", "checkrelease")
	defer os.RemoveAll(dir)
	// TraverseLevel-Setup.exe: Traverse & Level copies. TraversePro-Setup.exe: Traverse Pro 1.x.
	for _, name := range []string{"TraverseLevel-Setup.exe", "TraversePro-Setup.exe"} {
		dest := filepath.Join(dir, name)
		if err := update.Download(*base+name, dest, u.SHA256); err != nil {
			fail("%s: %v", name, err)
		}
		st, _ := os.Stat(dest)
		fmt.Printf("%s: downloaded %d bytes, SHA-256 matches\n", name, st.Size())
	}
	fmt.Println("The release downloads and verifies as installed copies will see it.")
}

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", a...)
	os.Exit(1)
}
