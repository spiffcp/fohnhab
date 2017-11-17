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
				req.Key = "kB+BmVehHNsxaCqOaZN/s+c/aqRfk4LSfNi52SJOb3I="
				req.ToEncrypt = "Hello Peter"
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
				req.Key = "kB+BmVehHNsxaCqOaZN/s+c/aqRfk4LSfNi52SJOb3I="
				req.ToDecrypt = "YzyMcxu8OIE9Tz2nMvpHHlsJq/2CiIx2lFe2/4Ua1pu2cQqIPSWM"
				BeforeEach(func() {
					dt, err = s.GCMD(ctx, req)
					if err != nil {
						Fail(err.Error())
					}
				})
				It("Should return the decrpyted text", func() {
					Expect(dt).To(Equal("Hello Peter"))
				})
			})
		})
	})
})
