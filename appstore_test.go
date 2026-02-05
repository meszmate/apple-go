package apple

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Test Infrastructure ---

// testCertChain holds a generated cert chain for testing.
type testCertChain struct {
	rootCert         *x509.Certificate
	rootKey          *ecdsa.PrivateKey
	rootPool         *x509.CertPool
	intermediateCert *x509.Certificate
	intermediateKey  *ecdsa.PrivateKey
	leafCert         *x509.Certificate
	leafKey          *ecdsa.PrivateKey
	// DER-encoded certificates for x5c header
	leafDER         []byte
	intermediateDER []byte
	rootDER         []byte
}

// generateTestCertChain creates a root -> intermediate -> leaf cert chain for testing.
func generateTestCertChain(t *testing.T) *testCertChain {
	t.Helper()

	// Generate root CA
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	rootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Test Root CA",
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	assert.NoError(t, err)

	rootCert, err := x509.ParseCertificate(rootDER)
	assert.NoError(t, err)

	// Generate intermediate CA
	intermediateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	intermediateTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName:   "Test Intermediate CA",
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	intermediateDER, err := x509.CreateCertificate(rand.Reader, intermediateTemplate, rootCert, &intermediateKey.PublicKey, rootKey)
	assert.NoError(t, err)

	intermediateCert, err := x509.ParseCertificate(intermediateDER)
	assert.NoError(t, err)

	// Generate leaf certificate
	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName:   "Test Leaf",
			Organization: []string{"Test"},
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, intermediateCert, &leafKey.PublicKey, intermediateKey)
	assert.NoError(t, err)

	leafCert, err := x509.ParseCertificate(leafDER)
	assert.NoError(t, err)

	rootPool := x509.NewCertPool()
	rootPool.AddCert(rootCert)

	return &testCertChain{
		rootCert:         rootCert,
		rootKey:          rootKey,
		rootPool:         rootPool,
		intermediateCert: intermediateCert,
		intermediateKey:  intermediateKey,
		leafCert:         leafCert,
		leafKey:          leafKey,
		leafDER:          leafDER,
		intermediateDER:  intermediateDER,
		rootDER:          rootDER,
	}
}

// createTestJWS builds a valid JWS token signed with the test leaf key.
func createTestJWS(t *testing.T, chain *testCertChain, payload []byte) string {
	t.Helper()

	header := jwsHeader{
		Alg: "ES256",
		X5C: []string{
			base64.StdEncoding.EncodeToString(chain.leafDER),
			base64.StdEncoding.EncodeToString(chain.intermediateDER),
		},
	}

	headerBytes, err := json.Marshal(header)
	assert.NoError(t, err)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)

	signingInput := headerB64 + "." + payloadB64
	hash := sha256.Sum256([]byte(signingInput))

	signature, err := ecdsa.SignASN1(rand.Reader, chain.leafKey, hash[:])
	assert.NoError(t, err)

	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return headerB64 + "." + payloadB64 + "." + signatureB64
}

// newTestAppStore creates an appStoreNotifications with a custom root cert pool for testing.
func newTestAppStore(pool *x509.CertPool) *appStoreNotifications {
	return &appStoreNotifications{rootCertPool: pool}
}

// --- Error Tests ---

func TestASError_Error(t *testing.T) {
	t.Run("with reason", func(t *testing.T) {
		err := &ASError{Code: ASErrorInvalidPayload, Reason: "bad json"}
		assert.Equal(t, "appstore: INVALID_PAYLOAD: bad json", err.Error())
	})

	t.Run("without reason", func(t *testing.T) {
		err := &ASError{Code: ASErrorInvalidJWS}
		assert.Equal(t, "appstore: INVALID_JWS", err.Error())
	})
}

// --- V1 Tests ---

