package aes

import (
  "errors"
  "fmt"
  "testing"
)

func TestRotWord(t *testing.T) {
  in := []byte {
    0x00, 0x11, 0x22, 0x33,
  }
  out := rotWord(in)
  expected := []byte {
    0x11, 0x22, 0x33, 0x00,
  }
  if ok, err := compare(out, expected); !ok {
    t.Errorf("rotWord(0x%X) returned 0x%X. Expected 0x%X. Error: %q", in, out, expected, err.Error())
  }
}

func TestEncrypt(t *testing.T) {
  textIn := []byte {
    0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
    0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
  }
  keyIn := []byte {
    0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
    0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
  }
  out := Encrypt(textIn, keyIn)
  expected := []byte {
    0x69, 0xc4, 0xe0, 0xd8, 0x6a, 0x7b, 0x04, 0x30,
    0xd8, 0xcd, 0xb7, 0x80, 0x70, 0xb4, 0xc5, 0x5a,
  }
  if ok, err := compare(out, expected); !ok {
    t.Errorf("Encrypt(0x%X, 0x%X) returned 0x%X. Expected 0x%X. Error: %q", textIn, keyIn, out, expected, err.Error())
  }
}

func compare(a, b []byte) (bool, error) {
  if len(a) != len(b) {
    return false, errors.New("lengths differ")
  }
  for i, v := range a {
    if v != b[i] {
      return false, errors.New(fmt.Sprintf("error at index %d", i))
    }
  }
  return true, nil
}