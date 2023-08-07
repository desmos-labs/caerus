package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	DefaultSessionLength = time.Hour * 24 * 7 // 7 days
)

type UserSession struct {
	DesmosAddress string
	Token         string
	CreationDate  time.Time
	ExpiryTime    time.Time
}

// NewUserSession returns a new UserSession instance
func NewUserSession(desmosAddress string, token string, creationDate time.Time, expiryTime time.Time) *UserSession {
	return &UserSession{
		DesmosAddress: desmosAddress,
		Token:         token,
		CreationDate:  creationDate,
		ExpiryTime:    expiryTime,
	}
}

// CreateUserSession creates a new session for the user having the given Desmos address
func CreateUserSession(desmosAddress string) (*UserSession, error) {
	// Get the creation time
	creationTime := time.Now()

	return NewUserSession(
		desmosAddress,
		uuid.NewString(),
		creationTime,
		creationTime.Add(DefaultSessionLength),
	), nil
}

type EncryptedUserSession struct {
	DesmosAddress  string
	EncryptedToken string
	CreationDate   time.Time
	ExpiryTime     *time.Time
}

// NewEncryptedUserSession returns a new EncryptedUserSession instance
func NewEncryptedUserSession(desmosAddress string, token string, creationDate time.Time, expiryTime *time.Time) *EncryptedUserSession {
	return &EncryptedUserSession{
		DesmosAddress:  desmosAddress,
		EncryptedToken: token,
		CreationDate:   creationDate,
		ExpiryTime:     expiryTime,
	}
}

// Refresh refreshes this session instance extending its expiration time
func (s *EncryptedUserSession) Refresh() *EncryptedUserSession {
	if s.ExpiryTime == nil {
		return s
	}

	updatedExpirationTime := time.Now().Add(DefaultSessionLength)
	return NewEncryptedUserSession(
		s.DesmosAddress,
		s.EncryptedToken,
		s.CreationDate,
		&updatedExpirationTime,
	)
}

func (s *EncryptedUserSession) IsExpired() bool {
	return s.ExpiryTime != nil && s.ExpiryTime.Before(time.Now())
}

// Validate checks the validity of this session based on the given Desmos address.
// It returns an error if anything goes wrong, and true if the session should be refresh and/or deleted from the store.
func (s *EncryptedUserSession) Validate() (shouldRefresh bool, shouldDelete bool, err error) {
	if s.IsExpired() {
		return false, true, fmt.Errorf("token expired")
	}

	return s.ExpiryTime != nil, false, nil
}