func TestParseV1_Valid(t *testing.T) {
	payload := []byte(`{
		"notification_type": "DID_RENEW",
		"password": "secret123",
		"environment": "Sandbox",
		"auto_renew_product_id": "com.example.sub",
		"auto_renew_status": "true",
		"bid": "com.example.app",
		"bvrs": "1.0",
		"unified_receipt": {
			"status": 0,
			"environment": "Sandbox",
			"latest_receipt": "base64receipt",
			"latest_receipt_info": [
				{
					"transaction_id": "1000000123",
					"original_transaction_id": "1000000100",
					"product_id": "com.example.sub",
					"quantity": "1",
					"expires_date_ms": "1700000000000",
					"is_trial_period": "false",
					"in_app_ownership_type": "PURCHASED"
				}
			],
			"pending_renewal_info": [
				{
					"auto_renew_product_id": "com.example.sub",
					"auto_renew_status": "1",
					"original_transaction_id": "1000000100",
					"product_id": "com.example.sub"
				}
			]
		}
	}`)

	as := NewAppStoreNotifications()
	n, err := as.ParseV1(payload)
	assert.NoError(t, err)
	assert.Equal(t, ASNotificationTypeV1DidRenew, n.NotificationType)
	assert.Equal(t, "secret123", n.Password)
	assert.Equal(t, ASEnvironmentSandbox, n.Environment)
	assert.Equal(t, "com.example.sub", n.AutoRenewProductID)
	assert.Equal(t, "true", n.AutoRenewStatus)
	assert.Equal(t, "com.example.app", n.BID)
	assert.Equal(t, "1.0", n.BVRS)

	// Unified receipt
	assert.Equal(t, 0, n.UnifiedReceipt.Status)
	assert.Equal(t, ASEnvironmentSandbox, n.UnifiedReceipt.Environment)
	assert.Equal(t, "base64receipt", n.UnifiedReceipt.LatestReceipt)
	assert.Len(t, n.UnifiedReceipt.LatestReceiptInfo, 1)
	assert.Equal(t, "1000000123", n.UnifiedReceipt.LatestReceiptInfo[0].TransactionID)
	assert.Equal(t, "com.example.sub", n.UnifiedReceipt.LatestReceiptInfo[0].ProductID)
	assert.Equal(t, "1", n.UnifiedReceipt.LatestReceiptInfo[0].Quantity)
	assert.Equal(t, "PURCHASED", n.UnifiedReceipt.LatestReceiptInfo[0].InAppOwnershipType)

	// Pending renewal info
	assert.Len(t, n.UnifiedReceipt.PendingRenewalInfo, 1)
	assert.Equal(t, "com.example.sub", n.UnifiedReceipt.PendingRenewalInfo[0].AutoRenewProductID)
	assert.Equal(t, "1", n.UnifiedReceipt.PendingRenewalInfo[0].AutoRenewStatus)
}

func TestParseV1_MinimalFields(t *testing.T) {
	payload := []byte(`{"notification_type": "INITIAL_BUY"}`)

	as := NewAppStoreNotifications()
	n, err := as.ParseV1(payload)
	assert.NoError(t, err)
	assert.Equal(t, ASNotificationTypeV1InitialBuy, n.NotificationType)
	assert.Nil(t, n.UnifiedReceipt)
}

func TestParseV1_InvalidJSON(t *testing.T) {
	as := NewAppStoreNotifications()
	_, err := as.ParseV1([]byte(`{invalid`))
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidPayload, asErr.Code)
}

func TestParseV1_AllNotificationTypes(t *testing.T) {
	types := []ASNotificationTypeV1{
		ASNotificationTypeV1InitialBuy,
		ASNotificationTypeV1Cancel,
		ASNotificationTypeV1DidChangeRenewalPref,
		ASNotificationTypeV1DidChangeRenewalStatus,
		ASNotificationTypeV1DidFailToRenew,
		ASNotificationTypeV1DidRecover,
		ASNotificationTypeV1DidRenew,
		ASNotificationTypeV1InteractiveRenewal,
		ASNotificationTypeV1PriceIncreaseConsent,
		ASNotificationTypeV1Refund,
		ASNotificationTypeV1Revoke,
		ASNotificationTypeV1ConsumptionRequest,
	}

	as := NewAppStoreNotifications()
	for _, nt := range types {
		payload, _ := json.Marshal(map[string]string{"notification_type": string(nt)})
		n, err := as.ParseV1(payload)
		assert.NoError(t, err)
		assert.Equal(t, nt, n.NotificationType)
	}
}

