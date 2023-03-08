package entity_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/entity"
)

var _ = Describe("Order", func() {
	Context("Validate should identify that", func() {
		// The Validate method is called by the constructor function, so the tests are covering both

		It("ID is invalid", func() {
			Expect(entity.NewOrder("", 10.0, 190.9)).Error().To(HaveOccurred())
		})

		It("price is invalid", func() {
			Expect(entity.NewOrder("opa", 0.0, 2.0)).Error().To(HaveOccurred())
		})

		It("tax is invalid", func() {
			Expect(entity.NewOrder("opa", 100.0, 0.0)).Error().To(HaveOccurred())
		})

		It("order is valid", func() {
			Expect(entity.NewOrder("opa", 879.0, 78.0)).Error().NotTo(HaveOccurred())
		})
	})

	Context("Calculate final price", func() {
		It("Succeed", func() {
			order, err := entity.NewOrder("some order id", 120.0, 32.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(order.CalculateFinalPrice()).To(Succeed())
			Expect(order.FinalPrice).To(BeEquivalentTo(152.0))
		})

		It("Fail because order is not valid", func() {
			order := entity.Order{}
			Expect(order.CalculateFinalPrice()).NotTo(Succeed())
		})
	})
})
