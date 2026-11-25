package goption_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/manicar2093/goption"
)

func BenchmarkOptional_Get(b *testing.B) {
	n := goption.Some(42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		val, ok := n.Get()
		if !ok || val != 42 {
			b.Fatal("unexpected value")
		}
	}
}

func BenchmarkZero_Get(b *testing.B) {
	z := goption.SomeZero(42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		val, ok := z.Get()
		if !ok || val != 42 {
			b.Fatal("unexpected value")
		}
	}
}

func BenchmarkPointer_Get(b *testing.B) {
	v := 42
	p := &v
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if p == nil || *p != 42 {
			b.Fatal("unexpected value")
		}
	}
}

func BenchmarkSQLNull_Get(b *testing.B) {
	sn := sql.Null[int]{V: 42, Valid: true}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !sn.Valid || sn.V != 42 {
			b.Fatal("unexpected value")
		}
	}
}

func BenchmarkOptional_JSONMarshal(b *testing.B) {
	n := goption.Some("benchmark_string")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(n)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOptional_JSONUnmarshal(b *testing.B) {
	data := []byte(`"benchmark_string"`)
	var n goption.Optional[string]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &n)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkZero_JSONMarshal(b *testing.B) {
	z := goption.SomeZero("benchmark_string")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(z)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkZero_JSONUnmarshal(b *testing.B) {
	data := []byte(`"benchmark_string"`)
	var z goption.Zero[string]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &z)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOptional_Map(b *testing.B) {
	n := goption.Some(42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapped := goption.Map(n, func(v int) int {
			return v * 2
		})
		if !mapped.IsValid() {
			b.Fatal("unexpected invalid")
		}
	}
}

var benchmarkInt int
var benchmarkBool bool
var benchmarkOptionalInt goption.Optional[int]
var benchmarkZeroInt goption.Zero[int]
var benchmarkOptionalPointer goption.Optional[*int]
var benchmarkZeroPointer goption.Zero[*int]

func BenchmarkConstructors(b *testing.B) {
	for _, tc := range []struct {
		name  string
		value int
	}{{"nonzero", 42}, {"zero", 0}} {
		b.Run("Some/"+tc.name, func(b *testing.B) {
			var result goption.Optional[int]
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				result = goption.Some(tc.value)
			}
			benchmarkOptionalInt = result
			if v, ok := result.Get(); !ok || v != tc.value {
				b.Fatal("unexpected result", v, ok)
			}
		})
		b.Run("SomeZero/"+tc.name, func(b *testing.B) {
			var result goption.Zero[int]
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				result = goption.SomeZero(tc.value)
			}
			benchmarkZeroInt = result
			if v, ok := result.Get(); ok != (tc.value != 0) || v != tc.value {
				b.Fatal("unexpected result", v, ok)
			}
		})
	}
	var input *int
	b.Run("Some/nil", func(b *testing.B) {
		var result goption.Optional[*int]
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result = goption.Some(input)
		}
		benchmarkOptionalPointer = result
		if v, ok := result.Get(); ok || v != nil {
			b.Fatal("unexpected result", v, ok)
		}
	})
	b.Run("SomeZero/nil", func(b *testing.B) {
		var result goption.Zero[*int]
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result = goption.SomeZero(input)
		}
		benchmarkZeroPointer = result
		if v, ok := result.Get(); ok || v != nil {
			b.Fatal("unexpected result", v, ok)
		}
	})
}

func BenchmarkAccessors(b *testing.B) {
	for _, present := range []bool{true, false} {
		name := "absent"
		var o goption.Optional[int]
		var z goption.Zero[int]
		if present {
			name = "present"
			o = goption.Some(42)
			z = goption.SomeZero(42)
		}
		for _, tc := range []struct {
			name   string
			get    func() (int, bool)
			orElse func(int) int
		}{{"Optional", o.Get, o.OrElse}, {"Zero", z.Get, z.OrElse}} {
			b.Run(tc.name+"/Get/"+name, func(b *testing.B) {
				var v int
				var ok bool
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					v, ok = tc.get()
				}
				benchmarkInt, benchmarkBool = v, ok
				want := 0
				if present {
					want = 42
				}
				if v != want || ok != present {
					b.Fatal("unexpected result", v, ok)
				}
			})
			b.Run(tc.name+"/OrElse/"+name, func(b *testing.B) {
				var v int
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					v = tc.orElse(7)
				}
				benchmarkInt = v
				want := 7
				if present {
					want = 42
				}
				if v != want {
					b.Fatal("unexpected result", v)
				}
			})
		}
	}
}

var benchmarkSQLValue driver.Value
var benchmarkSQLNull sql.Null[int64]
var benchmarkOptionalInt64 goption.Optional[int64]
var benchmarkZeroInt64 goption.Zero[int64]
var benchmarkJSON []byte
var benchmarkPatch any

