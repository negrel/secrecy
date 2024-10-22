package secrecy

import (
	"encoding/json"
	"testing"
)

func TestSerializableSecret(t *testing.T) {
	t.Run("json.Marshal", func(t *testing.T) {
		secret := NewSerializableSecret(NewSecretString([]byte("mysecret")))

		bytes, err := json.Marshal(secret)
		if err != nil {
			t.Fatal(err)
		}

		str := string(bytes)
		if str != `"mysecret"` {
			t.Fatal("json.Marshal(SerializableSecret) didn't serialize the secret:", str)
		}
	})
}
