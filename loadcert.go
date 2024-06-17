package main

import (
	"crypto/tls"
	"os"
)

func loadPEMCertificate(certPath, keyPath string) (*tls.Certificate, error) {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		// try loading pfx format certificate
		cert, err = tls.LoadPKSC12(certData, keyData)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &cert, nil
}
