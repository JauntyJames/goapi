package main

import (
	"database/sql"
	"log"
	"time"
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

	err = tx.QueryRow("SELECT * FROM products WHERE ID = (SELECT MAX(ID) FROM products)").Scan(&p.ID, &p.Name, &p.Price)
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

//User table
type User struct {
	UserID      int    `json:"userID"`
	DisplayName string `json:"displayName"`
	CreatedOn   string `json:"createdOn"`
	IdcsID      string `json:"idcsID"`
}

//Account table
type Account struct {
	AccountID  int     `json:"accountID"`
	Checking   float64 `json:"checking"`
	Investment float64 `json:"CreatedOn"`
	UserID     int     `json:"userID"`
}

//OpenPosition table
type OpenPosition struct {
	OpenPositionID int     `json:"openPositionID"`
	Price          float64 `json:"price"`
	Shares         int     `json:"shares"`
	UserID         int     `json:"userID"`
	SymbolID       int     `json:"symbolID"`
}

//ClosedPosition table
type ClosedPosition struct {
	ClosedPositionID int     `json:"closedPositionID"`
	Price            float64 `json:"price"`
	Shares           int     `json:"shares"`
	UserID           int     `json:"userID"`
	SymbolID         int     `json:"symbolID"`
}

//Portfolio table
type Portfolio struct {
	PortfolioID int `json:"portfolioID"`
	IsWatched   int `json:"isWatched"`
	SymbolID    int `json:"symbolID"`
	UserID      int `json:"userID"`
}

//News table
type News struct {
	NewsID   int    `json:"newsID"`
	Headline string `json:"headline"`
	Body     string `json:"body"`
	Link     string `json:"link"`
	SymbolID int    `json:"symbolID"`
}

//Symbol table
type Symbol struct {
	SymbolID int    `json:"symbolID"`
	IsNASDAQ int    `json:"isNadaq"`
	IsSP500  int    `json:"isSP500"`
	IsDOW    int    `json:"isDOW"`
	Symbol   string `json:"symbol"`
}

func (u *User) createUser(db *sql.DB) error {
	sqlStatement := "INSERT INTO users VALUES (:1, :2, :3, :4)"
	CreatedOn := time.Now()
	tx, err := db.Begin()
	_, errExec := tx.Exec(sqlStatement, nil, u.DisplayName, CreatedOn, u.IdcsID)
	if errExec != nil {
		tx.Rollback()
		log.Fatal(errExec)
	}

	err = tx.QueryRow("SELECT * FROM users WHERE USER_ID = (SELECT MAX(USER_ID) FROM users)").Scan(&u.UserID, &u.DisplayName, &u.CreatedOn, &u.IdcsID)
	if err != nil {
		log.Fatal(err)
		tx.Commit()
	}

	return err
}

// func (accnt *Account) createAccount(db *sql.DB) error {
// 	sqlStatement := "INSERT INTO account VALUES (:1, :2, :3, :4)"
// 	tx, err := db.Begin()
// 	_, errExec := tx.Exec(sqlStatement, nil, accnt.Checking, accnt.Investment, accnt.UserID)
// 	if errExec != nil {
// 		tx.Rollback()
// 		log.Fatal(errExec)
// 	}

// 	err = tx.QueryRow("SELECT * FROM account WHERE ACCOUNT_ID = (SELECT MAX(ACCOUNT_ID) FROM account)").Scan(&accnt.AccountID, &accnt.Checking, &accnt.Investment, &accnt.UserID)
// 	if err != nil {
// 		log.Fatal(err)
// 		tx.Commit()
// 	}

// 	return err
// }

func getOpenPositions(db *sql.DB, start, count int) ([]OpenPosition, error) {
	rows, err := db.Query(
		"SELECT OPEN_POSITION_ID, price, shares, user_id, symbol_id  FROM OPEN_POSITION")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	OpenPositions := []OpenPosition{}

	for rows.Next() {
		var op OpenPosition
		if err := rows.Scan(&op.OpenPositionID, &op.Price, &op.Shares, &op.UserID, &op.SymbolID); err != nil {
			return nil, err
		}
		OpenPositions = append(OpenPositions, op)
	}

	return OpenPositions, nil
}

func getClosedPositions(db *sql.DB, start, count int) ([]ClosedPosition, error) {
	rows, err := db.Query(
		"SELECT CLOSED_POSITION_ID, price, shares, user_id, symbol_id  FROM CLOSED_POSITION")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ClosedPositions := []ClosedPosition{}

	for rows.Next() {
		var cp ClosedPosition
		if err := rows.Scan(&cp.ClosedPositionID, &cp.Price, &cp.Shares, &cp.UserID, &cp.SymbolID); err != nil {
			return nil, err
		}
		ClosedPositions = append(ClosedPositions, cp)
	}

	return ClosedPositions, nil
}

func (prtflo *Portfolio) createPortfolioEntry(db *sql.DB) error {
	sqlStatement := "INSERT INTO portfolio VALUES (:1, :2, :3, :4)"
	tx, err := db.Begin()
	_, errExec := tx.Exec(sqlStatement, nil, prtflo.IsWatched, prtflo.SymbolID, prtflo.UserID)
	if errExec != nil {
		tx.Rollback()
		log.Fatal(errExec)
	}

	err = tx.QueryRow("SELECT * FROM portfolio WHERE PORTFOLIO_ID = (SELECT MAX(PORTFOLIO_ID) FROM portfolio)").Scan(&prtflo.PortfolioID, &prtflo.IsWatched, &prtflo.SymbolID, &prtflo.UserID)
	if err != nil {
		log.Fatal(err)
		tx.Commit()
	}

	return err
}

func getSymbols(db *sql.DB) ([]Symbol, error) {
	rows, err := db.Query("SELECT * FROM symbol")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	symbols := []Symbol{}

	for rows.Next() {
		var s Symbol
		if err := rows.Scan(&s.SymbolID, &s.IsNASDAQ, &s.IsSP500, &s.IsDOW, &s.Symbol); err != nil {
			return nil, err
		}
		symbols = append(symbols, s)
	}

	return symbols, nil
}

func (prtflo *Portfolio) getPortfolio(db *sql.DB) ([]Portfolio, error) {
	rows, err := db.Query("SELECT * FROM portfolio WHERE USER_ID=:1", prtflo.UserID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	portfolio := []Portfolio{}

	for rows.Next() {
		var p Portfolio
		if err := rows.Scan(&p.PortfolioID, &p.IsWatched, &p.SymbolID, &p.UserID); err != nil {
			return nil, err
		}
		portfolio = append(portfolio, p)
	}

	return portfolio, nil
}

func (prtflo *Portfolio) getNews(db *sql.DB, start, count int) ([]News, error) {
	rows, err := db.Query(
		"SELECT N.NEWS_ID, N.HEADLINE, N.BODY, N.LINK, P.SYMBOL_ID FROM News N, Portfolio P WHERE P.USER_ID=:1 AND N.SYMBOL_ID = P.SYMBOL_ID OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY", prtflo.UserID, start, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	news := []News{}

	for rows.Next() {
		var n News
		if err := rows.Scan(&n.NewsID, &n.Headline, &n.Body, &n.Link, &n.SymbolID); err != nil {
			return nil, err
		}
		news = append(news, n)
	}
	return news, nil
}

func (prtflo *Portfolio) setWatchedFalse(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE portfolio SET IS_WATCHED=:1 WHERE SYMBOL_ID=:2 AND USER_ID=:3",
			0, prtflo.SymbolID, prtflo.UserID)

	return err
}
