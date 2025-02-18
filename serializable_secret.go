package secrecy

import "encoding/json"

// SecretExposer define any secret wrapper types that can expose it's underlying
// secret of type T.
// Don't store returned value and prefer passing SecretExposer itself if needed.
type SecretExposer[T any] interface {
	ExposeSecret() T
}

// NewSerializableSecret wraps secret and return a SerializableSecret that implements
// json.Marshaler interface.
func NewSerializableSecret[S any, T SecretExposer[S]](secret T) SerializableSecret[S, T] {
	return SerializableSecret[S, T]{secret}
}

// SerializableSecret is a serializable wrapper around a SecretExposer.
type SerializableSecret[S any, T SecretExposer[S]] struct {
	secret T
}

// ExposeSecret implements SecretExposer.
func (ss SerializableSecret[S, T]) ExposeSecret() S {
	return ss.secret.ExposeSecret()
}

// MarshalJSON implements json.Marshaler.
func (ss SerializableSecret[S, T]) MarshalJSON() ([]byte, error) {
	secret := ss.ExposeSecret()
	return json.Marshal(secret)
}

func (ss *SerializableSecret[S, T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &ss.secret)
}
