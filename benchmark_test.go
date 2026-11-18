package goption_test

import (
	"database/sql"
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
