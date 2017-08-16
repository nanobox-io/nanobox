package router

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/nanobox-io/golang-portal-client"

	"github.com/nanobox-io/nanobox/models"
)

func BuildCert(appModel *models.App) ([]portal.CertBundle, error) {
	certs := []portal.CertBundle{}

	var err error
	var key, cert string

	if appModel.Key == "" {
		key, cert, err = generate()
		if err != nil {
			return nil, fmt.Errorf("Failed to generate cert %s", err.Error())
		}

		appModel.Key = key
		appModel.Cert = cert
		err = appModel.Save()
		if err != nil {
			return nil, fmt.Errorf("Failed to save cert %s", err.Error())
		}
	} else {
		key = appModel.Key
		cert = appModel.Cert
	}

	certs = append(certs, portal.CertBundle{
		Cert: cert,
		Key:  key,
	})

	// send to portal
	return certs, nil
}

func generate() (string, string, error) {
	host := "localhost"

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	notBefore := time.Now()

	notAfter := notBefore.Add(365 * 24 * 100 * time.Hour) // 100 years..

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = append(template.DNSNames, host)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return string(key), string(cert), err
}
