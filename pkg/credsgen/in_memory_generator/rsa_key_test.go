package inmemorygenerator_test

import (
	"code.cloudfoundry.org/cf-operator/pkg/credsgen"
	inmemorygenerator "code.cloudfoundry.org/cf-operator/pkg/credsgen/in_memory_generator"
	"code.cloudfoundry.org/cf-operator/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InMemoryGenerator", func() {
	var (
		generator credsgen.Generator
	)

	BeforeEach(func() {
		_, log := testing.NewTestLogger()
		generator = inmemorygenerator.NewInMemoryGenerator(log)
	})

	Describe("GenerateRSAKey", func() {
		It("generates an RSA key", func() {
			key, err := generator.GenerateRSAKey("foo")

			Expect(err).ToNot(HaveOccurred())
			Expect(key.PrivateKey).To(ContainSubstring("BEGIN RSA PRIVATE KEY"))
			Expect(key.PublicKey).To(ContainSubstring("BEGIN PUBLIC KEY"))
		})
	})
})
