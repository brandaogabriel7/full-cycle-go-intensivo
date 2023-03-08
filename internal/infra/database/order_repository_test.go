package database_test

import (
	"database/sql"

	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/entity"
	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/infra/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

var _ = Describe("OrderRepository", func() {
	var db *sql.DB
	BeforeEach(func() {
		var err error
		db, err = sql.Open("sqlite3", ":memory:")
		Expect(err).NotTo(HaveOccurred())
		_, err = db.Exec("CREATE TABLE orders (id VARCHAR(256) NOT NULL, price FLOAT NOT NULL, tax FLOAT NOT NULL, final_price FLOAT NOT NULL, PRIMARY KEY (id))")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		db.Close()
	})

	It("Saving order", func() {
		order, err := entity.NewOrder("123", 10.0, 1.0)
		Expect(err).NotTo(HaveOccurred())
		Expect(order.CalculateFinalPrice()).To(Succeed())
		repo := database.NewOrderRepository(db)
		err = repo.Save(order)
		Expect(err).NotTo(HaveOccurred())

		var orderResult entity.Order
		err = db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", order.ID).
			Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)
		Expect(err).NotTo(HaveOccurred())
		Expect(orderResult).To(BeEquivalentTo(*order))
	})

	It("Getting total orders", func() {
		repo := database.NewOrderRepository(db)

		order, err := entity.NewOrder("123", 10.0, 1.0)
		Expect(err).NotTo(HaveOccurred())
		Expect(order.CalculateFinalPrice()).To(Succeed())
		err = repo.Save(order)
		Expect(err).NotTo(HaveOccurred())

		total, err := repo.GetTotal()
		Expect(err).NotTo(HaveOccurred())
		Expect(total).To(Equal(1))
	})
})
