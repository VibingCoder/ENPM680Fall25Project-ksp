package repos

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type CartRepo struct{ db *sqlx.DB }

func NewCartRepo(db *sqlx.DB) *CartRepo { return &CartRepo{db: db} }

type CartItemRow struct {
	ProductID  string  `db:"product_id"`
	Title      string  `db:"title"`
	Condition  string  `db:"condition"`
	Qty        int     `db:"qty"`
	PriceAtAdd float64 `db:"price_at_add"`
	Subtotal   float64 `db:"subtotal"`
}

func (r *CartRepo) EnsureCart(sessionID string) (string, error) {
	var cartID string
	if err := r.db.Get(&cartID, `SELECT id FROM carts WHERE session_id = ?`, sessionID); err == nil {
		return cartID, nil
	}
	_, err := r.db.Exec(`INSERT INTO carts(id,session_id,updated_at) VALUES(?,?,?)`,
		sessionID, sessionID, time.Now().Format(time.RFC3339))
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (r *CartRepo) UpsertItem(cartID, productID string, qty int, price float64) error {
	_, err := r.db.Exec(`
		INSERT INTO cart_items(cart_id,product_id,qty,price_at_add,created_at)
		VALUES(?,?,?,?,CURRENT_TIMESTAMP)
		ON CONFLICT(cart_id,product_id) DO UPDATE
		SET qty = cart_items.qty + excluded.qty, updated_at = CURRENT_TIMESTAMP
	`, cartID, productID, qty, price)
	return err
}

func (r *CartRepo) View(cartID string) ([]CartItemRow, float64, error) {
	rows := []CartItemRow{}
	if err := r.db.Select(&rows, `
	  SELECT ci.product_id, p.title, p.condition, ci.qty, ci.price_at_add,
	         (ci.qty*ci.price_at_add) AS subtotal
	  FROM cart_items ci JOIN products p ON p.id=ci.product_id
	  WHERE ci.cart_id = ?
	`, cartID); err != nil {
		return nil, 0, err
	}
	total := 0.0
	for _, it := range rows {
		total += it.Subtotal
	}
	return rows, total, nil
}

type CartItem struct {
	ProductID string  `db:"product_id"`
	Qty       int     `db:"qty"`
	Price     float64 `db:"price"` // <— was price_at_add
	Condition string  `db:"condition"`
	Title     string  `db:"title"`
}

func (r *CartRepo) Items(cartID string) ([]CartItem, error) {
	var out []CartItem
	err := r.db.Select(&out, `
	  SELECT ci.product_id, ci.qty, ci.price_at_add AS price, p.condition, p.title
	  FROM cart_items ci JOIN products p ON p.id=ci.product_id
	  WHERE ci.cart_id = ?
	`, cartID)
	return out, err
}

func (r *CartRepo) Clear(cartID string) error {
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE cart_id = ?`, cartID)
	return err
}
