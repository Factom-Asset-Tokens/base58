// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"crypto/sha256"
	"errors"
)

// ErrChecksum indicates that the checksum of a check-encoded string does not verify against
// the checksum.
var ErrChecksum = errors.New("checksum error")

// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
var ErrInvalidFormat = errors.New("invalid format: version and/or checksum bytes missing")

// checksum: first four bytes of sha256^2
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

// CheckEncode prepends version bytes and appends a four byte checksum.
func CheckEncode(input []byte, version ...byte) string {
	b := make([]byte, len(version)+len(input)+4)
	i := copy(b, version)
	i += copy(b[i:], input)
	cksum := checksum(b[:i])
	copy(b[i:], cksum[:])
	return Encode(b)
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies
// the checksum. It returns the payload and numVersion version bytes.
func CheckDecode(input string,
	numVersion int) (result []byte, version []byte, err error) {
	decoded := Decode(input)
	if numVersion < 0 {
		numVersion = 0
	}
	if len(decoded) < 4+numVersion {
		return nil, nil, ErrInvalidFormat
	}
	version = decoded[0:numVersion]
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, nil, ErrChecksum
	}
	payload := decoded[numVersion : len(decoded)-4]
	result = append(result, payload...)
	return
}
