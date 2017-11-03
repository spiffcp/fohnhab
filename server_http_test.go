package fohnhab_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/spiffcp/fohnhab"
)

var _ = Describe("ServerHttp", func() {
	Describe("Building a new server", func() {
		var (
			js []byte
			r  *http.Request
			c  context.Context
			h  http.Handler
			w  *httptest.ResponseRecorder
		)
		Context("When passing in an endpoints struct to NewHttpServer", func() {
			BeforeEach(func() {
				js = []byte(`{"kind":"aes-256"}`)
				r = httptest.NewRequest("POST", "/keygen", bytes.NewBuffer(js))
				h = fohnhab.NewHTTPServer(c, e)
				w = httptest.NewRecorder()
				h.ServeHTTP(w, r)
			})
			It("Should return a mux server that resonds", func() {
				Expect(w.Code).To(Equal(200))
			})
		})
	})
})
