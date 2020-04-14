// The MIT License (MIT)
//
// Copyright (c) 2020 Ryan Fowler
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package uuid provides functions for generating and formatting UUIDs as
// specified in RFC 4122.
package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// UUID represents a universally unique identifier made up of 128 bits.
//
// For more information see: https://en.wikipedia.org/wiki/Universally_unique_identifier
type UUID [16]byte

// NewV3 uses the provided namespace and name to generate and return a new v3
// UUID using MD5 hashing, as per RFC 4122.
func NewV3(namespace UUID, name []byte) UUID {
	return usingHash(md5.New(), namespace, name, 3)
}

// NewV4 generates and returns a new v4 UUID using random bytes, as per RFC
// 4122. If an error occurs while reading from "crypto/rand", it is returned.
func NewV4() (UUID, error) {
	var u UUID
	if _, err := io.ReadFull(rand.Reader, u[:]); err != nil {
		return u, err
	}
	setVersion(&u, 4)
	setVariant(&u)
	return u, nil
}

// NewV5 uses the provided namespace and name to generate and return a new v5
// UUID using SHA1 hashing, as per RFC 4122.
func NewV5(namespace UUID, name []byte) UUID {
	return usingHash(sha1.New(), namespace, name, 5)
}

const dash = '-'

// Format returns the hexadecimal format of the UUID as an array of 36 bytes.
//
// Example: 9e754ef6-8dd9-5903-af43-7aea99bfb1fe
func (u UUID) Format() [36]byte {
	var buf [36]byte
	u.format(buf[:])
	return buf
}

func (u UUID) format(buf []byte) {
	hex.Encode(buf[0:], u[:4])
	buf[8] = dash
	hex.Encode(buf[9:], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])
}

// Bytes returns the hexadecimal format of the UUID as a slice of 36 bytes.
//
// Example: 9e754ef6-8dd9-5903-af43-7aea99bfb1fe
func (u UUID) Bytes() []byte {
	b := u.Format()
	return b[:]
}

// String returns the human-readable, hexadecimal format of the UUID as a
// string with a length of 36 bytes.
//
// Example: 9e754ef6-8dd9-5903-af43-7aea99bfb1fe
func (u UUID) String() string {
	return string(u.Bytes())
}

// MarshalBinary implements the BinaryMarshaler interface. It returns a byte
// slice representing the 16 byte binary representation of the UUID.
func (u UUID) MarshalBinary() ([]byte, error) {
	b := u
	return b[:], nil
}

// MarshalJSON implements the MarshalJSON interface. It returns a byte slice
// representing the JSON string of a 36 byte hexadecimal representation of the
// UUID.
func (u UUID) MarshalJSON() ([]byte, error) {
	var b [38]byte
	b[0] = '"'
	u.format(b[1:])
	b[37] = '"'
	return b[:], nil
}

// MarshalText implements the MarshalText interface. It returns a byte slice
// representing the 36 byte hexadecimal representation of the UUID.
func (u UUID) MarshalText() ([]byte, error) {
	return u.Bytes(), nil
}

// Value implements the SQL Driver Valuer interface. It returns a formatted byte
// slice representation of the UUID.
func (u UUID) Value() (driver.Value, error) {
	return u.Bytes(), nil
}

// Version returns the version number of the UUID, as specified in RFC 4122.
func (u UUID) Version() int {
	return int(u[6] >> 4)
}

// usingHash returns a new UUID using the provided hash function, namespace
// UUID, name byte slice, and version number.
func usingHash(h hash.Hash, namespace UUID, name []byte, version byte) UUID {
	var u UUID
	_, _ = h.Write(namespace[:])
	_, _ = h.Write(name)
	copy(u[:], h.Sum(nil))
	setVersion(&u, version)
	setVariant(&u)
	return u
}

// setVersion sets the appropriate version bits in the provided UUID pointed to
// by u.
func setVersion(u *UUID, v byte) {
	u[6] = u[6]&0x0f | v<<4
}

// setVariant sets the variant bits to '10' in the provided UUID pointed to by
// u.
func setVariant(u *UUID) {
	u[8] = u[8]&0x3f | 0x80
}

// ErrInvalidUUID represents the error returned during parsing when the provided
// bytes do not represent a valid UUID.
var ErrInvalidUUID = errors.New("uuid: invalid bytes provided")

// Parse parses the provided UUID bytes, returning the UUID or any error
// encountered. The following formats are provided:
//   16 byte raw, binary UUID
//   32 byte hexadecimal formatted UUID without dashes e.g. 9e754ef68dd94903af437aea99bfb1fe
//   36 byte hexadecimal formatted UUID e.g "9e754ef6-8dd9-4903-af43-7aea99bfb1fe"
func Parse(b []byte) (UUID, error) {
	switch len(b) {
	case 16:
		var u UUID
		copy(u[:], b)
		return u, nil
	case 32:
		var u UUID
		_, err := hex.Decode(u[:], b)
		if err != nil {
			return u, ErrInvalidUUID
		}
		return u, nil
	case 36:
		return parseFormatted(b)
	default:
		return UUID{}, ErrInvalidUUID
	}
}

// ParseString parses the provided UUID string using the same rules as Parse.
func ParseString(s string) (UUID, error) {
	return Parse([]byte(s))
}

var uuidHexLengths = [5]int{8, 4, 4, 4, 12}

// parses returns the parsed 36-byte string UUID into a 16-byte UUID.
func parseFormatted(b []byte) (UUID, error) {
	var u UUID
	var iu, ib int
	for idx, cnt := range uuidHexLengths {
		n, err := hex.Decode(u[iu:], b[ib:ib+cnt])
		if err != nil {
			return u, ErrInvalidUUID
		}
		if idx < 4 && b[ib+cnt] != dash {
			return u, ErrInvalidUUID
		}
		iu += n
		ib += cnt + 1
	}
	return u, nil
}
