package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	DefaultSessionLength = time.Hour * 24 * 7 // 7 days
)

type Session struct {
	DesmosAddress string
	Token         string
	CreationDate  time.Time
	ExpiryTime    time.Time
}

// NewSession returns a new Sessions instance
func NewSession(desmosAddress string, token string, creationDate time.Time, expiryTime time.Time) *Session {
	return &Session{
		DesmosAddress: desmosAddress,
		Token:         token,
		CreationDate:  creationDate,
		ExpiryTime:    expiryTime,
	}
}

// CreateSession creates a new session for the given Desmos address
func CreateSession(desmosAddress string) (*Session, error) {
	// Get the creation time
	creationTime := time.Now()

	return NewSession(
		desmosAddress,
		uuid.NewString(),
		creationTime,
		creationTime.Add(DefaultSessionLength),
	), nil
}

type EncryptedSession struct {
	DesmosAddress  string
	EncryptedToken string
	CreationDate   time.Time
	ExpiryTime     *time.Time
}

// NewEncryptedSession returns a new EncryptedSessions instance
func NewEncryptedSession(desmosAddress string, token string, creationDate time.Time, expiryTime *time.Time) *EncryptedSession {
	return &EncryptedSession{
		DesmosAddress:  desmosAddress,
		EncryptedToken: token,
		CreationDate:   creationDate,
		ExpiryTime:     expiryTime,
	}
}

// Refresh refreshes this session instance extending its expiration time
func (s *EncryptedSession) Refresh() *EncryptedSession {
	if s.ExpiryTime == nil {
		return s
	}

	updatedExpirationTime := time.Now().Add(DefaultSessionLength)
	return NewEncryptedSession(
		s.DesmosAddress,
		s.EncryptedToken,
		s.CreationDate,
		&updatedExpirationTime,
	)
}

func (s *EncryptedSession) IsExpired() bool {
	return s.ExpiryTime != nil && s.ExpiryTime.Before(time.Now())
}

// Validate checks the validity of this session based on the given Desmos address.
// It returns an error if anything goes wrong, and true if the session should be refresh and/or deleted from the store.
func (s *EncryptedSession) Validate() (shouldRefresh bool, shouldDelete bool, err error) {
	if s.IsExpired() {
		return false, true, fmt.Errorf("token expired")
	}

	return s.ExpiryTime != nil, false, nil
}