// --- V2 Tests ---

func TestParseV2_Valid(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	v2Payload := ASNotificationV2{
		NotificationType: ASNotificationTypeSubscribed,
		Subtype:          ASSubtypeInitialBuy,
		NotificationUUID: "test-uuid-1234",
		Version:          "2.0",
		SignedDate:       1700000000000,
		Data: &ASNotificationData{
			AppAppleID:    123456789,
			BundleID:      "com.example.app",
			BundleVersion: "1.0",
			Environment:   ASEnvironmentSandbox,
			Status:        ASStatusActive,
		},
	}

	payloadBytes, err := json.Marshal(v2Payload)
	assert.NoError(t, err)

	jws := createTestJWS(t, chain, payloadBytes)

	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: jws})

	n, err := as.ParseV2(envelope)
	assert.NoError(t, err)
	assert.Equal(t, ASNotificationTypeSubscribed, n.NotificationType)
	assert.Equal(t, ASSubtypeInitialBuy, n.Subtype)
	assert.Equal(t, "test-uuid-1234", n.NotificationUUID)
	assert.Equal(t, "2.0", n.Version)
	assert.Equal(t, int64(1700000000000), n.SignedDate)
	assert.Equal(t, int64(123456789), n.Data.AppAppleID)
	assert.Equal(t, "com.example.app", n.Data.BundleID)
	assert.Equal(t, ASEnvironmentSandbox, n.Data.Environment)
	assert.Equal(t, ASStatusActive, n.Data.Status)
}

func TestParseV2_InvalidJSON(t *testing.T) {
	as := NewAppStoreNotifications()
	_, err := as.ParseV2([]byte(`{invalid`))
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidPayload, asErr.Code)
}

func TestParseV2_EmptySignedPayload(t *testing.T) {
	as := NewAppStoreNotifications()
	_, err := as.ParseV2([]byte(`{"signedPayload": ""}`))
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidPayload, asErr.Code)
	assert.Contains(t, asErr.Reason, "missing signedPayload")
}

func TestParseV2_MalformedJWS(t *testing.T) {
	as := NewAppStoreNotifications()

	// Not enough parts
	payload, _ := json.Marshal(ASSignedPayload{SignedPayload: "only.two"})
	_, err := as.ParseV2(payload)
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidJWS, asErr.Code)
}

func TestParseV2_InvalidSignature(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	v2Payload, _ := json.Marshal(ASNotificationV2{
		NotificationType: ASNotificationTypeTest,
		Version:          "2.0",
	})

	jws := createTestJWS(t, chain, v2Payload)

	// Corrupt the signature
	parts := splitJWS(jws)
	parts[2] = "AAAA" + parts[2][4:]
	corruptedJWS := parts[0] + "." + parts[1] + "." + parts[2]

	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: corruptedJWS})

	_, err := as.ParseV2(envelope)
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorSignatureInvalid, asErr.Code)
}

func TestParseV2_WrongRootCA(t *testing.T) {
	chain := generateTestCertChain(t)

	// Use the default Apple Root CA (won't match test chain)
	as := NewAppStoreNotifications()

	v2Payload, _ := json.Marshal(ASNotificationV2{
		NotificationType: ASNotificationTypeTest,
		Version:          "2.0",
	})

	jws := createTestJWS(t, chain, v2Payload)
	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: jws})

	_, err := as.ParseV2(envelope)
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidCertChain, asErr.Code)
}

