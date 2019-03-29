package main

import (
	"database/sql"
	"log"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE ID=:1",
		p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=:1, price=:2 WHERE id=:3",
			p.Name, p.Price, p.ID)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=:1", p.ID)

	return err
}

func (p *product) createProduct(db *sql.DB) error {

	sqlStatement := "INSERT INTO products VALUES (:1, :2, :3)"
	tx, err := db.Begin()
	_, errExec := tx.Exec(sqlStatement, nil, p.Name, p.Price)
	if errExec != nil {
		tx.Rollback()
		log.Fatal(errExec)
	}

	err = tx.QueryRow("SELECT * FROM products WHERE ID = (SELECT MAX(ID) FROM PRODUCTS)").Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		log.Fatal(err)
		tx.Commit()
	}

	return err

}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name, price FROM products OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY",
		start, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
