package apple

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// AppStoreNotifications provides methods for parsing and verifying
// Apple App Store Server Notifications (V1 and V2).
type AppStoreNotifications interface {
	// ParseV1 parses a V1 App Store Server Notification from raw JSON bytes.
	ParseV1(payload []byte) (*ASNotificationV1, error)

	// ParseV2 parses and verifies a V2 App Store Server Notification.
	// The payload is the raw JSON body containing a signedPayload field.
	// The JWS signature is verified against the Apple Root CA - G3 certificate chain.
	ParseV2(payload []byte) (*ASNotificationV2, error)

	// DecodeTransactionInfo decodes and verifies a signed transaction info JWS string.
	DecodeTransactionInfo(signedTransactionInfo string) (*ASTransactionInfo, error)

	// DecodeRenewalInfo decodes and verifies a signed renewal info JWS string.
	DecodeRenewalInfo(signedRenewalInfo string) (*ASRenewalInfo, error)
}

type appStoreNotifications struct {
	rootCertPool *x509.CertPool
}

// NewAppStoreNotifications creates a new AppStoreNotifications instance.
// The returned instance is safe for concurrent use.
func NewAppStoreNotifications() AppStoreNotifications {
	return &appStoreNotifications{
		rootCertPool: appleRootCertPool(),
	}
}

// ParseV1 parses a V1 App Store Server Notification from raw JSON bytes.
func (a *appStoreNotifications) ParseV1(payload []byte) (*ASNotificationV1, error) {
	var notification ASNotificationV1
	if err := json.Unmarshal(payload, &notification); err != nil {
		return nil, &ASError{Code: ASErrorInvalidPayload, Reason: err.Error()}
	}
	return &notification, nil
}

// ParseV2 parses and verifies a V2 App Store Server Notification.
func (a *appStoreNotifications) ParseV2(payload []byte) (*ASNotificationV2, error) {
	var signed ASSignedPayload
	if err := json.Unmarshal(payload, &signed); err != nil {
		return nil, &ASError{Code: ASErrorInvalidPayload, Reason: err.Error()}
	}
	if signed.SignedPayload == "" {
		return nil, &ASError{Code: ASErrorInvalidPayload, Reason: "missing signedPayload"}
	}

	decoded, err := verifyAndDecodeJWS(signed.SignedPayload, a.rootCertPool)
	if err != nil {
		return nil, err
	}

	var notification ASNotificationV2
	if err := json.Unmarshal(decoded, &notification); err != nil {
		return nil, &ASError{Code: ASErrorInvalidPayload, Reason: err.Error()}
	}
	return &notification, nil
}

// DecodeTransactionInfo decodes and verifies a signed transaction info JWS string.
func (a *appStoreNotifications) DecodeTransactionInfo(signedTransactionInfo string) (*ASTransactionInfo, error) {
	decoded, err := verifyAndDecodeJWS(signedTransactionInfo, a.rootCertPool)
	if err != nil {
		return nil, err
	}

	var txn ASTransactionInfo
	if err := json.Unmarshal(decoded, &txn); err != nil {
		return nil, &ASError{Code: ASErrorInvalidPayload, Reason: err.Error()}
	}
	return &txn, nil
}

// DecodeRenewalInfo decodes and verifies a signed renewal info JWS string.
func (a *appStoreNotifications) DecodeRenewalInfo(signedRenewalInfo string) (*ASRenewalInfo, error) {
	decoded, err := verifyAndDecodeJWS(signedRenewalInfo, a.rootCertPool)
	if err != nil {
		return nil, err
	}

	var renewal ASRenewalInfo
	if err := json.Unmarshal(decoded, &renewal); err != nil {
		return nil, &ASError{Code: ASErrorInvalidPayload, Reason: err.Error()}
	}
	return &renewal, nil
}

// jwsHeader represents the JOSE header of a JWS token.
type jwsHeader struct {
	Alg string   `json:"alg"`
	X5C []string `json:"x5c"`
}

// verifyAndDecodeJWS verifies a JWS token's signature against the provided root CA
// certificate pool and returns the decoded payload.
func verifyAndDecodeJWS(token string, rootCertPool *x509.CertPool) ([]byte, error) {
	// Split into header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, &ASError{Code: ASErrorInvalidJWS, Reason: "expected 3 parts"}
	}

	headerB64, payloadB64, signatureB64 := parts[0], parts[1], parts[2]

	// Decode and parse header
	headerBytes, err := base64.RawURLEncoding.DecodeString(headerB64)
	if err != nil {
		return nil, &ASError{Code: ASErrorInvalidJWS, Reason: "invalid header encoding"}
	}

	var header jwsHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, &ASError{Code: ASErrorInvalidJWS, Reason: "invalid header JSON"}
	}

	// Verify algorithm
	if header.Alg != "ES256" {
		return nil, &ASError{Code: ASErrorUnsupportedAlgo, Reason: "expected ES256, got " + header.Alg}
	}

	// Decode x5c certificates
	if len(header.X5C) == 0 {
		return nil, &ASError{Code: ASErrorCertificateInvalid, Reason: "empty x5c chain"}
	}

	certs := make([]*x509.Certificate, len(header.X5C))
	for i, certB64 := range header.X5C {
		certDER, err := base64.StdEncoding.DecodeString(certB64)
		if err != nil {
			return nil, &ASError{Code: ASErrorCertificateInvalid, Reason: "invalid x5c certificate encoding"}
		}
		cert, err := x509.ParseCertificate(certDER)
		if err != nil {
			return nil, &ASError{Code: ASErrorCertificateInvalid, Reason: "invalid x5c certificate"}
		}
		certs[i] = cert
	}

	// Build intermediate pool from non-leaf certs
	intermediates := x509.NewCertPool()
	for _, cert := range certs[1:] {
		intermediates.AddCert(cert)
	}

	// Verify the leaf certificate against the Apple Root CA - G3
	leaf := certs[0]
	_, err = leaf.Verify(x509.VerifyOptions{
		Roots:         rootCertPool,
		Intermediates: intermediates,
	})
	if err != nil {
		return nil, &ASError{Code: ASErrorInvalidCertChain, Reason: err.Error()}
	}

	// Extract ECDSA public key from leaf
	pubKey, ok := leaf.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, &ASError{Code: ASErrorCertificateInvalid, Reason: "leaf certificate does not contain an ECDSA public key"}
	}

	// Verify ES256 signature
	signingInput := headerB64 + "." + payloadB64
	hash := sha256.Sum256([]byte(signingInput))

	signatureBytes, err := base64.RawURLEncoding.DecodeString(signatureB64)
	if err != nil {
		return nil, &ASError{Code: ASErrorSignatureInvalid, Reason: "invalid signature encoding"}
	}

	if !ecdsa.VerifyASN1(pubKey, hash[:], signatureBytes) {
		return nil, &ASError{Code: ASErrorSignatureInvalid, Reason: "signature verification failed"}
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, &ASError{Code: ASErrorDecodeError, Reason: "invalid payload encoding"}
	}

	return payloadBytes, nil
}
