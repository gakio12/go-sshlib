// Copyright (c) 2019 Blacknon. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package sshlib

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// CreateAuthMethodPassword returns ssh.AuthMethod generated from password.
//
func CreateAuthMethodPassword(password string) (auth ssh.AuthMethod) {
	return ssh.Password(password)
}

// CreateAuthMethodPublicKey returns ssh.AuthMethod generated from PublicKey.
//
func CreateAuthMethodPublicKey(key, password string) (auth ssh.AuthMethod, err error) {
	signer, err := CreateSignerPublicKey(key, password)
	if err != nil {
		return
	}

	auth = ssh.PublicKeys(signer)
	return
}

// CreateSignerPublicKey returns []ssh.Signer generated from public key.
//
func CreateSignerPublicKey(key, password string) (signer ssh.Signer, err error) {
	// get absolute path
	key = getAbsPath(key)

	// Read PrivateKey file
	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		return
	}

	if password != "" { // password is not empty
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(password))
	} else { // password is empty
		signer, err = ssh.ParsePrivateKey(keyData)
	}

	return
}

// CreateSignerPublicKeyPrompt rapper CreateSignerPKCS11.
// Output a passphase input prompt if the passphase is not entered or incorrect.
//
// TODO(blacknon): Create
// func CreateSignerPublicKeyPrompt() (signer ssh.Signer, err error) {}

// CreateAuthMethodCertificate returns []ssh.Signer generated from Certificate.
//
func CreateAuthMethodCertificate(cert string, keySigner ssh.Signer) (auth ssh.AuthMethod, err error) {
	signer, err := CreateSignerCertificate(cert, keySigner)
	if err != nil {
		return
	}

	auth = ssh.PublicKeys(signer)
	return
}

// CreateSignerCertificate returns []ssh.Signer generated from Certificate.
//
func CreateSignerCertificate(cert string, keySigner ssh.Signer) (certSigner ssh.Signer, err error) {
	// get absolute path
	cert = getAbsPath(cert)

	// Read Cert file
	certData, err := ioutil.ReadFile(cert)
	if err != nil {
		return
	}

	// Create PublicKey from Cert
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey(certData)
	if err != nil {
		return
	}

	// Create Certificate Struct
	certificate, ok := pubkey.(*ssh.Certificate)
	if !ok {
		err = fmt.Errorf("%s\n", "Error: Not create certificate struct data")
		return
	}

	// Create Certificate Signer
	certSigner, err = ssh.NewCertSigner(certificate, keySigner)
	if err != nil {
		return
	}

	return
}

// CreateAuthMethodPKCS11
//
func CreateAuthMethodPKCS11(provider, pin string) (auth []ssh.AuthMethod, err error) {
	signers, err := CreateSignerPKCS11(provider, pin)
	if err != nil {
		return
	}

	for _, signer := range signers {
		auth = append(auth, ssh.PublicKeys(signer))
	}
	return
}

// CreateSignerCertificate returns []ssh.Signer generated from PKCS11 token.
//
func CreateSignerPKCS11(provider, pin string) (signers []ssh.Signer, err error) {
	// get absolute path
	provider = getAbsPath(provider)

	// Create PKCS11 struct
	p11 := new(PKCS11)
	p11.Pkcs11Provider = provider
	p11.PIN = pin

	// Create pkcs11 ctx
	err = p11.CreateCtx()
	if err != nil {
		return
	}

	// Get token label
	err = p11.GetTokenLabel()
	if err != nil {
		return
	}

	// Recreate ctx (pkcs11=>crypto11)
	err = p11.RecreateCtx(p11.Pkcs11Provider)
	if err != nil {
		return
	}

	// Get KeyID
	err = p11.GetKeyID()
	if err != nil {
		return
	}

	// Get crypto.Signer
	cryptoSigners, err := p11.GetCryptoSigner()
	if err != nil {
		return
	}

	// Exchange crypto.signer to ssh.Signer
	for _, cryptoSigner := range cryptoSigners {
		signer, _ := ssh.NewSignerFromSigner(cryptoSigner)
		signers = append(signers, signer)
	}

	return
}

// CreateSignerPKCS11Prompt rapper CreateSignerPKCS11.
// Output a PIN input prompt if the PIN is not entered or incorrect.
//
// TODO(blacknon): Create
// func CreateSignerPKCS11Prompt() (signers []ssh.Signer, err error) {}

// CreateSignerCertificate returns []ssh.Signer generated from ssh-agent.
//
// TODO(blacknon): Create
// func CreateSignerAgent() (signers []ssh.Signer, err error) {}
