package config_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/terjelafton/yeti/internal/config"
)

var _ = Describe("Load", func() {
	It("parses a valid config", func() {
		input := `
collections:
  - name: test-collection
    display: "Test Collection"
    path: test
`
		cfg, err := config.Load(strings.NewReader(input))
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.Collections).To(HaveLen(1))

		c := cfg.Collections[0]
		Expect(c.Name).To(Equal("test-collection"))
		Expect(c.Display).To(Equal("Test Collection"))
		Expect(c.Path).To(Equal("test"))
	})

	DescribeTable("rejects invalid configs",
		func(input string) {
			_, err := config.Load(strings.NewReader(input))
			Expect(err).To(HaveOccurred())
		},
		Entry("missing name", `collections: [{ display: "X", path: "x" }]`),
		Entry("missing display", `collections: [{ name: "x", path: "x" }]`),
		Entry("missing path", `collections: [{ name: "x", display: "X" }]`),
		Entry("empty collections", `collections: []`),
	)
})
