package trimmer_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/soundtrackyourbrand/trimmer"
)

func TestMail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trimmer Suite")
}

var _ = Describe("trimmer", func() {
	Context(".TrimStrings(..)", func() {
		It("should trim string fields", func() {
			type stringalias string
			type test struct {
				Field1 string
				Field2 *string
				Field3 **string
				Field4 stringalias
			}

			s := " ptr field        "
			ptr := &s
			t := test{
				Field1: "  field1   ",
				Field2: ptr,
				Field3: &ptr,
				Field4: "  alias ",
			}

			err := TrimStrings(&t)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(t.Field1).To(Equal("field1"))
			Expect(*t.Field2).To(Equal("ptr field"))
			Expect(**t.Field3).To(Equal("ptr field"))
			Expect(t.Field4).To(Equal(stringalias("alias")))
		})

		It("should not trim anything to other field types", func() {

			type test struct {
				Field1 []byte
			}

			t := test{
				Field1: []byte(" byte! "),
			}

			err := TrimStrings(&t)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(t.Field1).To(Equal([]byte(" byte! ")))
		})

		It("should return error if not a pointer", func() {
			type test struct{}
			t := test{}
			err := TrimStrings(t)
			Expect(err).To(Equal(ErrInvalidType))
		})

		It("should return an error if pointer to a non-struct", func() {
			t := []string{" a", " b", " c"}
			err := TrimStrings(&t)
			Expect(err).To(Equal(ErrInvalidType))
		})

		It("shoud not trim strings with trim tag set to false", func() {
			type test struct {
				Field1 string  `trim:"false"`
				Field2 *string `trim:"false"`
			}

			s := " ptr "
			t := test{
				Field1: " not trimmed ",
				Field2: &s,
			}

			err := TrimStrings(&t)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(t.Field1).To(Equal(" not trimmed "))
			Expect(*t.Field2).To(Equal(" ptr "))
		})

		It("should trim string fields in nested objects", func() {
			type inner struct {
				Field1 string
				Field2 string
			}

			type middle struct {
				Field1 string
				Field2 string
				Inner1 *inner
				Inner2 *inner
				Inner3 inner
			}

			type outer struct {
				Field   string
				Middle1 middle
				Middle2 middle
			}

			t := outer{
				Field: " first ",
				Middle1: middle{
					Field1: " field1 ",
				},
				Middle2: middle{
					Field1: " field1 ",
					Field2: " field2 ",
					Inner1: &inner{
						Field1: " field1 ",
						Field2: " field2 ",
					},
					Inner2: &inner{
						Field1: " field1 ",
						Field2: " field2 ",
					},
					Inner3: inner{
						Field1: " field1 ",
						Field2: " field2 ",
					},
				},
			}

			err := TrimStrings(&t)

			Expect(err).ShouldNot(HaveOccurred())

			Expect(t.Field).To(Equal("first"))
			Expect(t.Middle1.Field1).To(Equal("field1"))
			Expect(t.Middle2.Field1).To(Equal("field1"))
			Expect(t.Middle2.Field2).To(Equal("field2"))
			Expect(t.Middle2.Inner1.Field1).To(Equal("field1"))
			Expect(t.Middle2.Inner1.Field2).To(Equal("field2"))
			Expect(t.Middle2.Inner2.Field1).To(Equal("field1"))
			Expect(t.Middle2.Inner2.Field2).To(Equal("field2"))
			Expect(t.Middle2.Inner3.Field1).To(Equal("field1"))
			Expect(t.Middle2.Inner3.Field2).To(Equal("field2"))
		})

		It("should not fail when there are wierd time types on the struct", func() {
			type specialTime struct{ time.Time }

			type test struct {
				Time specialTime
			}

			t := test{
				Time: specialTime{time.Now()},
			}

			Expect(TrimStrings(&t)).To(Succeed())
		})
	})
})
