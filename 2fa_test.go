package twofa_test

import (
	"crypto/sha256"
	"testing"

	twofa "github.com/jovijovi/two-fa"

	"github.com/stretchr/testify/assert"
)

const (
	mockKey         = "JBSWY3DPFQQHO33SNRSCC==="
	mockRawKey      = "Hello, world!"
	mockKeyWithHash = "SQ5HALIG6NCZTLXB7DNI56PXFFQDDVUZ"
)

type CustomHashFunc struct {
	twofa.HashFunc
}

// Custom impl of hash
func (h *CustomHashFunc) Hash(msg []byte) ([]byte, error) {
	provider := h.Provider()
	if _, err := provider.Write(msg); err != nil {
		return nil, err
	}

	return provider.Sum(nil), nil
}

// GetCustomHashFunc returns custom hash func
func GetCustomHashFunc() twofa.IHashFunc {
	customHashFunc := new(CustomHashFunc)
	customHashFunc.Provider = sha256.New
	return customHashFunc
}

func TestGetCode(t *testing.T) {
	code, err := twofa.GetCode(mockKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Code=", code)
}

func TestGetCodeByRawKey(t *testing.T) {
	code, err := twofa.GetCodeByRaw(twofa.Key{
		Raw:    []byte(mockRawKey),
		Digits: twofa.CodeDigits6,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Code=", code)

	_, err = twofa.GetCodeByRaw(twofa.Key{
		Raw:    []byte(mockRawKey),
		Digits: 0,
	})
	assert.NotEmpty(t, err)
	t.Log("Error=", err.Error())

	_, err = twofa.GetCodeByRaw(twofa.Key{
		Raw:    []byte(mockRawKey),
		Digits: 9,
	})
	assert.NotEmpty(t, err)
	t.Log("Error=", err.Error())
}

func TestEncodeKey(t *testing.T) {
	key := twofa.EncodeKey(mockRawKey)
	t.Log("Key=", key)
	assert.Equal(t, key, mockKey)

	code, err := twofa.GetCode(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Code=", code)

	raw, err := twofa.DecodeKey(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Raw=", string(raw))
	assert.Equal(t, mockRawKey, string(raw))
}

func TestEncodeKeyWithHash(t *testing.T) {
	key := twofa.EncodeKey(mockRawKey, twofa.WithDefaultHashFunc())
	t.Log("Key=", key)
	assert.Equal(t, key, mockKeyWithHash)

	code, err := twofa.GetCode(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Code=", code)

	rawHash, err := twofa.DecodeKey(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RawHash=%x", string(rawHash))
}

func TestGenKey(t *testing.T) {
	key := twofa.GenKey()
	t.Log("Key=", key)

	code, err := twofa.GetCode(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Code=", code)

	raw, err := twofa.DecodeKey(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Raw=%x", string(raw))
}

func TestEncodeKeyForIOS(t *testing.T) {
	key := twofa.EncodeKeyForIOS(mockRawKey)
	t.Log("Key=", key)
}

func TestEncodeKeyWithHashForIOS(t *testing.T) {
	key := twofa.EncodeKeyForIOS(mockRawKey, twofa.WithHashFunc(GetCustomHashFunc()))
	t.Log("Key=", key)

	rawHash, err := twofa.DecodeKey(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RawHash=%x", string(rawHash))
}

func TestGenKeyForIOS(t *testing.T) {
	key := twofa.GenKeyForIOS()
	t.Log("Key=", key)

	code, err := twofa.GetCode(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Code=", code)

	raw, err := twofa.DecodeKey(key)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Raw=%x", string(raw))
}
