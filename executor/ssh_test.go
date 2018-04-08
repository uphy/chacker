package executor

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
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

func TestSSHClientDownload(t *testing.T) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	tw.WriteHeader(&tar.Header{
		Name: "log.txt",
		Size: 5,
	})
	io.WriteString(tw, "hello")
	sshtest.SSHServer(t, &sshtest.Result{
		StdOutReader: buf,
	}, func() {
		client, err := NewSSHClientFromPassword(sshtest.Address, "user1", "user1")
		if err != nil {
			t.Fatal(err)
		}
		dir, err := ioutil.TempDir("", "chackertest")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)
		if err := client.Download("/tmp/hello.txt", dir); err != nil {
			t.Error(err)
		}
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Error(err)
		}
		if len(files) != 1 {
			t.Error("expected single file.")
		}
		file := files[0]
		if file.Name() != "log.txt" {
			t.Error("invalid name")
		}
	})
}
