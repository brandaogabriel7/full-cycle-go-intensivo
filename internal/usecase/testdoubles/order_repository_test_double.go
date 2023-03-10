package testdoubles

import (
	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/entity"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryTestDouble struct {
	mock.Mock
}

func (c *OrderRepositoryTestDouble) Save(order *entity.Order) error {
	args := c.Called(order)
	return args.Error(0)
}

func (c *OrderRepositoryTestDouble) GetTotal() (int, error) {
	args := c.Called()
	return args.Get(0).(int), args.Error(1)
}
