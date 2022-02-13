package main

import (
	"context"
	"github.com/go-pkgz/auth/token"
	"testing"
)

func TestRandom(t *testing.T) {
	s := RandomString(16)
	if s == "11111" {
		t.Errorf("too simple random string")
	}

	if s == RandomString(16) {
		t.Errorf("not random strings")
	}
}

func MakeTextContext(user token.User) context.Context {
	return context.WithValue(context.Background(), "user", user)
}
