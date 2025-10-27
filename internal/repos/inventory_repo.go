package repos

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type InventoryRepo struct{ db *sqlx.DB }

func NewInventoryRepo(db *sqlx.DB) *InventoryRepo { return &InventoryRepo{db: db} }

func (r *InventoryRepo) Qty(productID, region string) (int, error) {
	var qty int
	err := r.db.Get(&qty, `
		SELECT qty FROM inventory
		WHERE product_id = ? AND region_code = ?
	`, productID, region)
	if err != nil {
		return 0, err
	}
	return qty, nil
}

func (r *InventoryRepo) Decrement(productID, region string, by int) error {
	res, err := r.db.Exec(`
	  UPDATE inventory SET qty = qty - ?
	  WHERE product_id = ? AND region_code = ? AND qty >= ?
	`, by, productID, region, by)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("insufficient stock for %s", productID)
	}
	return nil
}
