package daraja

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

// Identifier types accepted by the APIs that take an IdentifierType /
// RecieverIdentifierType field.
const (
	IdentifierTypeMSISDN      = "1"
	IdentifierTypeTill        = "2"
	IdentifierTypeShortCode   = "4"
	IdentifierTypeOrgReversal = "11"
)

var ErrNoRSAPublicKey = errors.New("daraja: certificate does not contain an RSA public key")

// GenerateSecurityCredential builds the SecurityCredential required by the
// initiator APIs (B2C, Account Balance, Reversal, Transaction Status): the
// initiator password encrypted with M-PESA's public certificate using RSA
// PKCS #1 v1.5, base64 encoded.
//
// cert may be the PEM or DER encoding of the X.509 certificate downloaded from
// the Daraja portal, or a PEM-encoded public key. Sandbox and production use
// different certificates - encrypting with the wrong one is not a local error,
// it surfaces at M-PESA as result code 2001 (invalid initiator information).
//
// The output differs on every call because PKCS #1 v1.5 padding is randomised;
// that is expected, and the credential should be generated per request rather
// than cached.
func GenerateSecurityCredential(initiatorPassword string, cert []byte) (string, error) {
	if initiatorPassword == "" {
		return "", errors.New("daraja: initiator password is empty")
	}

	key, err := parseRSAPublicKey(cert)
	if err != nil {
		return "", err
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(initiatorPassword))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// GenerateSecurityCredentialFromFile reads the certificate at certPath and
// passes it to GenerateSecurityCredential.
func GenerateSecurityCredentialFromFile(initiatorPassword string, certPath string) (string, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return "", err
	}

	return GenerateSecurityCredential(initiatorPassword, cert)
}

func parseRSAPublicKey(cert []byte) (*rsa.PublicKey, error) {
	der := cert

	// A PEM block may hold either the certificate or a bare public key; fall
	// through to treating the input as raw DER when there is no block.
	if block, _ := pem.Decode(cert); block != nil {
		der = block.Bytes

		if block.Type == "PUBLIC KEY" {
			pub, err := x509.ParsePKIXPublicKey(der)
			if err != nil {
				return nil, err
			}

			key, ok := pub.(*rsa.PublicKey)
			if !ok {
				return nil, ErrNoRSAPublicKey
			}

			return key, nil
		}

		if block.Type == "RSA PUBLIC KEY" {
			return x509.ParsePKCS1PublicKey(der)
		}
	}

	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	key, ok := parsed.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNoRSAPublicKey
	}

	return key, nil
}
