package twofa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	CodeDigits6 = 6
	CodeDigits7 = 7
	CodeDigits8 = 8

	// TotpFormatter TOTP formatter
	// Example:
	// otpauth://totp/NAME?secret=KEY
	TotpFormatter = "otpauth://totp/%s?secret=%s"
)

var (
	power = map[uint32]uint32{
		CodeDigits6: 1e6,
		CodeDigits7: 1e7,
		CodeDigits8: 1e8,
	}
)

type Key struct {
	Raw    []byte
	Digits uint32
	Offset uint32 // counter offset
}

// GetCode returns time based code by encode key
func GetCode(key string) (string, error) {
	raw, err := DecodeKey(key)
	if err != nil {
		return "", err
	}

	code, err := tOTP(raw, time.Now(), CodeDigits6)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", code), nil
}

// GetCodeByRaw returns time based code by raw key
func GetCodeByRaw(key Key) (uint32, error) {
	return tOTP(key.Raw, time.Now(), key.Digits)
}

// EncodeKey returns encoded key
func EncodeKey(raw string, opt ...OptionFunc) string {
	opts := NewOptions(opt...)
	key := ""
	if opts.WithHash {
		digest, err := opts.HashFunc.Hash([]byte(raw))
		if err != nil {
			return ""
		}
		key = strings.ToUpper(base32.StdEncoding.EncodeToString(digest))
	} else {
		key = strings.ToUpper(base32.StdEncoding.EncodeToString([]byte(raw)))
	}

	if opts.WithQR {
		_ = PrintQR(fmt.Sprintf(TotpFormatter, opts.KeyName, key))
	}

	return key
}

// EncodeKeyForIOS returns encoded key for iOS
func EncodeKeyForIOS(raw string, opt ...OptionFunc) string {
	return strings.ReplaceAll(EncodeKey(raw, opt...), "=", "")
}

// GenKey returns random encoded key
func GenKey(opt ...OptionFunc) string {
	randValue, err := newRand()
	if err != nil {
		return ""
	}

	return EncodeKey(string(randValue), opt...)
}

// GenKeyForIOS returns random encoded key for iOS
func GenKeyForIOS(opt ...OptionFunc) string {
	return strings.ReplaceAll(GenKey(opt...), "=", "")
}

// padding returns string with padding for base32
func padding(s string) string {
	return s + strings.Repeat("=", len(s)%8)
}

// DecodeKey returns decoded key
func DecodeKey(key string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(padding(key)))
}

func checkDigits(digits uint32) error {
	_, ok := power[digits]
	if !ok {
		return errors.New("invalid digits")
	}

	return nil
}

func hOTP(key []byte, counter uint64, digits uint32) (uint32, error) {
	if err := checkDigits(digits); err != nil {
		return 0, err
	}

	h := hmac.New(sha1.New, key)
	if err := binary.Write(h, binary.BigEndian, counter); err != nil {
		return 0, err
	}

	digest := h.Sum(nil)
	bin := binary.BigEndian.Uint32(digest[digest[len(digest)-1]&0x0F:]) & 0x7FFFFFFF

	return bin % power[digits], nil
}

func tOTP(key []byte, t time.Time, digits uint32) (uint32, error) {
	return hOTP(key, uint64(t.UnixNano())/30e9, digits)
}
