package fohnhab_test

import (
	"context"

	"github.com/spiffcp/fohnhab"
)

var _ = Describe("Service", func() {
	Describe("Service Endpoints", func() {

		Describe("GenerateKey", func() {
			var (
				test  string
				t     []string
				err   error
				c     fohnhab.GenerateKeyRequest
				ctx   context.Context
				kinds []string
			)
			Context("When called with a valid key type argument", func() {
				kinds = []string{"aes-256", "aes-192", "aes-128"}
				BeforeEach(func() {
					t = []string{}
					for _, ty := range kinds {
						c.Kind = ty
						test, err = s.GenerateKey(ctx, c)
						t = append(t, test)
					}
				})
				It("Should not error", func() {
					Expect(err).To(Not(HaveOccurred()))
				})
				It("Should return a 256 bit key for the user as a base64 encoded string", func() {
					Expect(len(t[0])).To(Equal(44)) //44
				})
				It("Should return a 196 bit key for the user as a base64 encoded string", func() {
					Expect(len(t[1])).To(Equal(32)) //32
				})
				It("Should return a 128 bit key for the user as a base64 encoded string", func() {
					Expect(len(t[2])).To(Equal(24)) //24
				})
			})

			Context("When called with an incorrect argument", func() {
				BeforeEach(func() {
					c.Kind = "aes-2555"
					test, err = s.GenerateKey(ctx, c)
				})
				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("Galois/Counter Encryption (GCM)", func() {
			var (
				et   string
				err  error
				req  fohnhab.GCMERequest
				breq fohnhab.GCMERequest
				ctx  context.Context
			)
			Context("When called with a valid key and plaintext to encrypt", func() {
				req.Key = "lwL3V4RoI7vnoh8TbgQ16mr+M60cdLrPhpHJ923Oldw="
				req.PlainText = "This is the best message ever"
				BeforeEach(func() {
					et, err = s.GCME(ctx, req)
				})
				It("Should not error", func() {
					Expect(err).To(Not(HaveOccurred()))
				})
				It("Should return the encrypted text", func() {
					Expect(et).NotTo(Equal(""))
				})
			})

			Context("When called with an invalid key length for encryption", func() {
				breq.Key = "lwL3V4RoI7vnoh8TbgQ16mr+M60cdLrPhpHJ9"
				breq.PlainText = "This message will not be encrypted"
				BeforeEach(func() {
					et, err = s.GCME(ctx, breq)
				})
				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
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
