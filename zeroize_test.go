package secrecy

import (
	"testing"
)

func TestZeroize(t *testing.T) {
	t.Run("Primitive", func(t *testing.T) {
		t.Run("Bool", func(t *testing.T) {
			data := true
			Zeroize(&data)
			if data != false {
				t.Fatal("Zeroize(true) didn't change bool to false")
			}
		})
		t.Run("String", func(t *testing.T) {
			data := "foo"
			Zeroize(&data)
			if data != "" {
				t.Fatal(`Zeroize("foo") didn't change string to ""`)
			}
		})
	})
	t.Run("Slice", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			data := []int{1, 2, 3}
			Zeroize(&data)
			if data != nil {
				t.Fatal("Zeroize([]int{1, 2, 3}) didn't change slice to nil")
			}
		})
		t.Run("*int", func(t *testing.T) {
			arr := [3]int{1, 2, 3}
			data := []*int{&arr[0], &arr[1], &arr[2]}
			Zeroize(&data)
			if data != nil {
				t.Fatal("Zeroize([]*int{1, 2, 3}) didn't change slice to nil")
			}
			if arr != [3]int{0, 0, 0} {
				t.Fatal("Zeroize([]*int{1, 2, 3}) didn't zeroize pointers memory")
			}
		})
		t.Run("**string", func(t *testing.T) {
			arr := [3]string{"1", "2", "3"}
			arrPtr := [3]*string{&arr[0], &arr[1], &arr[2]}
			data := []**string{&arrPtr[0], &arrPtr[1], &arrPtr[2]}
			Zeroize(&data)
			if data != nil {
				t.Fatal(`Zeroize([]**string{"1", "2", "3"}) didn't change slice nil`)
			}
			if arrPtr != [3]*string{} {
				t.Fatal(`Zeroize([]*string{"1", "2", "3"}) didn't change array to []*string{nil, nil, nil}`)
			}
			if arr != [3]string{"", "", ""} {
				t.Fatal(`Zeroize([]string{"", "", ""}) didn't zeroize pointers memory:`, arr)
			}
		})
	})
	t.Run("Array", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			data := [3]int{1, 2, 3}
			Zeroize(&data)
			if data != [3]int{} {
				t.Fatal("Zeroize([]int{1, 2, 3}) didn't change slice to []int{0, 0, 0}")
			}
		})
		t.Run("*int", func(t *testing.T) {
			arr := [3]int{1, 2, 3}
			data := [3]*int{&arr[0], &arr[1], &arr[2]}
			Zeroize(&data)
			if data != [3]*int{} {
				t.Fatal("Zeroize([]*int{1, 2, 3}) didn't change array to []*int{nil, nil, nil}")
			}
			if arr != [3]int{0, 0, 0} {
				t.Fatal("Zeroize([]*int{1, 2, 3}) didn't zeroize pointers memory")
			}
		})
		t.Run("**int", func(t *testing.T) {
			arr := [3]int{1, 2, 3}
			arrPtr := [3]*int{&arr[0], &arr[1], &arr[2]}
			data := [3]**int{&arrPtr[0], &arrPtr[1], &arrPtr[2]}
			Zeroize(&data)
			if data != [3]**int{} {
				t.Fatal(`Zeroize([]**int{1, 2, 3}) didn't change array to []**string{nil, nil, nil}`)
			}
			if arrPtr != [3]*int{} {
				t.Fatal(`Zeroize([]**int{1, 2, 3}) didn't change pointers of pointers to nil`)
			}
			if arr != [3]int{0, 0, 0} {
				t.Fatal(`Zeroize([]**int{1, 2, 3}) didn't zeroize pointers memory`)
			}
		})
	})
	t.Run("Struct", func(t *testing.T) {
		t.Run("WithUnexportedFields", func(t *testing.T) {
			t.Run("ImplementsZeroizer", func(t *testing.T) {
				num := 1
				instance := ZeroizerStruct{ptr: &num, str: "foo"}
				Zeroize(&instance)
				if (instance != ZeroizerStruct{}) {
					t.Fatal(`Zeroize(ZeroizerStruct{ptr: 1, str: "foo"}) didn't change struct to ZeroizerStruct{ptr: nil, str: ""}`)
				}
				if num != 0 {
					t.Fatal(`Zeroize(ZeroizerStruct{ptr: 1, str: "foo"}) didn't zeroize pointers memory`)
				}
			})

			t.Run("DoesNotImplementsZeroizer", func(t *testing.T) {
				type Struct struct {
					ptr *int
					str string
				}

				num := 1
				instance := Struct{ptr: &num, str: "foo"}
				var recovered error

				defer func() {
					if recovered == nil {
						t.Fatal("Zeroize didn't panic")
					}
					if recovered.Error() != "type Struct contains unexported field and doesn't implements Zeroizer interface" {
						t.Fatal("wrong panic reason:", recovered.Error())
					}
				}()
				defer func() {
					if v := recover(); v != nil {
						recovered = v.(error)
					}
				}()

				Zeroize(&instance)
			})
		})

		t.Run("WithIgnoredUnexportedFields", func(t *testing.T) {
			type Struct struct {
				Ptr *int
				str string `zeroize:"ignore"`
			}

			num := 1
			instance := Struct{Ptr: &num, str: "foo"}
			Zeroize(&instance)
			if (instance != Struct{}) {
				t.Fatal(`Zeroize(Struct{Ptr: 1, str: "foo"}) didn't change struct to Struct{Ptr: nil, str: "foo"}`)
			}
			if num != 0 {
				t.Fatal(`Zeroize(Struct{ptr: 1, str: "foo"}) didn't zeroize pointers memory`)
			}
		})
	})

	t.Run("Map", func(t *testing.T) {
		t.Run("PrimitiveKeyValue", func(t *testing.T) {
			data := map[string]int{"1": 1, "2": 2, "3": 3}
			data2 := data
			Zeroize(&data)
			if data != nil {
				t.Fatal(`Zeroize(map[string]int{"1": 1, "2": 2, "3": 3}) didn't change map to nil`)
			}
			if len(data2) != 0 {
				t.Fatal(`Zeroize(map[string]int{"1": 1, "2": 2, "3": 3}) didn't empty map`)
			}
		})
		t.Run("StructKeyValue", func(t *testing.T) {
			one := 1
			zero := ZeroizerStruct{ptr: &one, str: "foo"}
			data := map[ZeroizerStruct]ZeroizerStruct{zero: zero}
			data2 := data
			Zeroize(&data)
			if data != nil {
				t.Fatal(`Zeroize(map[ZeroizerStruct]ZeroizerStruct{ZeroizerStruct{ptr: 1, str: "foo"}: ZeroizerStruct{ptr: 1, str: "foo"}}) didn't change map to nil`)
			}
			if len(data2) != 0 {
				t.Fatal(`Zeroize(map[ZeroizerStruct]ZeroizerStruct{ZeroizerStruct{ptr: 1, str: "foo"}: ZeroizerStruct{ptr: 1, str: "foo"}}) didn't empty map`)
			}
			if one != 0 {
				t.Fatal(`Zeroize(map[ZeroizerStruct]ZeroizerStruct{ZeroizerStruct{ptr: 1, str: "foo"}: ZeroizerStruct{ptr: 1, str: "foo"}}) didn't zeroize pointers memory`)
			}
			// Don't check for zero struct as it copied into the map.
		})
	})
}

type ZeroizerStruct struct {
	ptr *int
	str string
}

func (zs ZeroizerStruct) Zeroize() {
	Zeroize(zs.ptr)
}
