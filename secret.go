package secrecy

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

// NewSecret wraps given secret and returns *Secret[T]. Wrapped value and inner
// values will be entirely wiped from memory when garbage collected, this
// includes maps, slices pointers, etc. Returned secret owns wrapped value and
// inner value, you must not share underlying data.
func NewSecret[T any](value T) *Secret[T] {
	secret := new(Secret[T])
	secret.value = value
	runtime.SetFinalizer(secret, Zeroize)
	return secret
}

// Secret is a wrapper type for values that contains secrets, which attempts to
// limit accidental exposure and ensure secrets are wiped from memory when
// garbage collected. (e.g. passwords, cryptographic keys, access tokens or
// other credentials)
//
// Prefer SecretString over Secret[string] as go string are immutable and can't
// be wiped from memory.
type Secret[T any] struct {
	value T
}

// ExposeSecret returns a copy of inner secret.
// Don't store returned value and prefer passing Secret itself if needed.
func (s *Secret[T]) ExposeSecret() T {
	return s.value
}

// String implements fmt.Stringer.
func (s *Secret[T]) String() string {
	typeName := reflect.TypeOf(s.value).String()
	return fmt.Sprintf("Secret[%v](******)", typeName)
}

// Disable zeroize on garbage collection for this secret.
func (s *Secret[T]) DisableZeroize() {
	runtime.SetFinalizer(s, nil)
}

// Enable zeroize on garbage collection for this secret.
// By default zeroize is enabled, you don't need to call this function if didn't
// call WithoutZeroize before.
func (s *Secret[T]) EnableZeroize() {
	runtime.SetFinalizer(s, Zeroize)
}

// Zeroize implements Zeroizer.
func (s *Secret[T]) Zeroize() {
	Zeroize(s.value)
}

// NewSecretString wraps given secret and returns SecretString.
func NewSecretString(secret []byte) SecretString {
	return SecretString{NewSecret(secret)}
}

// SecretString is a wrapper around Secret[[]byte] that expose its secret
// as string.
type SecretString struct {
	*Secret[[]byte]
}

// ExposeSecret exposes underlying secret as a string using unsafe.
// Don't store returned value and prefer passing SecretString itself if needed.
func (ss SecretString) ExposeSecret() string {
	bytes := ss.Secret.ExposeSecret()
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}
