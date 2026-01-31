package tokens

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Params Parameters tuned for server-side token verification.
// Adjust based on your latency/CPU budget.
type Argon2Params struct {
	Time    uint32 // iterations
	Memory  uint32 // KiB
	Threads uint8
	KeyLen  uint32
	SaltLen int
}

var DefaultParams = Argon2Params{
	Time:    2,
	Memory:  64 * 1024, // 64 MiB
	Threads: 1,
	KeyLen:  32,
	SaltLen: 16,
}

// GenerateToken returns (tokenToShowUser, tokenIDPrefix, secretPart).
// Store only tokenIDPrefix + hashed verifier in DB, not the token itself.
func GenerateToken(prefixLen int, secretBytes int) (token string, prefix string, secret string, err error) {
	if prefixLen <= 0 || secretBytes <= 0 {
		return "", "", "", errors.New("invalid sizes")
	}

	pb := make([]byte, prefixLen)
	sb := make([]byte, secretBytes)

	if _, err = rand.Read(pb); err != nil {
		return "", "", "", err
	}
	if _, err = rand.Read(sb); err != nil {
		return "", "", "", err
	}

	enc := base64.RawURLEncoding
	prefix = enc.EncodeToString(pb)
	secret = enc.EncodeToString(sb)
	token = prefix + "." + secret
	return token, prefix, secret, nil
}

// HashSecret produces a PHC-like encoded string to store in DB.
// Store this alongside the token prefix (ID).
//
// Format (single string):
// $argon2id$v=19$m=<mem>,t=<time>,p=<threads>$<saltB64>$<hashB64>
func HashSecret(secret string, p Argon2Params) (string, error) {
	if secret == "" {
		return "", errors.New("empty secret")
	}
	if p.SaltLen < 8 {
		return "", errors.New("salt too short")
	}

	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(secret), salt, p.Time, p.Memory, p.Threads, p.KeyLen)

	enc := base64.RawStdEncoding
	saltB64 := enc.EncodeToString(salt)
	hashB64 := enc.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory, p.Time, p.Threads, saltB64, hashB64)

	return encoded, nil
}

// VerifySecret checks whether secret matches the stored encoded verifier.
func VerifySecret(secret string, encoded string) (bool, error) {
	if secret == "" || encoded == "" {
		return false, errors.New("missing input")
	}

	parts := strings.Split(encoded, "$")
	// ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid encoded format")
	}
	if parts[2] != "v=19" {
		return false, errors.New("unsupported argon2 version")
	}

	mem, time, threads, err := parseParams(parts[3])
	if err != nil {
		return false, err
	}

	dec := base64.RawStdEncoding
	salt, err := dec.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("invalid salt b64")
	}
	want, err := dec.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid hash b64")
	}

	got := argon2.IDKey([]byte(secret), salt, time, mem, uint8(threads), uint32(len(want)))

	// Constant-time compare
	if subtle.ConstantTimeCompare(got, want) == 1 {
		return true, nil
	}
	return false, nil
}

func parseParams(s string) (mem uint32, time uint32, threads uint32, err error) {
	// s: "m=65536,t=2,p=1"
	items := strings.Split(s, ",")
	if len(items) != 3 {
		return 0, 0, 0, errors.New("invalid params")
	}

	get := func(prefix string) (uint64, error) {
		for _, it := range items {
			it = strings.TrimSpace(it)
			if strings.HasPrefix(it, prefix) {
				return strconv.ParseUint(strings.TrimPrefix(it, prefix), 10, 32)
			}
		}
		return 0, fmt.Errorf("missing %s", prefix)
	}

	m, e := get("m=")
	if e != nil {
		return 0, 0, 0, e
	}
	t, e := get("t=")
	if e != nil {
		return 0, 0, 0, e
	}
	p, e := get("p=")
	if e != nil {
		return 0, 0, 0, e
	}

	return uint32(m), uint32(t), uint32(p), nil
}
