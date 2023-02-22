package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type Params struct {
	Memory      uint32
	Iter        uint32
	Parallelism uint8
	LenSalt     uint32
	LenKey      uint32
}

func GenerateHashedPassword(pwd string, p *Params) (string, error) {
	salt, err := generateRandomBytes(p.LenSalt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(pwd),
		salt,
		p.Iter,
		p.Memory,
		p.Parallelism,
		p.LenKey,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iter,
		p.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func CompareSavedWithIncomingPassword(incomingPwd, savedPwd string) (bool, error) {
	p, salt, savedHash, err := decodeHash(savedPwd)
	if err != nil {
		return false, err
	}

	incomingHash := argon2.IDKey(
		[]byte(incomingPwd),
		salt,
		p.Iter,
		p.Memory,
		p.Parallelism,
		p.LenKey,
	)

	if subtle.ConstantTimeCompare(savedHash, incomingHash) == 1 {
		return true, nil
	}

	return false, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func decodeHash(encodedHash string) (p *Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iter, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.LenSalt = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.LenKey = uint32(len(hash))

	return p, salt, hash, nil
}