func BenchmarkSQLScan(b *testing.B) {
	for _, tc := range []struct {
		name  string
		src   any
		want  int64
		valid bool
	}{
		{"nonzero", int64(42), 42, true},
		{"zero", int64(0), 0, true},
		{"null", nil, 0, false},
		{"conversion", int32(42), 42, true},
	} {
		b.Run("Optional/"+tc.name, func(b *testing.B) {
			var result goption.Optional[int64]
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = result.Scan(tc.src)
			}
			benchmarkOptionalInt64 = result
			v, ok := result.Get()
			if err != nil || v != tc.want || ok != tc.valid {
				b.Fatal("unexpected scan result", v, ok, err)
			}
		})
		b.Run("Zero/"+tc.name, func(b *testing.B) {
			var result goption.Zero[int64]
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = result.Scan(tc.src)
			}
			benchmarkZeroInt64 = result
			v, ok := result.Get()
			if err != nil || v != tc.want || ok != (tc.valid && tc.want != 0) {
				b.Fatal("unexpected scan result", v, ok, err)
			}
		})
		// sql.Null preserves zero, so it matches Optional, not Zero, for zero.
		b.Run("SQLNull/"+tc.name, func(b *testing.B) {
			var result sql.Null[int64]
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = result.Scan(tc.src)
			}
			benchmarkSQLNull = result
			if err != nil || result.V != tc.want || result.Valid != tc.valid {
				b.Fatal("unexpected scan result", result, err)
			}
		})
	}
}

func BenchmarkSQLValue(b *testing.B) {
	for _, tc := range []struct {
		name  string
		value int64
		valid bool
	}{{"nonzero", 42, true}, {"zero", 0, true}, {"null", 0, false}} {
		var optional goption.Optional[int64]
		var zero goption.Zero[int64]
		if tc.valid {
			optional = goption.Some(tc.value)
			zero = goption.SomeZero(tc.value)
		}
		null := sql.Null[int64]{V: tc.value, Valid: tc.valid}
		for _, v := range []struct {
			name   string
			valuer driver.Valuer
			want   driver.Value
		}{
			{"Optional", optional, sqlExpected(tc.value, tc.valid)},
			{"Zero", zero, sqlExpected(tc.value, tc.valid && tc.value != 0)},
			{"SQLNull", null, sqlExpected(tc.value, tc.valid)},
		} {
			b.Run(v.name+"/"+tc.name, func(b *testing.B) {
				var result driver.Value
				var err error
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					result, err = v.valuer.Value()
				}
				benchmarkSQLValue = result
				if err != nil || result != v.want {
					b.Fatal("unexpected value", result, err)
				}
			})
		}
	}
}

func sqlExpected(value int64, valid bool) driver.Value {
	if !valid {
		return nil
	}
	return value
}

func BenchmarkJSONStates(b *testing.B) {
	for _, tc := range []struct {
		name  string
		input string
		value int64
		valid bool
	}{{"nonzero", "42", 42, true}, {"zero", "0", 0, true}, {"null", "null", 0, false}} {
		var o goption.Optional[int64]
		var z goption.Zero[int64]
		if tc.valid {
			o = goption.Some(tc.value)
			z = goption.SomeZero(tc.value)
		}
		zeroJSON := tc.input
		if tc.value == 0 {
			zeroJSON = "null"
		}
		for _, item := range []struct {
			name  string
			value any
			want  string
		}{{"Optional", o, tc.input}, {"Zero", z, zeroJSON}} {
			b.Run("Marshal/"+item.name+"/"+tc.name, func(b *testing.B) {
				var result []byte
				var err error
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					result, err = json.Marshal(item.value)
				}
				benchmarkJSON = result
				if err != nil || string(result) != item.want {
					b.Fatal("unexpected JSON", string(result), err)
				}
			})
		}
		data := []byte(tc.input)
		b.Run("Unmarshal/Optional/"+tc.name, func(b *testing.B) {
			var result goption.Optional[int64]
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = json.Unmarshal(data, &result)
			}
			benchmarkOptionalInt64 = result
			v, ok := result.Get()
			if err != nil || v != tc.value || ok != tc.valid {
				b.Fatal("unexpected JSON result", v, ok, err)
			}
		})
		b.Run("Unmarshal/Zero/"+tc.name, func(b *testing.B) {
			var result goption.Zero[int64]
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = json.Unmarshal(data, &result)
			}
			benchmarkZeroInt64 = result
			v, ok := result.Get()
			if err != nil || v != tc.value || ok != (tc.valid && tc.value != 0) {
				b.Fatal("unexpected JSON result", v, ok, err)
			}
		})
	}
}

