package goption_test

import (
	"fmt"
	"strconv"

	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Combinators", func() {

	Context("Optional[T] Combinators", func() {

		Describe("Map", func() {
			It("transforms valid value to another type", func() {
				opt := goption.Some(42)
				result := goption.Map(opt, func(v int) string {
					return fmt.Sprintf("val: %d", v)
				})

				Expect(result.IsValid()).To(BeTrue())
				Expect(result.MustGet()).To(Equal("val: 42"))
			})

			It("preserves zero values during mapping", func() {
				opt := goption.Some(0)
				result := goption.Map(opt, func(v int) string {
					return strconv.Itoa(v)
				})

				Expect(result.IsValid()).To(BeTrue())
				Expect(result.MustGet()).To(Equal("0"))
			})

			It("returns invalid Optional for invalid input", func() {
				var opt goption.Optional[int]
				result := goption.Map(opt, func(v int) string {
					return strconv.Itoa(v)
				})

				Expect(result.IsValid()).To(BeFalse())
			})
		})

		Describe("FlatMap", func() {
			It("chains and returns valid transformed Optional", func() {
				opt := goption.Some("123")
				result := goption.FlatMap(opt, func(s string) goption.Optional[int] {
					i, err := strconv.Atoi(s)
					if err != nil {
						return goption.Optional[int]{}
					}
					return goption.Some(i)
				})

				Expect(result.IsValid()).To(BeTrue())
				Expect(result.MustGet()).To(Equal(123))
			})

			It("chains and returns invalid Optional when transformation produces invalid", func() {
				opt := goption.Some("invalid_number")
				result := goption.FlatMap(opt, func(s string) goption.Optional[int] {
					i, err := strconv.Atoi(s)
					if err != nil {
						return goption.Optional[int]{}
					}
					return goption.Some(i)
				})

				Expect(result.IsValid()).To(BeFalse())
			})

			It("returns invalid Optional for invalid input", func() {
				var opt goption.Optional[string]
				result := goption.FlatMap(opt, func(s string) goption.Optional[int] {
					return goption.Some(10)
				})

				Expect(result.IsValid()).To(BeFalse())
			})
		})

		Describe("Filter", func() {
			It("returns original Optional if predicate is satisfied", func() {
				opt := goption.Some(10)
				filtered := opt.Filter(func(v int) bool {
					return v > 5
				})

				Expect(filtered.IsValid()).To(BeTrue())
				Expect(filtered.MustGet()).To(Equal(10))
			})

			It("returns invalid Optional if predicate is not satisfied", func() {
				opt := goption.Some(3)
				filtered := opt.Filter(func(v int) bool {
					return v > 5
				})

				Expect(filtered.IsValid()).To(BeFalse())
			})

			It("returns invalid Optional if input is invalid", func() {
				var opt goption.Optional[int]
				filtered := goption.Filter(opt, func(v int) bool {
					return true
				})

				Expect(filtered.IsValid()).To(BeFalse())
			})
		})

		Describe("IfPresent", func() {
			It("executes action if valid", func() {
				opt := goption.Some("test")
				var captured string

				opt.IfPresent(func(v string) {
					captured = v
				})

				Expect(captured).To(Equal("test"))
			})

			It("does not execute action if invalid", func() {
				var opt goption.Optional[string]
				var called bool

				goption.IfPresent(opt, func(v string) {
					called = true
				})

				Expect(called).To(BeFalse())
			})
		})

		Describe("IfPresentOrElse", func() {
			It("executes action when valid", func() {
				opt := goption.Some("active")
				var actionVal string
				var emptyCalled bool

				opt.IfPresentOrElse(
					func(v string) { actionVal = v },
					func() { emptyCalled = true },
				)

				Expect(actionVal).To(Equal("active"))
				Expect(emptyCalled).To(BeFalse())
			})

			It("executes emptyAction when invalid", func() {
				var opt goption.Optional[string]
				var actionVal string
				var emptyCalled bool

				goption.IfPresentOrElse(
					opt,
					func(v string) { actionVal = v },
					func() { emptyCalled = true },
				)

				Expect(actionVal).To(BeEmpty())
				Expect(emptyCalled).To(BeTrue())
			})
		})

		Describe("Or", func() {
			It("returns self if valid", func() {
				opt := goption.Some("primary")
				fallback := goption.Some("fallback")

				result := opt.Or(fallback)
				Expect(result.MustGet()).To(Equal("primary"))
			})

			It("returns fallback if invalid", func() {
				var opt goption.Optional[string]
				fallback := goption.Some("fallback")

				result := goption.Or(opt, fallback)
				Expect(result.MustGet()).To(Equal("fallback"))
			})
		})
	})

	Context("Zero[T] Combinators", func() {

		Describe("MapZero", func() {
			It("transforms present value to another type", func() {
				z := goption.SomeZero(100)
				result := goption.MapZero(z, func(v int) string {
					return fmt.Sprintf("num: %d", v)
				})

				Expect(result.IsPresent()).To(BeTrue())
				Expect(result.MustGet()).To(Equal("num: 100"))
			})

			It("returns empty Zero when mapping yields zero value", func() {
				z := goption.SomeZero(100)
				result := goption.MapZero(z, func(v int) int {
					return 0 // returns zero value
				})

				Expect(result.IsPresent()).To(BeFalse())
			})

			It("returns empty Zero for empty input", func() {
				var z goption.Zero[int]
				result := goption.MapZero(z, func(v int) string {
					return strconv.Itoa(v)
				})

				Expect(result.IsPresent()).To(BeFalse())
			})
		})

		Describe("FlatMapZero", func() {
			It("chains and returns present transformed Zero", func() {
				z := goption.SomeZero("50")
				result := goption.FlatMapZero(z, func(s string) goption.Zero[int] {
					i, err := strconv.Atoi(s)
					if err != nil {
						return goption.Zero[int]{}
					}
					return goption.SomeZero(i)
				})

				Expect(result.IsPresent()).To(BeTrue())
				Expect(result.MustGet()).To(Equal(50))
			})

			It("returns empty Zero for empty input", func() {
				var z goption.Zero[string]
				result := goption.FlatMapZero(z, func(s string) goption.Zero[int] {
					return goption.SomeZero(10)
				})

				Expect(result.IsPresent()).To(BeFalse())
			})
		})

		Describe("FilterZero", func() {
			It("returns original Zero if predicate is satisfied", func() {
				z := goption.SomeZero(20)
				filtered := z.Filter(func(v int) bool {
					return v >= 20
				})

				Expect(filtered.IsPresent()).To(BeTrue())
				Expect(filtered.MustGet()).To(Equal(20))
			})

			It("returns empty Zero if predicate is not satisfied", func() {
				z := goption.SomeZero(10)
				filtered := z.Filter(func(v int) bool {
					return v >= 20
				})

				Expect(filtered.IsPresent()).To(BeFalse())
			})
		})

		Describe("IfPresent, IfPresentOrElse, Or for Zero[T]", func() {
			It("IfPresent executes action when present", func() {
				z := goption.SomeZero("running")
				var captured string
				z.IfPresent(func(v string) { captured = v })
				Expect(captured).To(Equal("running"))
			})

			It("IfPresentOrElse executes appropriate branch", func() {
				var z goption.Zero[int]
				var branch string

				z.IfPresentOrElse(
					func(v int) { branch = "present" },
					func() { branch = "empty" },
				)
				Expect(branch).To(Equal("empty"))
			})

			It("Or returns fallback if empty", func() {
				var z goption.Zero[string]
				fallback := goption.SomeZero("backup")

				result := z.Or(fallback)
				Expect(result.MustGet()).To(Equal("backup"))
			})
		})
	})
})