func TestParseV2_WrongAlgorithm(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	// Build a JWS with wrong algorithm in header
	header := jwsHeader{
		Alg: "RS256",
		X5C: []string{
			base64.StdEncoding.EncodeToString(chain.leafDER),
			base64.StdEncoding.EncodeToString(chain.intermediateDER),
		},
	}

	headerBytes, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadB64 := base64.RawURLEncoding.EncodeToString([]byte(`{}`))
	sigB64 := base64.RawURLEncoding.EncodeToString([]byte("fakesig"))

	jws := headerB64 + "." + payloadB64 + "." + sigB64
	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: jws})

	_, err := as.ParseV2(envelope)
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorUnsupportedAlgo, asErr.Code)
}

func TestParseV2_EmptyX5C(t *testing.T) {
	as := NewAppStoreNotifications()

	header := jwsHeader{Alg: "ES256", X5C: []string{}}
	headerBytes, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadB64 := base64.RawURLEncoding.EncodeToString([]byte(`{}`))
	sigB64 := base64.RawURLEncoding.EncodeToString([]byte("fakesig"))

	jws := headerB64 + "." + payloadB64 + "." + sigB64
	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: jws})

	_, err := as.ParseV2(envelope)
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorCertificateInvalid, asErr.Code)
}

// --- DecodeTransactionInfo Tests ---

func TestDecodeTransactionInfo_Valid(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	txn := ASTransactionInfo{
		TransactionID:         "1000000123456",
		OriginalTransactionID: "1000000100000",
		BundleID:              "com.example.app",
		ProductID:             "com.example.sub.monthly",
		PurchaseDate:          1700000000000,
		ExpiresDate:           1702592000000,
		Type:                  ASTransactionTypeAutoRenewable,
		InAppOwnershipType:    ASOwnershipTypePurchased,
		Environment:           ASEnvironmentSandbox,
		TransactionReason:     ASTransactionReasonPurchase,
		Storefront:            "USA",
		StorefrontID:          "143441",
		Quantity:              1,
		SignedDate:            1700000000000,
	}

	payloadBytes, _ := json.Marshal(txn)
	jws := createTestJWS(t, chain, payloadBytes)

	result, err := as.DecodeTransactionInfo(jws)
	assert.NoError(t, err)
	assert.Equal(t, "1000000123456", result.TransactionID)
	assert.Equal(t, "1000000100000", result.OriginalTransactionID)
	assert.Equal(t, "com.example.app", result.BundleID)
	assert.Equal(t, "com.example.sub.monthly", result.ProductID)
	assert.Equal(t, int64(1700000000000), result.PurchaseDate)
	assert.Equal(t, int64(1702592000000), result.ExpiresDate)
	assert.Equal(t, ASTransactionTypeAutoRenewable, result.Type)
	assert.Equal(t, ASOwnershipTypePurchased, result.InAppOwnershipType)
	assert.Equal(t, ASEnvironmentSandbox, result.Environment)
	assert.Equal(t, ASTransactionReasonPurchase, result.TransactionReason)
	assert.Equal(t, "USA", result.Storefront)
	assert.Equal(t, int32(1), result.Quantity)
}

func TestDecodeTransactionInfo_InvalidJWS(t *testing.T) {
	as := NewAppStoreNotifications()
	_, err := as.DecodeTransactionInfo("not-a-jws")
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidJWS, asErr.Code)
}

// --- DecodeRenewalInfo Tests ---

func TestDecodeRenewalInfo_Valid(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	renewal := ASRenewalInfo{
		AutoRenewProductID:    "com.example.sub.monthly",
		AutoRenewStatus:       1,
		Environment:           ASEnvironmentProduction,
		OriginalTransactionID: "1000000100000",
		ProductID:             "com.example.sub.monthly",
		RenewalDate:           1702592000000,
		SignedDate:            1700000000000,
	}

	payloadBytes, _ := json.Marshal(renewal)
	jws := createTestJWS(t, chain, payloadBytes)

	result, err := as.DecodeRenewalInfo(jws)
	assert.NoError(t, err)
	assert.Equal(t, "com.example.sub.monthly", result.AutoRenewProductID)
	assert.Equal(t, int32(1), result.AutoRenewStatus)
	assert.Equal(t, ASEnvironmentProduction, result.Environment)
	assert.Equal(t, "1000000100000", result.OriginalTransactionID)
	assert.Equal(t, int64(1702592000000), result.RenewalDate)
}

