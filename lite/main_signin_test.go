package main

import (
	"strings"
	"testing"
)

func TestSignInCode(t *testing.T) {
	good := strings.Repeat("a", 64)
	if got := signInCode([]string{"--x", "--signin-code=" + good}); got != good {
		t.Fatalf("code not read: %q", got)
	}
	for _, args := range [][]string{nil, {"--signin-code="}, {"--signin-code=abc"}, {"--signin-code=" + good + "&evil=1"}} {
		if got := signInCode(args); got != "" {
			t.Fatalf("%v: want no code, got %q", args, got)
		}
	}
	if firstURL("") != startURL {
		t.Fatal("no code: the app opens as before")
	}
	if u := firstURL(good); u != "https://www.exceltokml.com/api/toolkit/app-code/redeem?app=traverse&code="+good {
		t.Fatalf("redeem address: %s", u)
	}
}