type benchmarkOptionalPatch struct {
	Value goption.Optional[int64] `json:"value,omitzero"`
}
type benchmarkZeroPatch struct {
	Value goption.Zero[int64] `json:"value,omitzero"`
}
type benchmarkPointerPatch struct {
	Value *int64 `json:"value,omitempty"`
}

func BenchmarkJSONPatch(b *testing.B) {
	// Each decode starts with a fresh destination inside the timed loop: omitted
	// fields otherwise retain prior state. Pointers distinguish zero from null,
	// but, unlike a separate presence flag, cannot distinguish null from omitted.
	for _, tc := range []struct {
		name, input string
		value       int64
		valid       bool
	}{{"omitted", `{}`, 0, false}, {"null", `{"value":null}`, 0, false}, {"zero", `{"value":0}`, 0, true}, {"nonzero", `{"value":42}`, 42, true}} {
		data := []byte(tc.input)
		b.Run("Optional/"+tc.name, func(b *testing.B) {
			var result benchmarkOptionalPatch
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				result = benchmarkOptionalPatch{}
				err = json.Unmarshal(data, &result)
			}
			benchmarkPatch = result
			v, ok := result.Value.Get()
			if err != nil || v != tc.value || ok != tc.valid {
				b.Fatal("unexpected patch", v, ok, err)
			}
		})
		b.Run("Zero/"+tc.name, func(b *testing.B) {
			var result benchmarkZeroPatch
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				result = benchmarkZeroPatch{}
				err = json.Unmarshal(data, &result)
			}
			benchmarkPatch = result
			v, ok := result.Value.Get()
			if err != nil || v != tc.value || ok != (tc.valid && tc.value != 0) {
				b.Fatal("unexpected patch", v, ok, err)
			}
		})
		b.Run("Pointer/"+tc.name, func(b *testing.B) {
			var result benchmarkPointerPatch
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				result = benchmarkPointerPatch{}
				err = json.Unmarshal(data, &result)
			}
			benchmarkPatch = result
			if err != nil || (result.Value != nil) != tc.valid || (result.Value != nil && *result.Value != tc.value) {
				b.Fatal("unexpected patch", result, err)
			}
		})
	}
	// Go 1.23 ignores omitzero; Go 1.24+ consults IsZero. Detect the
	// behavior outside the timed loop rather than assuming a toolchain version.
	probe, err := json.Marshal(benchmarkOptionalPatch{})
	if err != nil {
		b.Fatal(err)
	}
	omitsZero := string(probe) == `{}`
	if !omitsZero && string(probe) != `{"value":null}` {
		b.Fatal("unexpected omitzero behavior", string(probe))
	}
	for _, tc := range []struct {
		name  string
		value int64
		valid bool
	}{{"null", 0, false}, {"zero", 0, true}, {"nonzero", 42, true}} {
		var o benchmarkOptionalPatch
		var z benchmarkZeroPatch
		var p benchmarkPointerPatch
		if tc.valid {
			o.Value = goption.Some(tc.value)
			z.Value = goption.SomeZero(tc.value)
			v := tc.value
			p.Value = &v
		}
		absent := `{"value":null}`
		if omitsZero {
			absent = `{}`
		}
		optionalWant, zeroWant, pointerWant := absent, absent, `{}`
		if tc.valid {
			optionalWant, pointerWant = `{"value":0}`, `{"value":0}`
			if tc.value != 0 {
				optionalWant, zeroWant, pointerWant = `{"value":42}`, `{"value":42}`, `{"value":42}`
			}
		}
		for _, item := range []struct {
			name  string
			value any
			want  string
		}{
			{"Optional", o, optionalWant}, {"Zero", z, zeroWant}, {"Pointer", p, pointerWant},
		} {
			b.Run("Marshal/"+item.name+"/"+tc.name, func(b *testing.B) {
				var result []byte
				var err error
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					result, err = json.Marshal(item.value)
				}
				benchmarkJSON = result
				if err != nil || string(result) != item.want {
					b.Fatal("unexpected patch JSON", string(result), err)
				}
			})
		}
	}
}

func benchmarkTextDecode[T comparable](b *testing.B, input string, want T, optionalPresent, zeroPresent, wantError bool) {
	data := []byte(input)
	b.Run("Optional", func(b *testing.B) {
		var result goption.Optional[T]
		var err error
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err = result.UnmarshalText(data)
		}
		benchmarkPatch = result
		value, present := result.Get()
		if (err != nil) != wantError || value != want || present != optionalPresent {
			b.Fatal("unexpected text result", value, present, err)
		}
	})
	b.Run("Zero", func(b *testing.B) {
		var result goption.Zero[T]
		var err error
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err = result.UnmarshalText(data)
		}
		benchmarkPatch = result
		value, present := result.Get()
		if (err != nil) != wantError || value != want || present != zeroPresent {
			b.Fatal("unexpected text result", value, present, err)
		}
	})
}

