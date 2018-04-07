package executor

import (
	"testing"

	"github.com/uphy/chacker/executor/sshtest"
)

func TestNewSSHClientFromPassword(t *testing.T) {
	sshtest.SSHServer(t, &sshtest.Result{
		StdOut:      "hello\n",
		StdErr:      "",
		ExitStatus:  0,
		AuthFailure: false,
	}, func() {
		if client, err := NewSSHClientFromPassword(sshtest.Address, "user1", "password"); err != nil {
			t.Error(err)
		} else {
			r, err := client.Exec("echo hello")
			if err != nil {
				t.Error("exec failed: ", err)
			}
			if r.StdOut != "hello\n" {
				t.Error("unexpected stdout")
			}
			if r.ExitStatus != 0 {
				t.Error("unexpected exit status")
			}
		}
	})
}

func TestNewSSHClientFromPrivateKey(t *testing.T) {
	sshtest.SSHServer(t, &sshtest.Result{
		StdOut:      "hello\n",
		StdErr:      "",
		ExitStatus:  0,
		AuthFailure: false,
	}, func() {
		keypair := sshtest.GenerateKeyPair()
		defer keypair.Delete()
		if client, err := NewSSHClientFromPrivateKey(sshtest.Address, "user1", keypair.PrivateKeyFile, nil); err != nil {
			t.Error(err)
		} else {
			r, err := client.Exec("echo hello")
			if err != nil {
				t.Error("exec failed: ", err)
			}
			if r.StdOut != "hello\n" {
				t.Error("unexpected stdout")
			}
			if r.ExitStatus != 0 {
				t.Error("unexpected exit status")
			}
		}
	})
}
