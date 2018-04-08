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
	StdOut       string
	StdOutReader io.Reader
	StdErr       string
	StdErrReader io.Reader
	ExitStatus   int
	AuthFailure  bool
}

func write(w io.Writer, r io.Reader, s string) error {
	if r != nil {
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, s); err != nil {
		return err
	}
	return nil
}

func SSHServer(t *testing.T, expected *Result, f func()) {
	srv := &ssh.Server{
		Addr: fmt.Sprint(":", Port),
	}
	srv.Handler = func(s ssh.Session) {
		if expected != nil {
			write(s, expected.StdOutReader, expected.StdOut)
			write(s, expected.StdErrReader, expected.StdErr)
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
