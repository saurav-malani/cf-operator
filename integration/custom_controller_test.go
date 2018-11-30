package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Custom controller for core object", func() {
	var ()

	FContext("when correctly setup", func() {
		AfterEach(func() {
			for _, msg := range env.AllLogMessages() {
				fmt.Println(msg)
			}
		})

		It("should notice pod events", func() {
			tearDown, err := env.CreateSecret(env.Namespace, env.AnnotatedSecret("asecret"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			tearDown, err = env.CreatePod(env.Namespace, env.CustomPod("custom-pod"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			err = env.WaitForPod(env.Namespace, "custom-pod")
			Expect(err).NotTo(HaveOccurred(), "error waiting for initial pod")

			msgs := env.AllLogMessages()
			Expect(msgs[len(msgs)-1]).To(ContainSubstring("custom pod controller add triggered for custom-pod"))

			// Delete pod
			tearDown()
			err = env.WaitForPodDelete(env.Namespace, "custom-pod")
			Expect(err).NotTo(HaveOccurred(), "error waiting for pod delete")

			msgs = env.AllLogMessages()
			Expect(msgs[len(msgs)-1]).To(ContainSubstring("custom pod controller delete triggered for custom-pod"))

		})

		It("should ignore default pods", func() {
			tearDown, err := env.CreatePod(env.Namespace, env.DefaultPod("default-pod"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			err = env.WaitForPod(env.Namespace, "default-pod")
			Expect(err).NotTo(HaveOccurred(), "error waiting for initial pod")

			msgs := env.AllLogMessages()
			Expect(msgs[len(msgs)-1]).NotTo(ContainSubstring("custom pod controller add triggered for default-pod"))
		})

		It("should delete despite of finalizer", func() {
			tearDown, err := env.CreatePod(env.Namespace, env.CustomPod("custo-pod"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			err = env.WaitForPod(env.Namespace, "custo-pod")
			Expect(err).NotTo(HaveOccurred(), "error waiting for initial pod")

			// Delete pod
			tearDown()
			err = env.WaitForPodDelete(env.Namespace, "custo-pod")
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
