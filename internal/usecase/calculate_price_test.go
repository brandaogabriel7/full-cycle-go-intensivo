package usecase_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/usecase"
	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/usecase/testdoubles"
)

var _ = Describe("CalculatePrice", func() {
	var cp usecase.CalculateFinalPrice
	var repoTestDouble *testdoubles.OrderRepositoryTestDouble

	Context("successfully calculate price and store order", func() {
		BeforeEach(func() {
			repoTestDouble = &testdoubles.OrderRepositoryTestDouble{}
			repoTestDouble.On("Save", mock.Anything).Return(nil)

			cp = *usecase.NewCalculateFinalPrice(repoTestDouble)
		})

		It("should calculate final price", func() {
			input := usecase.OrderInputDTO{
				ID:    "123",
				Price: 100,
				Tax:   10,
			}
			output, err := cp.Execute(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(output.FinalPrice).To(Equal(110.0))
		})
	})

	Context("fail to calculate price", func() {
		BeforeEach(func() {
			repoTestDouble = &testdoubles.OrderRepositoryTestDouble{}
			cp = *usecase.NewCalculateFinalPrice(repoTestDouble)
		})

		DescribeTable("should return error when order is invalid", func(input usecase.OrderInputDTO) {
			repoTestDouble.On("Save", mock.Anything).Return(nil)
			_, err := cp.Execute(input)
			Expect(err).To(HaveOccurred())
		},
			Entry("when there's no id", usecase.OrderInputDTO{
				Price: 100,
				Tax:   10,
			}),
			Entry("when there's no price", usecase.OrderInputDTO{
				ID:  "123",
				Tax: 10,
			}),
			Entry("when there's no tax", usecase.OrderInputDTO{
				ID:    "123",
				Price: 100,
			}),
			Entry("when there's no id or price", usecase.OrderInputDTO{
				Tax: 10,
			}),
			Entry("when there's no id or tax", usecase.OrderInputDTO{
				Price: 100,
			}),
			Entry("when there's no tax or price", usecase.OrderInputDTO{
				ID: "123",
			}),
			Entry("when there's no tax or price or id", usecase.OrderInputDTO{}),
		)

		It("should return error when repository fails", func() {
			repoTestDouble.On("Save", mock.Anything).Return(errors.New("couldn't save order"))
			input := usecase.OrderInputDTO{
				ID:    "123",
				Price: 100,
				Tax:   10,
			}
			_, err := cp.Execute(input)
			Expect(err).To(HaveOccurred())
		})
	})
})