func TestDecodeRenewalInfo_InvalidJWS(t *testing.T) {
	as := NewAppStoreNotifications()
	_, err := as.DecodeRenewalInfo("bad")
	assert.Error(t, err)

	var asErr *ASError
	assert.ErrorAs(t, err, &asErr)
	assert.Equal(t, ASErrorInvalidJWS, asErr.Code)
}

// --- V2 Notification with Summary ---

func TestParseV2_WithSummary(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	v2Payload := ASNotificationV2{
		NotificationType: ASNotificationTypeRenewalExtension,
		Subtype:          ASSubtypeSummary,
		NotificationUUID: "summary-uuid",
		Version:          "2.0",
		SignedDate:       1700000000000,
		Summary: &ASNotificationSummary{
			RequestIdentifier:      "req-123",
			Environment:            ASEnvironmentProduction,
			AppAppleID:             123456789,
			BundleID:               "com.example.app",
			ProductID:              "com.example.sub",
			StorefrontCountryCodes: []string{"USA", "GBR"},
			SucceededCount:         100,
			FailedCount:            5,
		},
	}

	payloadBytes, _ := json.Marshal(v2Payload)
	jws := createTestJWS(t, chain, payloadBytes)
	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: jws})

	n, err := as.ParseV2(envelope)
	assert.NoError(t, err)
	assert.Equal(t, ASNotificationTypeRenewalExtension, n.NotificationType)
	assert.Equal(t, ASSubtypeSummary, n.Subtype)
	assert.NotNil(t, n.Summary)
	assert.Equal(t, "req-123", n.Summary.RequestIdentifier)
	assert.Equal(t, int64(100), n.Summary.SucceededCount)
	assert.Equal(t, int64(5), n.Summary.FailedCount)
	assert.Equal(t, []string{"USA", "GBR"}, n.Summary.StorefrontCountryCodes)
}

// --- V2 Notification with ExternalPurchaseToken ---

func TestParseV2_WithExternalPurchaseToken(t *testing.T) {
	chain := generateTestCertChain(t)
	as := newTestAppStore(chain.rootPool)

	v2Payload := ASNotificationV2{
		NotificationType: ASNotificationTypeExternalPurchaseToken,
		NotificationUUID: "ext-uuid",
		Version:          "2.0",
		SignedDate:       1700000000000,
		ExternalPurchaseToken: &ASExternalPurchaseToken{
			ExternalPurchaseID: "ext-purchase-123",
			TokenCreationDate:  1700000000000,
			AppAppleID:         987654321,
			BundleID:           "com.example.app",
			Environment:        ASEnvironmentSandbox,
		},
	}

	payloadBytes, _ := json.Marshal(v2Payload)
	jws := createTestJWS(t, chain, payloadBytes)
	envelope, _ := json.Marshal(ASSignedPayload{SignedPayload: jws})

	n, err := as.ParseV2(envelope)
	assert.NoError(t, err)
	assert.Equal(t, ASNotificationTypeExternalPurchaseToken, n.NotificationType)
	assert.NotNil(t, n.ExternalPurchaseToken)
	assert.Equal(t, "ext-purchase-123", n.ExternalPurchaseToken.ExternalPurchaseID)
	assert.Equal(t, int64(987654321), n.ExternalPurchaseToken.AppAppleID)
}

// --- Certificate Pool Test ---

func TestAppleRootCertPool(t *testing.T) {
	pool := appleRootCertPool()
	assert.NotNil(t, pool)
}

// --- Helper ---

func splitJWS(jws string) [3]string {
	var parts [3]string
	first := 0
	idx := 0
	for i := 0; i < len(jws); i++ {
		if jws[i] == '.' {
			parts[idx] = jws[first:i]
			idx++
			first = i + 1
		}
	}
	parts[idx] = jws[first:]
	return parts
}
