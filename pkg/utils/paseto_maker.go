package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type PayloadSession struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayloadSession(email string, duration time.Duration) (*PayloadSession, error) {
	tokenID := uuid.NewString()

	payloadSession := &PayloadSession{
		ID:        tokenID,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payloadSession, nil
}

func (payloadSession *PayloadSession) Valid() error {
	if time.Now().After(payloadSession.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

type Maker interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*PayloadSession, error)
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (p *PasetoMaker) CreateToken(email string, duration time.Duration) (string, error) {
	payloadSession, err := NewPayloadSession(email, duration)
	if err != nil {
		return "", err
	}

	return p.paseto.Encrypt(p.symmetricKey, payloadSession, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*PayloadSession, error) {
	payloadSession := &PayloadSession{}

	err := p.paseto.Decrypt(token, p.symmetricKey, payloadSession, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payloadSession.Valid()
	if err != nil {
		return nil, err
	}

	return payloadSession, nil
}
