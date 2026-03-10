package jwtutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func LoadKeys(privatePath, publicPath string) error {
	if publicPath != "" {
		pubBytes, err := os.ReadFile(publicPath)
		if err != nil {
			return err
		}
		block, _ := pem.Decode(pubBytes)
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return err
		}
		publicKey = pub.(*rsa.PublicKey)
	}

	if privatePath != "" {
		privBytes, err := os.ReadFile(privatePath)
		if err != nil {
			return err
		}
		block, _ := pem.Decode(privBytes)
		priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		privateKey = priv.(*rsa.PrivateKey)
	}
	return nil
}

func GetPrivateKey() (*rsa.PrivateKey, error) {
	if privateKey == nil {
		return nil, errors.New("private key not loaded")
	}
	return privateKey, nil
}

func GetPublicKey() (*rsa.PublicKey, error) {
	if publicKey == nil {
		return nil, errors.New("public key not loaded")
	}
	return publicKey, nil
}
