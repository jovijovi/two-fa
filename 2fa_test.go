package twofa_test

import (
	"testing"

	twofa "github.com/jovijovi/two-fa"

	"github.com/stretchr/testify/assert"
)

const (
	mockKey    = "JBSWY3DPFQQHO33SNRSCC==="
	mockRawKey = "Hello, world!"
)

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

func TestEncodeKeyForIOS(t *testing.T) {
	key := twofa.EncodeKeyForIOS(mockRawKey)
	t.Log("Key=", key)
}
