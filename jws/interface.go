package jws

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"hash"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

// Errors for JWS
var (
	ErrInvalidCompactPartsCount  = errors.New("compact JWS format must have three parts")
	ErrInvalidHeaderValue        = errors.New("invalid value for header key")
	ErrInvalidEcdsaSignatureSize = errors.New("invalid signature size of ecdsa algorithm")
	ErrInvalidSignature          = errors.New("invalid signature")
	ErrMissingPrivateKey         = errors.New("missing private key")
	ErrMissingPublicKey          = errors.New("missing public key")
	ErrUnsupportedAlgorithm      = errors.New("unspported algorithm")
)

// EssentialHeader is a set of headers that are already defined in RFC 7515
type EssentialHeader struct {
	Algorithm              jwa.SignatureAlgorithm `json:"alg,omitempty"`
	ContentType            string                 `json:"cty,omitempty"`
	Critical               []string               `json:"crit,omitempty"`
	Jwk                    jwk.Key                `json:"jwk,omitempty"` // public key
	JwkSetURL              *url.URL               `json:"jku,omitempty"`
	KeyID                  string                 `json:"kid,omitempty"`
	Type                   string                 `json:"typ,omitempty"` // e.g. "JWT"
	X509Url                *url.URL               `json:"x5u,omitempty"`
	X509CertChain          []string               `json:"x5c,omitempty"`
	X509CertThumbprint     string                 `json:"x5t,omitempty"`
	X509CertThumbprintS256 string                 `json:"x5t#S256,omitempty"`
}

// Header represents a JWS header.
type Header struct {
	*EssentialHeader `json:"-"`
	PrivateParams    map[string]interface{} `json:"-"`
}

// EncodedHeader represents a header value that is base64 encoded
// in JSON format
type EncodedHeader struct {
	*Header
	// This is a special field. It's ONLY set when parsed from a serialized form.
	// It's used for verification purposes, because header representations (such as
	// JSON key order) may differ from what the source encoded with and what the
	// go json package uses
	//
	// If this field is populated (Len() > 0), it will be used for signature
	// calculation.
	// If you change the header values, make sure to clear this field, too
	Source buffer.Buffer `json:"-"`
}

// PayloadSigner generates signature for the given payload
type PayloadSigner interface {
	PayloadSign([]byte) ([]byte, error)
	PublicHeaders() *Header
	ProtectedHeaders() *Header
	SetPublicHeaders(*Header)
	SetProtectedHeaders(*Header)
	SignatureAlgorithm() jwa.SignatureAlgorithm
}

// Verifier is used to verify the signature against the payload
type Verifier interface {
	Verify(*Message) error
}

// RsaSign is a signer using RSA
type RsaSign struct {
	Public     *Header
	Protected  *Header
	PrivateKey *rsa.PrivateKey
}

// EcdsaSign is a signer using ECDSA
type EcdsaSign struct {
	Public     *Header
	Protected  *Header
	PrivateKey *ecdsa.PrivateKey
}

// MergedHeader is a provides an interface to query both protected
// and public headers
type MergedHeader struct {
	ProtectedHeader *EncodedHeader
	PublicHeader    *Header
}

// Signature represents a signature generated by one of the signers
type Signature struct {
	PublicHeader    *Header        `json:"header"`              // Raw JWS Unprotected Heders
	ProtectedHeader *EncodedHeader `json:"protected,omitempty"` // Base64 encoded JWS Protected Headers
	Signature       buffer.Buffer  `json:"signature"`           // Base64 encoded signature
}

// Message represents a full JWS encoded message. Flattened serialization
// is not supported as a struct, but rather it's represented as a
// Message struct with only one `signature` element
type Message struct {
	Payload    buffer.Buffer `json:"payload"`
	Signatures []Signature   `json:"signatures"`
}

// MultiSign is a signer that creates multiple signatures for the message.
type MultiSign struct {
	Signers []PayloadSigner
}

// HmacSign is a symmetric signer using HMAC
type HmacSign struct {
	Public    *Header
	Protected *Header
	Key       []byte
	hash      func() hash.Hash
}

// Serializer defines the interface for things that can serialize JWS messages
type Serializer interface {
	Serialize(*Message) ([]byte, error)
}

// CompactSerialize is serializer that produces compact JSON JWS representation
type CompactSerialize struct{}

// JSONSerialize is serializer that produces full JSON JWS representation
type JSONSerialize struct {
	Pretty bool
}

// RsaVerify is a sign verifider using RSA
type RsaVerify struct {
	alg    jwa.SignatureAlgorithm
	hash   crypto.Hash
	pubkey *rsa.PublicKey
}

// EcdsaVerify is a sign verifider using ECDSA
type EcdsaVerify struct {
	alg    jwa.SignatureAlgorithm
	hash   crypto.Hash
	pubkey *ecdsa.PublicKey
}

// HmacVerify is a symmetric sign verifier using HMAC
type HmacVerify struct {
	signer *HmacSign
}
