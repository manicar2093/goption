package goption_test

import (
	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type sampleStruct struct {
	Name string
	Age  int
}

var _ = Describe("Optional", func() {

	Describe("Some - Non-zero primitive, struct, and pointer types", func() {
		DescribeTable("creates a valid Optional for all non-zero types", func(val any) {
			got := goption.Some(val)

			Expect(got.IsValid()).To(BeTrue())
			Expect(got.IsPresent()).To(BeTrue())
			Expect(got.IsZero()).To(BeFalse())

			v, ok := got.Get()
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(val))
			Expect(got.OrElse(nil)).To(Equal(val))
			Expect(got.MustGet()).To(Equal(val))
		},
			Entry("bool true", true),
			Entry("string", "hello"),
			Entry("int", int(6)),
			Entry("int8", int8(6)),
			Entry("int16", int16(6)),
			Entry("int32", int32(6)),
			Entry("int64", int64(6)),
			Entry("uint", uint(6)),
			Entry("uint8", uint8(6)),
			Entry("uint16", uint16(6)),
			Entry("uint32", uint32(6)),
			Entry("uint64", uint64(6)),
			Entry("float32", float32(6.5)),
			Entry("float64", float64(6.5)),
			Entry("complex64", complex64(6+2i)),
			Entry("complex128", complex128(6+2i)),
			Entry("struct", sampleStruct{Name: "name", Age: 30}),
			Entry("empty struct value", sampleStruct{}),
			Entry("non-nil empty slice", []string{}),
			Entry("non-nil filled slice", []string{"a", "b"}),
			Entry("non-nil empty map", map[string]int{}),
			Entry("pointer bool", toPointer(true)),
			Entry("pointer string", toPointer("hello")),
			Entry("pointer int", toPointer(int(6))),
			Entry("pointer int8", toPointer(int8(6))),
			Entry("pointer int16", toPointer(int16(6))),
			Entry("pointer int32", toPointer(int32(6))),
			Entry("pointer int64", toPointer(int64(6))),
			Entry("pointer uint", toPointer(uint(6))),
			Entry("pointer uint8", toPointer(uint8(6))),
			Entry("pointer uint16", toPointer(uint16(6))),
			Entry("pointer uint32", toPointer(uint32(6))),
			Entry("pointer uint64", toPointer(uint64(6))),
			Entry("pointer float32", toPointer(float32(6.5))),
			Entry("pointer float64", toPointer(float64(6.5))),
			Entry("pointer complex64", toPointer(complex64(6+2i))),
			Entry("pointer complex128", toPointer(complex128(6+2i))),
			Entry("pointer struct", toPointer(sampleStruct{Name: "name", Age: 30})),
		)
	})

	Describe("Some - Zero values preservation", func() {
		DescribeTable("treats Go zero values as valid present values", func(val any) {
			got := goption.Some(val)

			Expect(got.IsValid()).To(BeTrue())
			Expect(got.IsPresent()).To(BeTrue())
			Expect(got.IsZero()).To(BeFalse())

			v, ok := got.Get()
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(val))
			Expect(got.MustGet()).To(Equal(val))
		},
			Entry("zero bool", false),
			Entry("zero string", ""),
			Entry("zero int", int(0)),
			Entry("zero int8", int8(0)),
			Entry("zero int16", int16(0)),
			Entry("zero int32", int32(0)),
			Entry("zero int64", int64(0)),
			Entry("zero uint", uint(0)),
			Entry("zero uint8", uint8(0)),
			Entry("zero uint16", uint16(0)),
			Entry("zero uint32", uint32(0)),
			Entry("zero uint64", uint64(0)),
			Entry("zero float32", float32(0.0)),
			Entry("zero float64", float64(0.0)),
			Entry("zero complex64", complex64(0)),
			Entry("zero complex128", complex128(0)),
			Entry("empty struct", sampleStruct{}),
			Entry("pointer to zero int", toPointer(0)),
			Entry("pointer to zero string", toPointer("")),
			Entry("pointer to zero bool", toPointer(false)),
			Entry("pointer to empty struct", toPointer(sampleStruct{})),
		)
	})

	Describe("Some - Nil values recognition", func() {
		DescribeTable("treats nil pointers, slices, maps, channels, and interfaces as invalid/null", func(val any) {
			got := goption.Some(val)

			Expect(got.IsValid()).To(BeFalse())
			Expect(got.IsPresent()).To(BeFalse())
			Expect(got.IsZero()).To(BeTrue())

			_, ok := got.Get()
			Expect(ok).To(BeFalse())

			Expect(func() {
				got.MustGet()
			}).To(PanicWith(goption.ErrNoSuchElement))
		},
			Entry("nil interface", nil),
			Entry("nil *bool", (*bool)(nil)),
			Entry("nil *string", (*string)(nil)),
			Entry("nil *int", (*int)(nil)),
			Entry("nil *int8", (*int8)(nil)),
			Entry("nil *int16", (*int16)(nil)),
			Entry("nil *int32", (*int32)(nil)),
			Entry("nil *int64", (*int64)(nil)),
			Entry("nil *uint", (*uint)(nil)),
			Entry("nil *uint8", (*uint8)(nil)),
			Entry("nil *uint16", (*uint16)(nil)),
			Entry("nil *uint32", (*uint32)(nil)),
			Entry("nil *uint64", (*uint64)(nil)),
			Entry("nil *float32", (*float32)(nil)),
			Entry("nil *float64", (*float64)(nil)),
			Entry("nil *complex64", (*complex64)(nil)),
			Entry("nil *complex128", (*complex128)(nil)),
			Entry("nil *struct", (*sampleStruct)(nil)),
			Entry("nil slice of strings", ([]string)(nil)),
			Entry("nil slice of pointers", ([]*string)(nil)),
			Entry("nil map", (map[string]int)(nil)),
			Entry("nil channel", (chan int)(nil)),
		)
	})

	Describe("Uninitialized / Empty Optional[T]", func() {
		It("is invalid by default with correct accessor behavior", func() {
			var opt goption.Optional[string]

			Expect(opt.IsValid()).To(BeFalse())
			Expect(opt.IsPresent()).To(BeFalse())
			Expect(opt.IsZero()).To(BeTrue())

			v, ok := opt.Get()
			Expect(ok).To(BeFalse())
			Expect(v).To(Equal(""))

			Expect(opt.OrElse("fallback")).To(Equal("fallback"))
			Expect(opt.OrElseZero()).To(Equal(""))

			Expect(func() {
				opt.MustGet()
			}).To(PanicWith(goption.ErrNoSuchElement))
		})

		It("OrElse and OrElseZero return contained value when valid", func() {
			opt := goption.Some("active")

			Expect(opt.OrElse("fallback")).To(Equal("active"))
			Expect(opt.OrElseZero()).To(Equal("active"))
		})
	})
})
