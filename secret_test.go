package secrecy

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"testing"
)

func TestSecret(t *testing.T) {
	t.Run("fmt.Sprintf(%+v)", func(t *testing.T) {
		secret := NewSecret("mysecret")
		str := fmt.Sprintf("%+v", secret)
		if str != "<!SECRET_LEAKED!>" {
			t.Fatal("fmt.Sprintf(% +v) leak secret:", str)
		}

		// Should work on a non pointer receiver too
		str = fmt.Sprintf("%+v", *secret)
		if str != "<!SECRET_LEAKED!>" {
			t.Fatal("fmt.Sprintf(% +v) leak secret:", str)
		}
	})

	t.Run("fmt.Sprintf(%#v)", func(t *testing.T) {
		secret := NewSecret("mysecret")
		str := fmt.Sprintf("%#v", secret)
		if str != "<!SECRET_LEAKED!> Secret[string](******)" {
			t.Fatal("fmt.Sprintf(% #v) leak secret:", str)
		}

		// Should work on a non pointer receiver too
		str = fmt.Sprintf("%#v", *secret)
		if str != "<!SECRET_LEAKED!> Secret[string](******)" {
			t.Fatal("fmt.Sprintf(% #v) leak secret:", str)
		}
	})

	t.Run("slog.Log", func(t *testing.T) {
		t.Run("WithTextHandler", func(t *testing.T) {
			buf := bytes.Buffer{}
			logger := slog.New(slog.NewTextHandler(&buf, nil))

			secret := NewSecret("mysecret")
			logger.Info("", "secret", secret, "derefSecret", *secret)
			str := buf.String()
			if strings.Contains(str, "mysecret") {
				t.Fatal("fmt.Sprintf leak secret:", str)
			}
		})
		t.Run("WithJsonHandler", func(t *testing.T) {
			buf := bytes.Buffer{}
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			secret := NewSecret("mysecret")
			logger.Info("", "secret", secret, "derefSecret", *secret)
			str := buf.String()
			if strings.Contains(str, "mysecret") {
				t.Fatal("fmt.Sprintf leak secret:", str)
			}
		})
	})

	t.Run("json.Marshal", func(t *testing.T) {
		secret := NewSecret("mysecret")
		bytes, err := json.Marshal(secret)
		if err != nil {
			t.Fatal(err)
		}
		str := string(bytes)

		if str != `"\u003c!SECRET_LEAKED!\u003e"` {
			t.Fatal("json.Marshal leak secret:", str)
		}

		// Should work on a non pointer receiver too
		bytes, err = json.Marshal(*secret)
		if err != nil {
			t.Fatal(err)
		}
		str = string(bytes)

		if str != `"\u003c!SECRET_LEAKED!\u003e"` {
			t.Fatal("json.Marshal leak secret:", str)
		}
	})

	t.Run("json.Unmarshal", func(t *testing.T) {
		bytes := []byte(`"mysecret"`)
		var secret Secret[string]
		err := json.Unmarshal(bytes, &secret)
		if err != nil {
			t.Fatal(err)
		}
		str := secret.value

		if str != "mysecret" {
			t.Fatal("json.Unmarshal invalid value:", str)
		}
	})

	t.Run("xml.Marshal", func(t *testing.T) {
		secret := NewSecret("mysecret")
		bytes, err := xml.Marshal(secret)
		if err != nil {
			t.Fatal(err)
		}
		str := string(bytes)

		if str != `<Secret[string]>&lt;!SECRET_LEAKED!&gt;</Secret[string]>` {
			t.Fatal("xml.Marshal leak secret:", str)
		}

		// Should work on a non pointer receiver too
		bytes, err = xml.Marshal(*secret)
		if err != nil {
			t.Fatal(err)
		}
		str = string(bytes)

		if str != `<Secret[string]>&lt;!SECRET_LEAKED!&gt;</Secret[string]>` {
			t.Fatal("xml.Marshal leak secret:", str)
		}
	})

	// Run garbage collector to trigger finalizer and zeroize memory.
	runtime.GC()
}
