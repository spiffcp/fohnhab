package fohnhab_test

import (
	"context"

	"github.com/spiffcp/fohnhab"
)

var _ = Describe("Service", func() {
	Describe("Service Endpoints", func() {

		Describe("GenerateKey", func() {
			var (
				t   string
				err error
				c   fohnhab.GenerateKeyRequest
				ctx context.Context
			)
			Context("When called with a valid key type argument", func() {
				BeforeEach(func() {
					c.Kind = "aes-256"
					t, err = s.GenerateKey(ctx, c)
				})
				It("Should not error", func() {
					Expect(err).To(Not(HaveOccurred()))
				})
				It("Should return a 256 bit key for the user as a string", func() {
					Expect(len(t)).To(Equal(44))
				})
			})

			Context("When called with an incorrect argument", func() {
				BeforeEach(func() {
					c.Kind = "aes-2555"
					t, err = s.GenerateKey(ctx, c)
				})
				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("Galois/Counter Encryption (GCM)", func() {
			var (
				et  string
				err error
				req fohnhab.GCMERequest
				ctx context.Context
			)
			Context("When called with a valid key and plaintext to encrypt", func() {
				req.Key = "lwL3V4RoI7vnoh8TbgQ16mr+M60cdLrPhpHJ923Oldw="
				req.PlainText = "This is the best message ever"
				BeforeEach(func() {
					et, err = s.GCME(ctx, req)
				})
				It("Should return the encrypted text", func() {
					Expect(et).NotTo(Equal(""))
				})
			})
		})

		Describe("Galois/Counter Decryption", func() {
			var (
				dt  string
				err error
				req fohnhab.GCMDRequest
				ctx context.Context
			)
			Context("When called with a valid key and cyphertext to decrypt", func() {
				req.Key = "lwL3V4RoI7vnoh8TbgQ16mr+M60cdLrPhpHJ923Oldw="
				req.CipherText = "wFRV/VSW4GoeDD2bVC6IGphCv93wgyBxIhRSr61S0Aq0RunXO8VOYeXOyIVCWYhSObwrLRlaOINT"
				BeforeEach(func() {
					dt, err = s.GCMD(ctx, req)
					if err != nil {
						Fail(err.Error())
					}
				})
				It("Should return the decrpyted text", func() {
					Expect(dt).To(Equal("This is the best message ever"))
				})
			})
		})
	})
})
