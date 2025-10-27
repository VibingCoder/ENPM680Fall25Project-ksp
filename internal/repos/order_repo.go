package repos

import "github.com/jmoiron/sqlx"

type OrderRepo struct{ db *sqlx.DB }

func NewOrderRepo(db *sqlx.DB) *OrderRepo { return &OrderRepo{db: db} }

func (r *OrderRepo) Create(orderID, sessionID, region, fulfillment, name, email string, total float64) error {
	_, err := r.db.Exec(`
	  INSERT INTO orders(id,session_id,region_code,fulfillment,customer_name,customer_email,total,status,created_at)
	  VALUES(?,?,?,?,?,?,?,'PLACED',CURRENT_TIMESTAMP)
	`, orderID, sessionID, region, fulfillment, name, email, total)
	return err
}

func (r *OrderRepo) InsertItem(orderID, productID string, qty int, price float64, condition string) error {
	_, err := r.db.Exec(`
	  INSERT INTO order_items(order_id,product_id,qty,price,condition)
	  VALUES(?,?,?,?,?)
	`, orderID, productID, qty, price, condition)
	return err
}

type OrderRow struct {
	ID          string  `db:"id"`
	Region      string  `db:"region_code"`
	Fulfillment string  `db:"fulfillment"`
	Customer    string  `db:"customer_name"`
	Email       string  `db:"customer_email"`
	Total       float64 `db:"total"`
	Status      string  `db:"status"`
	CreatedAt   string  `db:"created_at"`
}

type OrderItemRow struct {
	Title     string  `db:"title"`
	Condition string  `db:"condition"`
	Qty       int     `db:"qty"`
	Price     float64 `db:"price"`
	Subtotal  float64 `db:"subtotal"`
}

func (r *OrderRepo) Get(orderID string) (OrderRow, []OrderItemRow, error) {
	var o OrderRow
	if err := r.db.Get(&o, `
	  SELECT id, region_code, fulfillment, customer_name, customer_email, total, status, created_at
	  FROM orders WHERE id = ?
	`, orderID); err != nil {
		return OrderRow{}, nil, err
	}

	var items []OrderItemRow
	if err := r.db.Select(&items, `
	  SELECT p.title, oi.condition, oi.qty, oi.price, (oi.qty*oi.price) AS subtotal
	  FROM order_items oi JOIN products p ON p.id = oi.product_id
	  WHERE oi.order_id = ?
	`, orderID); err != nil {
		return OrderRow{}, nil, err
	}

	return o, items, nil
}
