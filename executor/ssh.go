package executor

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/uphy/chacker/config"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

func NewSSHClientFromHostConfig(host *config.HostConfig) (*SSHClient, error) {
	dest := fmt.Sprint(host.Address, ":", host.Port)
	if host.Key != "" {
		return NewSSHClientFromPrivateKey(dest, host.User, host.Key, []byte(host.PassPhrase))
	}
	return NewSSHClientFromPassword(dest, host.User, host.Password)
}

func NewSSHClientFromPassword(dest string, user string, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", dest, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %v", err)
	}
	return &SSHClient{client}, nil
}

func NewSSHClientFromPrivateKey(dest string, user string, privateKey string, passphrase []byte) (*SSHClient, error) {
	buf, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}
	var key ssh.Signer
	if passphrase == nil || len(passphrase) == 0 {
		key, err = ssh.ParsePrivateKey(buf)
	} else {
		key, err = ssh.ParsePrivateKeyWithPassphrase(buf, passphrase)
	}
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", dest, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return &SSHClient{client}, nil
}

func (s *SSHClient) Download(src, dest string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	in, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	result := make(chan error)
	go func() {
		err := extractTar(in, dest)
		result <- err
	}()
	err = session.Run(fmt.Sprintf(`cd "%s";tar cf - .`, src))
	tarErr := <-result
	if err != nil {
		return err
	}
	return tarErr
}

func extractTar(reader io.Reader, dest string) error {
	tr := tar.NewReader(reader)
	for {
		h, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		outfile := filepath.Join(dest, h.Name)
		if h.FileInfo().IsDir() {
			os.MkdirAll(outfile, 0700)
			continue
		}

		outdir, _ := filepath.Split(outfile)
		if outdir != "" {
			if _, err := os.Stat(outdir); os.IsNotExist(err) {
				if err := os.MkdirAll(outdir, 0700); err != nil {
					return err
				}
			}
		}
		w, err := os.Create(outfile)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, tr); err != nil {
			w.Close()
			return err
		}
		w.Close()
	}
}

// Upload uploads local file to remote.
// permission examples:
// C0644: file with 0644 permission
// D0755: directory with 0755 permission
func (s *SSHClient) Upload(src, dest string, permission string) error {
	info, err := os.Stat(src)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not exist: %v", src)
	}
	fileSize := info.Size()
	filename := path.Base(dest)
	directory := path.Dir(dest)

	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		w, err := session.StdinPipe()
		if err != nil {
			errCh <- fmt.Errorf("failed to pipe stdin: %v", err)
			session.Close()
			return
		}
		defer w.Close()

		// Write permission, file size, destination path
		fmt.Fprintln(w, fmt.Sprint("C", permission), fileSize, filename)

		// Write content
		f, err := os.Open(src)
		if err != nil {
			errCh <- fmt.Errorf("cannot open file: %v", err)
			session.Close()
			return
		}
		if _, err := io.Copy(w, f); err != nil && err != io.EOF {
			errCh <- fmt.Errorf("cannot copy src to dst: %v", err)
			session.Close()
			return
		}
		fmt.Fprint(w, "\x00")
	}()
	if err := session.Run("scp -qrt " + directory); err != nil {
		return fmt.Errorf("failed to run scp: %v", err)
	}
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (s *SSHClient) Exec(command string) (*CommandResult, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(command); err != nil {
		exitStatus := 1
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitStatus = exitErr.Waitmsg.ExitStatus()
		}
		return &CommandResult{
			ExitStatus: exitStatus,
			StdOut:     stdout.String(),
			StdErr:     stderr.String(),
		}, fmt.Errorf("failed to run the command: %v", err)
	}
	return &CommandResult{
		ExitStatus: 0,
		StdOut:     stdout.String(),
		StdErr:     stderr.String(),
	}, nil
}
