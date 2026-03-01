package yang_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	yeti "github.com/terjelafton/yeti/internal/yang"
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

	Describe("GetNode", func() {
		It("returns a top-level container", func() {
			node, err := tree.GetNode("/interfaces")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindContainer))
			Expect(node.Config).To(BeTrue())
		})

		It("returns a nested leaf with type info", func() {
			node, err := tree.GetNode("/interfaces/interface/mtu")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindLeaf))
			Expect(node.Type).NotTo(BeNil())
			Expect(node.Type.Name).To(Equal("uint16"))
			Expect(node.Type.Range).To(Equal("68..9192"))
			Expect(node.Default).To(Equal("1500"))
		})

		It("marks state containers as config false", func() {
			node, err := tree.GetNode("/interfaces/interface/counters")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Config).To(BeFalse())
		})

		It("returns a list node with its key", func() {
			node, err := tree.GetNode("/interfaces/interface")
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Kind).To(Equal(yeti.KindList))
			Expect(node.Key).To(Equal("name"))
		})

		It("returns an error for a nonexistent path", func() {
			_, err := tree.GetNode("/nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetChildren", func() {
		It("returns the children of a list node", func() {
			children, err := tree.GetChildren("/interfaces/interface")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(HaveLen(4))

			childMap := map[string]yeti.NodeKind{}
			for _, c := range children {
				childMap[c.Name] = c.Kind
			}
			Expect(childMap).To(Equal(map[string]yeti.NodeKind{
				"counters": yeti.KindContainer,
				"enabled":  yeti.KindLeaf,
				"mtu":      yeti.KindLeaf,
				"name":     yeti.KindLeaf,
			}))
		})

		It("returns nil for a leaf node", func() {
			children, err := tree.GetChildren("/interfaces/interface/mtu")
			Expect(err).NotTo(HaveOccurred())
			Expect(children).To(BeNil())
		})
	})
})