func BenchmarkTextDecode(b *testing.B) {
	b.Run("string", func(b *testing.B) { benchmarkTextDecode(b, "hello", "hello", true, true, false) })
	b.Run("string_empty", func(b *testing.B) { benchmarkTextDecode(b, "", "", false, false, false) })
	b.Run("integer", func(b *testing.B) { benchmarkTextDecode(b, "42", 42, true, true, false) })
	b.Run("integer_zero", func(b *testing.B) { benchmarkTextDecode(b, "0", 0, true, false, false) })
	b.Run("integer_empty", func(b *testing.B) { benchmarkTextDecode(b, "", 0, false, false, false) })
	b.Run("integer_invalid", func(b *testing.B) { benchmarkTextDecode(b, "bad", 0, false, false, true) })
	b.Run("on", func(b *testing.B) { benchmarkTextDecode(b, "on", true, true, true, false) })
	b.Run("off", func(b *testing.B) { benchmarkTextDecode(b, "off", false, true, false, false) })
	b.Run("bool_empty", func(b *testing.B) { benchmarkTextDecode(b, "", false, false, false, false) })
	b.Run("bool_invalid", func(b *testing.B) { benchmarkTextDecode(b, "bad", false, false, false, true) })
}

func BenchmarkCombinators(b *testing.B) {
	for _, present := range []bool{true, false} {
		name := "absent"
		var o goption.Optional[int]
		var z goption.Zero[int]
		if present {
			name = "present"
			o, z = goption.Some(42), goption.SomeZero(42)
		}
		fallback, zeroFallback := goption.Some(7), goption.SomeZero(7)
		double := func(v int) int { return v * 2 }
		flat := func(v int) goption.Optional[int] { return goption.Some(v * 2) }
		flatZero := func(v int) goption.Zero[int] { return goption.SomeZero(v * 2) }
		positive := func(v int) bool { return v > 0 }
		want := 0
		if present {
			want = 42
		}
		orWant := 7
		if present {
			orWant = 42
		}
		for _, tc := range []struct {
			name  string
			run   func() goption.Optional[int]
			want  int
			valid bool
		}{
			{"Map", func() goption.Optional[int] { return goption.Map(o, double) }, want * 2, present},
			{"FlatMap", func() goption.Optional[int] { return goption.FlatMap(o, flat) }, want * 2, present},
			{"Filter", func() goption.Optional[int] { return goption.Filter(o, positive) }, want, present},
			{"Or", func() goption.Optional[int] { return goption.Or(o, fallback) }, orWant, true},
		} {
			b.Run(tc.name+"/"+name, func(b *testing.B) {
				var result goption.Optional[int]
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					result = tc.run()
				}
				benchmarkOptionalInt = result
				if v, ok := result.Get(); v != tc.want || ok != tc.valid {
					b.Fatal("unexpected result", v, ok)
				}
			})
		}
		for _, tc := range []struct {
			name  string
			run   func() goption.Zero[int]
			want  int
			valid bool
		}{
			{"MapZero", func() goption.Zero[int] { return goption.MapZero(z, double) }, want * 2, present},
			{"FlatMapZero", func() goption.Zero[int] { return goption.FlatMapZero(z, flatZero) }, want * 2, present},
			{"FilterZero", func() goption.Zero[int] { return goption.FilterZero(z, positive) }, want, present},
			{"OrZero", func() goption.Zero[int] { return goption.OrZero(z, zeroFallback) }, orWant, true},
		} {
			b.Run(tc.name+"/"+name, func(b *testing.B) {
				var result goption.Zero[int]
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					result = tc.run()
				}
				benchmarkZeroInt = result
				if v, ok := result.Get(); v != tc.want || ok != tc.valid {
					b.Fatal("unexpected result", v, ok)
				}
			})
		}
	}
	// Select a configured timeout, then a service default, then a hard default.
	for _, tc := range []struct {
		name                      string
		configured, service, want int
	}{
		{"configured", 30, 60, 30}, {"service", 0, 60, 60}, {"default", 0, 0, 120},
	} {
		configured, service, fallback := goption.SomeZero(tc.configured), goption.SomeZero(tc.service), goption.SomeZero(120)
		b.Run("OrZeroChain/"+tc.name, func(b *testing.B) {
			var result goption.Zero[int]
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				result = goption.OrZero(goption.OrZero(configured, service), fallback)
			}
			benchmarkZeroInt = result
			if v, ok := result.Get(); v != tc.want || !ok {
				b.Fatal("unexpected result", v, ok)
			}
		})
	}
}
