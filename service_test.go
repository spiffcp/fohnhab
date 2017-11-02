package fohnhab_test

import (
	"context"
	"encoding/base64"

	"github.com/spiffcp/fohnhab"
)

var _ = Describe("Service", func() {
	Describe("Service Endpoints", func() {

		Describe("GenerateKey", func() {
			var (
				s   fohnhab.Service
				t   string
				err error
				c   fohnhab.GenerateKeyRequest
				ctx context.Context
			)
			Context("When called with correct arguments", func() {
				BeforeEach(func() {
					c.Kind = "aes-256"
					s = fohnhab.NewService()
					t, err = s.GenerateKey(ctx, c)
				})
				It("Should not error", func() {
					Expect(err).To(Not(HaveOccurred()))
				})
				It("Should return a 256 bit key for the user as a string", func() {
					data, _ := base64.StdEncoding.DecodeString(t)
					Expect(len(data)).To(Equal(32))
				})
			})

			Context("When called with an incorrect argument", func() {
				BeforeEach(func() {
					c.Kind = "aes-2555"
					s = fohnhab.NewService()
					t, err = s.GenerateKey(ctx, c)
				})
				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
