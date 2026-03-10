package jwtutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func LoadKeys(privatePath, publicPath string) error {
	if publicPath != "" {
		pubBytes, err := os.ReadFile(publicPath)
		if err != nil {
			logrus.Errorf("Failed to read public key: %v", err)
			return err
		}
		logrus.Infof("Public key file size: %d bytes", len(pubBytes))
		block, _ := pem.Decode(pubBytes)
		if block == nil {
			logrus.Error("Failed to parse PEM block for public key")
			return errors.New("failed to parse PEM block for public key")
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logrus.Errorf("Failed to parse public key: %v", err)
			return err
		}
		publicKey = pub.(*rsa.PublicKey)
		logrus.Info("Public key loaded successfully")
	}

	if privatePath != "" {
		privBytes, err := os.ReadFile(privatePath)
		if err != nil {
			return err
		}
		block, _ := pem.Decode(privBytes)
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		privateKey = priv
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
