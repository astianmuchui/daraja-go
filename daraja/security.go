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



const (
	IdentifierTypeMSISDN      = "1"
	IdentifierTypeTill        = "2"
	IdentifierTypeShortCode   = "4"
	IdentifierTypeOrgReversal = "11"
)

var ErrNoRSAPublicKey = errors.New("daraja: certificate does not contain an RSA public key")





//




//



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



func GenerateSecurityCredentialFromFile(initiatorPassword string, certPath string) (string, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return "", err
	}

	return GenerateSecurityCredential(initiatorPassword, cert)
}

func parseRSAPublicKey(cert []byte) (*rsa.PublicKey, error) {
	der := cert

	
	
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
