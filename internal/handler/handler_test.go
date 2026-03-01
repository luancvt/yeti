package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/terjelafton/yeti/internal/handler"
	"github.com/terjelafton/yeti/internal/yang"
)

var _ = Describe("Handler", func() {
	var mux *http.ServeMux

	BeforeEach(func() {
		modelsFS := os.DirFS("../../models")
		tree, err := yang.ParseCollection(modelsFS, "test")
		Expect(err).NotTo(HaveOccurred())

		trees := map[string]*yang.CollectionTree{"test": tree}
		h := handler.New(trees)

		mux = http.NewServeMux()
		mux.HandleFunc("GET /{$}", h.Index)
		mux.HandleFunc("GET /{collection}/{module}", h.Browse)
		mux.HandleFunc("GET /models/{collection}", h.Models)
		mux.HandleFunc("GET /tree/{collection}/{module}", h.Tree)
		mux.HandleFunc("GET /tree/{collection}/{module}/{path...}", h.Tree)
		mux.HandleFunc("GET /detail/{collection}/{module}/{path...}", h.Detail)
	})

	Describe("GET /", func() {
		It("returns 200 with the page shell", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("Yeti"))
			Expect(body).To(ContainSubstring("test"))
		})
	})

	Describe("GET /{collection}/{module}", func() {
		It("returns the full page with collection and module pre-selected", func() {
			req := httptest.NewRequest("GET", "/test/test-module", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("Yeti"))
			Expect(body).To(ContainSubstring("/tree/test/test-module"))
		})

		It("returns 404 for unknown collection", func() {
			req := httptest.NewRequest("GET", "/nonexistent/foo", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /models/{collection}", func() {
		It("returns module names for a collection", func() {
			req := httptest.NewRequest("GET", "/models/test", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("test-module"))
		})

		It("returns 404 for unknown collection", func() {
			req := httptest.NewRequest("GET", "/models/nonexistent", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /tree/{collection}/{module}/{path...}", func() {
		It("returns top-level children of a module", func() {
			req := httptest.NewRequest("GET", "/tree/test/test-module", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("interfaces"))
		})

		It("returns children of a container", func() {
			req := httptest.NewRequest("GET", "/tree/test/test-module/interfaces/interface", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("name"))
			Expect(body).To(ContainSubstring("mtu"))
			Expect(body).To(ContainSubstring("enabled"))
			Expect(body).To(ContainSubstring("counters"))
		})

		It("returns 404 for unknown collection", func() {
			req := httptest.NewRequest("GET", "/tree/nonexistent/foo", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("returns 404 for unknown module", func() {
			req := httptest.NewRequest("GET", "/tree/test/nonexistent", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /detail/{collection}/{module}/{path...}", func() {
		It("returns detail for a leaf node", func() {
			req := httptest.NewRequest("GET", "/detail/test/test-module/interfaces/interface/mtu", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("mtu"))
			Expect(body).To(ContainSubstring("uint16"))
			Expect(body).To(ContainSubstring("68..9192"))
			Expect(body).To(ContainSubstring("1500"))
		})

		It("returns detail for a container", func() {
			req := httptest.NewRequest("GET", "/detail/test/test-module/interfaces", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("interfaces"))
			Expect(body).To(ContainSubstring("container"))
		})

		It("returns 404 for unknown path", func() {
			req := httptest.NewRequest("GET", "/detail/test/test-module/nonexistent", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
