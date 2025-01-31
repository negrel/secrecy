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

	t.Run("json.Unmarshal", func(t *testing.T) {
		bytes := []byte(`"mysecret"`)
		var secret SerializableSecret[string, *Secret[string]]
		err := json.Unmarshal(bytes, &secret)
		if err != nil {
			t.Fatal(err)
		}
		str := secret.ExposeSecret()

		if str != "mysecret" {
			t.Fatal("json.Unmarshal invalid value:", str)
		}
	})
}
