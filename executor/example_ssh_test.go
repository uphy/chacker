package executor_test

import (
	"fmt"
	"log"

	"github.com/uphy/chacker/executor"
)

const Password = "ishikura"

func ExampleSSH_Pubkey() {
	client, err := executor.NewSSHClientFromPrivateKey("192.168.100.138:22", "ishikura", "./id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	user, err := client.Exec("whoami")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	// Output:
	// ishikura
}

func ExampleSSH() {
	client, err := executor.NewSSHClientFromPassword("192.168.100.138:22", "ishikura", Password)
	if err != nil {
		log.Fatal(err)
	}
	user, err := client.Exec("whoami")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	// Output:
	// ishikura
}

func ExampleSCP() {
	client, err := executor.NewSSHClientFromPassword("192.168.100.138:22", "ishikura", Password)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Copy("ssh.go", "ssh.go", "0755"); err != nil {
		log.Fatal(err)
	}
	// Output:
	// a
}
