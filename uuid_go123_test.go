//go:build go1.23

package uuid

import (
	"crypto/rand"
	"io"
	mrand "math/rand/v2"
	"testing"
	"time"
)

func BenchmarkNewV4ChaCha8(b *testing.B) {
	var seed [32]byte
	_, err := io.ReadFull(rand.Reader, seed[:])
	if err != nil {
		b.Fatalf("unable to read from crypto/rand: %s", err.Error())
	}
	r := mrand.NewChaCha8(seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewV4FromRand(r)
	}
}

func BenchmarkNewV7ChaCha8(b *testing.B) {
	var seed [32]byte
	_, err := io.ReadFull(rand.Reader, seed[:])
	if err != nil {
		b.Fatalf("unable to read from crypto/rand: %s", err.Error())
	}
	now := time.Unix(1000, 0)
	r := mrand.NewChaCha8(seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewV7FromRand(now, r)
	}
}
