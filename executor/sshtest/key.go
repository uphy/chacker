package sshtest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
	"os"
)

type KeyPair struct {
	PrivateKeyFile string
	PublicKeyFile  string
}

func GenerateKeyPair() *KeyPair {
	if keypair, err := generateKeyPair(); err != nil {
		panic(err)
	} else {
		return keypair
	}
}

func generateKeyPair() (*KeyPair, error) {
	reader := rand.Reader
	bitSize := 2048

	// Generate key pair
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	publicKey := key.PublicKey

	// Write private key file
	privateKeyFile, err := ioutil.TempFile("./", "private.pem")
	if err != nil {
		return nil, err
	}
	defer privateKeyFile.Close()
	if err := pem.Encode(privateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return nil, err
	}

	// Write public key file
	publicKeyFile, err := ioutil.TempFile("./", "id_rsa.pub")
	if err != nil {
		return nil, err
	}
	defer publicKeyFile.Close()
	asn1Bytes, err := asn1.Marshal(publicKey)
	if err != nil {
		return nil, err
	}
	if err := pem.Encode(publicKeyFile, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}); err != nil {
		return nil, err
	}
	return &KeyPair{privateKeyFile.Name(), publicKeyFile.Name()}, nil
}

func (k *KeyPair) Delete() {
	os.Remove(k.PrivateKeyFile)
	os.Remove(k.PublicKeyFile)
}
