package fohnhab_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/spiffcp/fohnhab"
)

var _ = Describe("Transport", func() {
	var (
		ctx context.Context
		err error
	)

	Describe("Request layer", func() {
		var (
			r  *http.Request
			js []byte
			t  interface{}
		)
		Context("When Decoding a valid GenerateKeyRequest", func() {
			var (
				ex fohnhab.GenerateKeyRequest
			)
			ex.Kind = "aes-256"
			BeforeEach(func() {
				js = []byte(`{"kind":"aes-256"}`)
				r = httptest.NewRequest("POST", "http://encryptionService.com", bytes.NewBuffer(js))
				t, err = fohnhab.DecodeGenerateKeyRequest(ctx, r)
			})
			It("Does not error", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
			It("Returns a valid Generate Key Request type", func() {
				Expect(t.(fohnhab.GenerateKeyRequest)).To(Equal(ex))
			})
		})

		Context("When decoding an invalid GenerateKeyRequest", func() {
			BeforeEach(func() {
				js = []byte(`{""kind"""aes-256"}`)
				r = httptest.NewRequest("POST", "http://encryptionService.com", bytes.NewBuffer(js))
				t, err = fohnhab.DecodeGenerateKeyRequest(ctx, r)
			})
			It("Should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Response layer", func() {
		var (
			w    *httptest.ResponseRecorder
			r    interface{}
			body []byte
			resp *http.Response
		)
		Context("When returning a json response for a successful service request", func() {
			BeforeEach(func() {
				r = fohnhab.GenerateKeyResponse{
					Key: "57129476B8A8421DC968DA99B2B3F",
				}
				w = httptest.NewRecorder()
				err = fohnhab.EncodeGenerateKeyResponse(ctx, w, r)
				resp = w.Result()
				body, _ = ioutil.ReadAll(resp.Body)
			})
			It("Should not error", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
			It("Should return status code 200", func() {
				Expect(resp.Status).To(Equal("200 OK"))
			})
			It("Should return correctly formatted json", func() {
				var js map[string]interface{}
				Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(json.Unmarshal([]byte(body), &js)).To(BeNil())
			})
		})
		Context("When returning a response for an unsuccesful service request", func() {
			BeforeEach(func() {
				r = fohnhab.GenerateKeyResponse{
					Key: "",
					Err: "Type aes-243 not found",
				}
				w = httptest.NewRecorder()
				err = fohnhab.EncodeGenerateKeyResponse(ctx, w, r)
				resp = w.Result()
				body, _ = ioutil.ReadAll(resp.Body)
			})
			It("Should not error", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
			It("Should have a 400 status code", func() {
				Expect(resp.StatusCode).To(Equal(400))
			})
		})
	})

	Describe("Endpoints", func() {
		var (
			req fohnhab.GenerateKeyRequest
			res interface{}
		)
		Describe("MakeGenerateKeyEndpoint", func() {
			Context("When passed a valid Service", func() {
				BeforeEach(func() {
					req.Kind = "aes-256"
					res, err = e.GenerateKeyEndpoint(ctx, req)
				})
				It("Should return an endpoint implementation", func() {
					Expect(reflect.TypeOf(e.GenerateKeyEndpoint).Kind()).To(Equal(reflect.Func))
				})
				It("Should not error when called", func() {
					Expect(err).To(Not(HaveOccurred()))
				})
				It("Should call the GenerateKey function", func() {
					Expect(res.(fohnhab.GenerateKeyResponse).Err).To(Equal(""))
				})
				It("Should put errors on the response correctly", func() {
					req.Kind = "aes-246"
					res, _ = e.GenerateKeyEndpoint(ctx, req)
					Expect(res.(fohnhab.GenerateKeyResponse).Err).To(Not(BeNil()))
				})
			})
		})
	})
})
