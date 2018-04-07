package sshtest

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/gliderlabs/ssh"
)

const (
	Port           = 2222
	Address string = Host + ":2222"
	Host           = "localhost"
)

type Result struct {
	StdOut      string
	StdErr      string
	ExitStatus  int
	AuthFailure bool
}

func SSHServer(t *testing.T, expected *Result, f func()) {
	srv := &ssh.Server{
		Addr: fmt.Sprint(":", Port),
	}
	srv.Handler = func(s ssh.Session) {
		if expected != nil {
			io.WriteString(s, expected.StdOut)
			io.WriteString(s.Stderr(), expected.StdErr)
			s.Exit(expected.ExitStatus)
		} else {
			s.Exit(0)
		}
	}
	srv.SetOption(ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		if expected != nil && expected.AuthFailure {
			return false
		}
		return true
	}))
	srv.SetOption(ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		if expected != nil && expected.AuthFailure {
			return false
		}
		return true
	}))
	go func() {
		srv.ListenAndServe()
	}()
	defer srv.Shutdown(context.Background())
	f()
	srv.Close()
}
