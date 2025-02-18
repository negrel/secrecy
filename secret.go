package secrecy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

// SecretLeakedMarker is a marker string returned by Secret[T] when accidentally
// leaking the secret. This value is returned when you try to fmt.Println(secret)
// or serialize it to JSON/XML. This enable you to easily monitor your logs for
// secret leak.
var SecretLeakedMarker = "<!SECRET_LEAKED!>"

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
// This function returns the [SecretLeakedMarker] marker string so it can be
// easily searched.
func (s Secret[T]) String() string {
	return SecretLeakedMarker
}

// Format implements fmt.GoStringer.
// This function return the [SecretLeakedMarker] marker string and appends
// "Secret[T](******)" as a suffix.
func (s Secret[T]) GoString() string {
	typeName := reflect.TypeOf(s.value).String()
	return fmt.Sprintf(SecretLeakedMarker+" Secret[%v](******)", typeName)
}

// MarshalText implements encoding.TextMarshaler
// This function returns the [SecretLeakedMarker] marker byte slice so it can be
// easily searched.
func (s Secret[T]) MarshalText() ([]byte, error) {
	return []byte(SecretLeakedMarker), nil
}

func (s *Secret[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.value)
}

// Disable zeroize on garbage collection for this secret.
func (s *Secret[T]) DisableZeroize() {
	runtime.SetFinalizer(s, nil)
}

// Enable zeroize on garbage collection for this secret.
// By default zeroize is enabled, you don't need to call this function if didn't
// call DisableZeroize before.
func (s *Secret[T]) EnableZeroize() {
	runtime.SetFinalizer(s, Zeroize)
}

// Zeroize implements Zeroizer.
func (s Secret[T]) Zeroize() {
	Zeroize(s.value)
}

// NewSecretString wraps given secret and returns SecretString.
// This function takes a byte slice as parameter so the secret can be wiped
// on garbage collection. Use UnsafeStringToBytes to convert your string to a
// byte slice if needed.
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
	return UnsafeBytesToString(bytes)
}

// UnsafeStringToBytes converts given string into a bytes slice using unsafe.
// This function doesn't perform any copy so mutating the byte slice will mutate
// the string.
func UnsafeStringToBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

// UnsafeStringToBytes converts given string into a bytes slice using unsafe.
// This function doesn't perform any copy so mutating the byte slice will mutate
// the returned string.
func UnsafeBytesToString(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}
