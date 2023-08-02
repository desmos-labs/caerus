package types

import (
	"time"

	"github.com/google/uuid"
)

type Nonce struct {
	DesmosAddress  string
	Value          string
	ExpirationTime time.Time
}

// NewNonce returns a new Nonce instance
func NewNonce(desmosAddress string, value string, expirationTime time.Time) *Nonce {
	return &Nonce{
		DesmosAddress:  desmosAddress,
		Value:          value,
		ExpirationTime: expirationTime,
	}
}

// CreateNonce creates a new nonce for the given Desmos address
func CreateNonce(desmosAddress string) (*Nonce, error) {
	creationTime := time.Now()
	return NewNonce(
		desmosAddress,
		uuid.NewString(),

		// We set an expiration time of 30 minutes in order to avoid the nonce to be used
		// multiple times (e.g. by a malicious user that tries to brute force the nonce).
		creationTime.Add(time.Minute*30),
	), nil
}

// -------------------------------------------------------------------------------------------------------------------

// EncryptedNonce represents a nonce that has been encrypted and safely stored inside the database
type EncryptedNonce struct {
	// DesmosAddress is the address of the user that is associated to the nonce.
	DesmosAddress string

	// EncryptedValue is the SHA-256 encrypted nonce value.
	EncryptedValue string

	// ExpirationTime is the time at which the nonce will expire.
	ExpirationTime time.Time
}

// NewEncryptedNonce returns a new EncryptedNonce instance
func NewEncryptedNonce(desmosAddress string, encryptedValue string, expirationTime time.Time) *EncryptedNonce {
	return &EncryptedNonce{
		DesmosAddress:  desmosAddress,
		EncryptedValue: encryptedValue,
		ExpirationTime: expirationTime,
	}
}

func (e *EncryptedNonce) IsExpired() bool {
	return e.ExpirationTime.Before(time.Now())
}
