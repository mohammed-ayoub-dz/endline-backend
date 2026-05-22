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

type Argon2Params struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLength  uint32
    KeyLength   uint32
}

var DefaultParams = &Argon2Params{
    Memory:      64 * 1024, 
    Iterations:  3,
    Parallelism: 2,       
    SaltLength:  16,
    KeyLength:   32,
}

func GenerateSalt() ([]byte, error) {
    salt := make([]byte, DefaultParams.SaltLength)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }
    return salt, nil
}

func HashPassword(password string) (string, error) {
    salt, err := GenerateSalt()
    if err != nil {
        return "", err
    }

    hash := argon2.IDKey(
        []byte(password),
        salt,
        DefaultParams.Iterations,
        DefaultParams.Memory,
        DefaultParams.Parallelism,
        DefaultParams.KeyLength,
    )

    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version, DefaultParams.Memory, DefaultParams.Iterations,
        DefaultParams.Parallelism, b64Salt, b64Hash)

    return encoded, nil
}

func VerifyPassword(encodedHash, password string) (bool, error) {
    parts := strings.Split(encodedHash, "$")
    if len(parts) != 6 {
        return false, errors.New("invalid hash format")
    }

    var memory uint32
    var iterations uint32
    var parallelism uint8

    _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
    if err != nil {
        return false, err
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return false, err
    }

    hash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false, err
    }

    keyLength := uint32(len(hash))

    newHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

    return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}