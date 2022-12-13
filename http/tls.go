package http

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"github.com/shimmeris/SCFProxy/fileutil"

	"github.com/google/martian/v3/mitm"
)

func GetX509KeyPair(certPath, keyPath string) (*x509.Certificate, crypto.PrivateKey, error) {

	if !fileutil.PathExists(certPath) || !fileutil.PathExists(keyPath) {
		if err := GenerateCert(certPath, keyPath); err != nil {
			return nil, nil, err
		}
	}

	tlsc, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, nil, err
	}
	pk := tlsc.PrivateKey

	cert, err := x509.ParseCertificate(tlsc.Certificate[0])
	if err != nil {
		return nil, nil, err
	}

	return cert, pk, err
}

func GenerateCert(certPath, keyPath string) error {
	cert, pk, err := mitm.NewAuthority("SCFProxy", "Martian Authority", 365*24*time.Hour)
	if err != nil {
		return err
	}

	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		return err
	}
	return pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

}
