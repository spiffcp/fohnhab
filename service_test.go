package fohnhab_test

import (
	"github.com/spiffcp/fohnhab"
)

var _ = Describe("Service", func() {
	Describe("Service Endpoints", func() {

		Describe("GenerateKey", func() {
			var (
				s   fohnhab.Service
				t   []byte
				err error
				c   string
				i   string
			)
			Context("When called with correct arguments", func() {
				c = "aes-256"
				BeforeEach(func() {
					s = fohnhab.NewService()
					t, err = s.GenerateKey(c)
				})
				It("Should not error", func() {
					Expect(err).To(Not(HaveOccurred()))
				})
				It("Should return a 256 bit key for the user", func() {
					Expect(len(t)).To(Equal(32))
				})
			})

			Context("When called with an incorrect argument", func() {
				i = "aes-75309"
				BeforeEach(func() {
					s = fohnhab.NewService()
					t, err = s.GenerateKey(i)
				})
				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
