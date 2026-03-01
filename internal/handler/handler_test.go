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
		mux.HandleFunc("GET /", h.Index)
		mux.HandleFunc("GET /tree/{collection}/{path...}", h.TreeChildren)
		mux.HandleFunc("GET /detail/{collection}/{path...}", h.Detail)
	})

	Describe("GET /", func() {
		It("returns 200 with HTML containing the tree", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("YANG Tree"))
			Expect(body).To(ContainSubstring("interfaces"))
		})
	})

	Describe("GET /tree/{collection}/{path...}", func() {
		It("returns children of a container", func() {
			req := httptest.NewRequest("GET", "/tree/test/interfaces", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("interface"))
		})

		It("returns children of a list", func() {
			req := httptest.NewRequest("GET", "/tree/test/interfaces/interface", nil)
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

		It("returns 404 for unknown path", func() {
			req := httptest.NewRequest("GET", "/tree/test/nonexistent", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /detail/{collection}/{path...}", func() {
		It("returns detail for a leaf node", func() {
			req := httptest.NewRequest("GET", "/detail/test/interfaces/interface/mtu", nil)
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
			req := httptest.NewRequest("GET", "/detail/test/interfaces", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body := w.Body.String()
			Expect(body).To(ContainSubstring("interfaces"))
			Expect(body).To(ContainSubstring("container"))
		})

		It("returns 404 for unknown path", func() {
			req := httptest.NewRequest("GET", "/detail/test/nonexistent", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
