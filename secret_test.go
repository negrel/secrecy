package secrecy

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"runtime"
	"testing"
)

func TestSecret(t *testing.T) {
	t.Run("fmt.Sprintf", func(t *testing.T) {
		secret := NewSecret("mysecret")
		str := fmt.Sprintf("%+v", secret)
		if str != "Secret[string](******)" {
			t.Fatal("fmt.Sprintf leak secret:", str)
		}
	})

	t.Run("json.Marshal", func(t *testing.T) {
		secret := NewSecret("mysecret")
		bytes, err := json.Marshal(secret)
		if err != nil {
			t.Fatal(err)
		}
		str := string(bytes)

		if str != "{}" {
			t.Fatal("json.Marshal leak secret:", str)
		}
	})

	t.Run("xml.Marshal", func(t *testing.T) {
		secret := NewSecret("mysecret")
		bytes, err := xml.Marshal(secret)
		if err != nil {
			t.Fatal(err)
		}
		str := string(bytes)

		if str != `<Secret></Secret>` {
			t.Fatal("json.Marshal leak secret:", str)
		}
	})
	runtime.GC()
}
