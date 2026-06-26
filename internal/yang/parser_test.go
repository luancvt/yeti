package yang_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	yeti "github.com/luancvt/yeti/internal/yang"
)

var _ = Describe("ParseCollection", func() {
	var tree *yeti.CollectionTree

	BeforeEach(func() {
		modelsFS := os.DirFS("../../models")
		var err error
		tree, err = yeti.ParseCollection(modelsFS, "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(tree).NotTo(BeNil())
	})

	It("finds the top-level interfaces container", func() {
		children := tree.Children()
		names := make([]string, len(children))
		for i, c := range children {
			names[i] = c.Name
		}
		Expect(names).To(ContainElement("interfaces"))
	})

	Describe("ModuleNames", func() {
		It("returns sorted module names from both modules", func() {
			names := tree.ModuleNames()
			Expect(names).To(Equal([]string{"test-aux", "test-module"}))
		})
	})

	Describe("ModuleChildren", func() {
		It("returns only the specified module's top-level nodes", func() {
			children, err := tree.ModuleChildren("test-module")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(HaveLen(2))
			names := []string{children[0].Name, children[1].Name}
			Expect(names).To(ConsistOf("interfaces", "logging"))
		})

		It("returns the aux module's top-level nodes", func() {
			children, err := tree.ModuleChildren("test-aux")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(HaveLen(1))
			Expect(children[0].Name).To(Equal("system"))
		})

		It("returns error for nonexistent module", func() {
			_, err := tree.ModuleChildren("nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetNode", func() {
		It("returns a top-level container", func() {
			node, err := tree.GetNode("test-module", "/interfaces")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindContainer))
			Expect(node.Config).To(BeTrue())
		})

		It("returns a nested leaf with type info", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/mtu")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindLeaf))
			Expect(node.Type).NotTo(BeNil())
			Expect(node.Type.Name).To(Equal("uint16"))
			Expect(node.Type.Range).To(Equal("68..9192"))
			Expect(node.Default).To(Equal("1500"))
		})

		It("marks state containers as config false", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/counters")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Config).To(BeFalse())
		})

		It("returns a list node with its key", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindList))
			Expect(node.Key).To(Equal("name"))
		})

		It("returns a leaf-list node", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindLeafList))
		})

		It("returns an enum leaf with enum values populated", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/speed")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindLeaf))
			Expect(node.Type).NotTo(BeNil())
			Expect(node.Type.EnumValues).To(ConsistOf("10M", "100M", "1G", "10G"))
		})

		It("returns a pattern leaf with pattern populated", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/intf-description")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Type).NotTo(BeNil())
			Expect(node.Type.Pattern).To(HaveLen(1))
			Expect(node.Type.Pattern[0]).To(Equal("[A-Za-z][A-Za-z0-9 ]*"))
		})

		It("returns a mandatory leaf with Mandatory true", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/admin-status")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Mandatory).To(BeTrue())
		})

		It("returns a bandwidth leaf with uint64 type", func() {
			node, err := tree.GetNode("test-module", "/interfaces/interface/bandwidth")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindLeaf))
			Expect(node.Type).NotTo(BeNil())
			Expect(node.Type.Name).To(Equal("uint64"))
		})

		It("extracts Extra fields from a presence container", func() {
			node, err := tree.GetNode("test-module", "/logging")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindContainer))
			Expect(node.Extra).NotTo(BeNil())
			Expect(node.Extra["presence"]).To(ConsistOf("Enables logging"))
		})

		It("returns an error for a nonexistent path", func() {
			_, err := tree.GetNode("test-module", "/nonexistent")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an empty path", func() {
			_, err := tree.GetNode("test-module", "/")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetChildren", func() {
		It("returns the children of a list node", func() {
			children, err := tree.GetChildren("test-module", "/interfaces/interface")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(HaveLen(9))

			childMap := map[string]yeti.NodeKind{}
			for _, c := range children {
				childMap[c.Name] = c.Kind
			}
			Expect(childMap).To(Equal(map[string]yeti.NodeKind{
				"admin-status":     yeti.KindLeaf,
				"bandwidth":        yeti.KindLeaf,
				"counters":         yeti.KindContainer,
				"enabled":          yeti.KindLeaf,
				"intf-description": yeti.KindLeaf,
				"mtu":              yeti.KindLeaf,
				"name":             yeti.KindLeaf,
				"speed":            yeti.KindLeaf,
				"tags":             yeti.KindLeafList,
			}))
		})

		It("returns nil for a leaf node", func() {
			children, err := tree.GetChildren("test-module", "/interfaces/interface/mtu")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(BeNil())
		})

		It("returns children at 3+ depth levels", func() {
			children, err := tree.GetChildren("test-module", "/interfaces/interface/counters/error-counters")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(HaveLen(2))
			names := []string{children[0].Name, children[1].Name}
			Expect(names).To(ConsistOf("in-errors", "out-errors"))
		})

		It("returns error for nonexistent module", func() {
			_, err := tree.GetChildren("nonexistent", "/foo")
			Expect(err).To(HaveOccurred())
		})
	})
})
