package main_test

import (
	"github.com/cloudfoundry-incubator/bbs/cmd/bbs/testrunner"
	"github.com/hashicorp/consul/api"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceRegistration", func() {
	Context("when the bbs service starts", func() {
		BeforeEach(func() {
			bbsRunner = testrunner.New(bbsBinPath, bbsArgs)
			bbsProcess = ginkgomon.Invoke(bbsRunner)
		})

		AfterEach(func() {
			ginkgomon.Kill(bbsProcess)
		})

		It("registers itself with consul", func() {
			client := consulRunner.NewConsulClient()
			services, err := client.Agent().Services()
			Expect(err).ToNot(HaveOccurred())

			Expect(services).To(HaveKeyWithValue("bbs",
				&api.AgentService{
					Service: "bbs",
					ID:      "bbs",
					Port:    bbsPort,
				}))
		})

		It("registers a TTL healthcheck", func() {
			client := consulRunner.NewConsulClient()
			checks, err := client.Agent().Checks()
			Expect(err).ToNot(HaveOccurred())

			Expect(checks).To(HaveKeyWithValue("service:bbs",
				&api.AgentCheck{
					Node:        "0",
					CheckID:     "service:bbs",
					Name:        "Service 'bbs' check",
					Status:      "passing",
					ServiceID:   "bbs",
					ServiceName: "bbs",
				}))
		})
	})
})
