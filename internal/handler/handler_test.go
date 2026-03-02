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
		displayNames := map[string]string{"test": "Test Collection"}
		h := handler.New(trees, displayNames)

		mux = http.NewServeMux()
		mux.HandleFunc("GET /{$}", h.Index)
		mux.HandleFunc("GET /{collection}/{module}", h.Browse)
		mux.HandleFunc("GET /models", h.Models)
		mux.HandleFunc("GET /models/{collection}", h.Models)
		mux.HandleFunc("GET /tree/{collection}/{module}", h.Tree)
		mux.HandleFunc("GET /tree/{collection}/{module}/{path...}", h.Tree)
		mux.HandleFunc("GET /detail/{collection}/{module}/{path...}", h.Detail)
		mux.HandleFunc("GET /empty/tree", h.EmptyTree)
		mux.HandleFunc("GET /empty/detail", h.EmptyDetail)
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

	Describe("GET /models (query param)", func() {
		It("returns module names via collection query param", func() {
			req := httptest.NewRequest("GET", "/models?collection=test", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("test-module"))
		})

		It("returns 404 for nonexistent collection query param", func() {
			req := httptest.NewRequest("GET", "/models?collection=nonexistent", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /empty/tree", func() {
		It("returns 200 with empty state content", func() {
			req := httptest.NewRequest("GET", "/empty/tree", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("GET /empty/detail", func() {
		It("returns 200 with empty state content", func() {
			req := httptest.NewRequest("GET", "/empty/detail", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("Collections", func() {
		It("returns sorted collection info with display names", func() {
			modelsFS := os.DirFS("../../models")
			tree, err := yang.ParseCollection(modelsFS, "test")
			Expect(err).NotTo(HaveOccurred())

			trees := map[string]*yang.CollectionTree{
				"beta":  tree,
				"alpha": tree,
			}
			dn := map[string]string{"alpha": "Alpha Display", "beta": "Beta Display"}
			h := handler.New(trees, dn)
			collections := h.Collections()
			Expect(collections).To(HaveLen(2))
			Expect(collections[0].Name).To(Equal("alpha"))
			Expect(collections[0].Display).To(Equal("Alpha Display"))
			Expect(collections[1].Name).To(Equal("beta"))
			Expect(collections[1].Display).To(Equal("Beta Display"))
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
