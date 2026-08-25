package daraja

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func testCert(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test.safaricom.co.ke"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	return key, der
}

func decrypt(t *testing.T, key *rsa.PrivateKey, credential string) string {
	t.Helper()

	raw, err := base64.StdEncoding.DecodeString(credential)
	if err != nil {
		t.Fatalf("credential is not base64: %v", err)
	}

	plain, err := rsa.DecryptPKCS1v15(rand.Reader, key, raw)
	if err != nil {
		t.Fatalf("credential does not decrypt: %v", err)
	}

	return string(plain)
}

func TestGenerateSecurityCredential(t *testing.T) {
	key, der := testCert(t)
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	for _, tc := range []struct {
		name string
		cert []byte
	}{
		{"pem certificate", pemBytes},
		{"der certificate", der},
		{"pkix public key", pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: mustMarshalPKIX(t, &key.PublicKey),
		})},
		{"pkcs1 public key", pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
		})},
	} {
		t.Run(tc.name, func(t *testing.T) {
			credential, err := GenerateSecurityCredential("s3cret", tc.cert)
			if err != nil {
				t.Fatal(err)
			}

			if got := decrypt(t, key, credential); got != "s3cret" {
				t.Fatalf("got %q, want %q", got, "s3cret")
			}
		})
	}
}

// PKCS #1 v1.5 padding is randomised, so the same password must not produce the
// same credential twice - callers should generate one per request.
func TestGenerateSecurityCredentialIsNotDeterministic(t *testing.T) {
	_, der := testCert(t)

	first, err := GenerateSecurityCredential("s3cret", der)
	if err != nil {
		t.Fatal(err)
	}

	second, err := GenerateSecurityCredential("s3cret", der)
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Fatal("credential is deterministic, padding is not being randomised")
	}
}

func TestGenerateSecurityCredentialErrors(t *testing.T) {
	_, der := testCert(t)

	if _, err := GenerateSecurityCredential("", der); err == nil {
		t.Fatal("expected an error for an empty initiator password")
	}

	if _, err := GenerateSecurityCredential("s3cret", []byte("not a certificate")); err == nil {
		t.Fatal("expected an error for an unparseable certificate")
	}
}

func mustMarshalPKIX(t *testing.T, key *rsa.PublicKey) []byte {
	t.Helper()

	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatal(err)
	}

	return der
}
