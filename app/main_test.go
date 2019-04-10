package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
	a.Initialize(
		os.Getenv("ATP_TEST_USERNAME"),
		os.Getenv("ATP_TEST_PASSWORD"),
		os.Getenv("ATP_TEST_NAME"),
	)

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test product","price":11.22}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}

}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts(4)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	payload := []byte(`{"name":"test product - updated name","price":11.22}`)

	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the ID to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])

	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}
func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products VALUES(:1, :2, :3)", nil, "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		CreatedOn := time.Now()
		a.DB.Exec("INSERT INTO users VALUES(:1, :2, :3, :4)", nil, "test user", CreatedOn, "01")
	}
}

func addSymbols(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO symbol VALUES(:1, :2, :3, :4, :5)", nil, 1, 1, 1, "ORCL")
		a.DB.Exec("INSERT INTO symbol VALUES(:1, :2, :3, :4, :5)", nil, 0, 1, 0, "SAN")
	}
}

func addTwoPortfolio(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO portfolio VALUES(:1, :2, :3, :4)", nil, 1, 1, 1)
		a.DB.Exec("INSERT INTO portfolio VALUES(:1, :2, :3, :4)", nil, 1, 2, 1)
	}
}

func addThreeNews(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO news VALUES(:1, :2, :3, :4, :5)", nil, "headline test 1", "body test 1", "www.linkTest.com", 1)
		a.DB.Exec("INSERT INTO news VALUES(:1, :2, :3, :4, :5)", nil, "headline test 2", "body test 2", "www.linkTest.com", 2)
		a.DB.Exec("INSERT INTO news VALUES(:1, :2, :3, :4, :5)", nil, "headline test 3", "body test 3", "www.linkTest.com", 1)
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("DELETE FROM account")
	a.DB.Exec("DELETE FROM symbol")
	a.DB.Exec("DELETE FROM portfolio")
	a.DB.Exec("DELETE FROM news")
	a.DB.Exec("ALTER SEQUENCE products_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE USERS_SEQ1 RESTART")
	a.DB.Exec("ALTER SEQUENCE ACCOUNT_SEQ RESTART")
	a.DB.Exec("ALTER SEQUENCE SYMBOL_SEQ RESTART")
	a.DB.Exec("ALTER SEQUENCE PORTFOLIO_SEQ RESTART")
	a.DB.Exec("ALTER SEQUENCE NEWS_SEQ RESTART")
}
func ensureTableExists() {
	tx, err := a.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(`DROP TABLE products`)
	_, err = tx.Exec(`DROP SEQUENCE products_seq`)
	_, err = tx.Exec(`DROP SEQUENCE USERS_SEQ1`)
	_, err = tx.Exec(`DROP SEQUENCE ACCOUNT_SEQ`)
	_, err = tx.Exec(`DROP SEQUENCE PORTFOLIO_SEQ`)
	_, err = tx.Exec(`DROP SEQUENCE SYMBOL_SEQ`)
	_, err = tx.Exec(`DROP SEQUENCE NEWS_SEQ`)
	_, err = tx.Exec(`CREATE TABLE products (
		id NUMBER(10) GENERATED BY DEFAULT ON NULL AS IDENTITY,
		name VARCHAR2(50) NOT NULL,
		price NUMBER NOT NULL,
		CONSTRAINT products_pk PRIMARY KEY (id))`)
	_, err = tx.Exec(`CREATE SEQUENCE products_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE`)
	_, err = tx.Exec(`CREATE SEQUENCE ACCOUNT_SEQ START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE`)
	_, err = tx.Exec(`CREATE SEQUENCE SYMBOL_SEQ START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE`)
	_, err = tx.Exec(`CREATE SEQUENCE PORTFOLIO_SEQ START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE`)
	_, err = tx.Exec(`CREATE SEQUENCE USERS_SEQ1 START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE`)
	_, err = tx.Exec(`CREATE SEQUENCE NEWS_SEQ START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE`)
	_, err = tx.Exec(`CREATE OR REPLACE TRIGGER products_on_insert
						BEFORE INSERT ON products
						FOR EACH ROW
					BEGIN
						SELECT products_seq.nextval
						INTO :new.id
						FROM dual;
					END;`)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func TestCreateUser(t *testing.T) {
	clearTable()

	payload := []byte(`{"displayname":"test user","idcsid":"010101"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["userID"] != 1.0 {
		t.Errorf("Expected UserID to be '1'. Got '%v'", m["userID"])
	}
	if m["displayName"] != "test user" {
		t.Errorf("Expected display name to be 'test user'. Got '%v'", m["displayName"])
	}

	if m["idcsID"] != "010101" {
		t.Errorf("Expected IdcsID to be '010101'. Got '%v'", m["idcsID"])
	}

}

func TestCreatePortfolioEntry(t *testing.T) {
	clearTable()
	addUsers(1)
	addSymbols(1)

	payload := []byte(`{"isWatched":1,"userID":1}`)

	req, _ := http.NewRequest("POST", "/addwatch/portfolio/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["portfolioID"] != 1.0 {
		t.Errorf("Expected Portfolio ID to be '1'. Got '%v'", m["portfolioID"])
	}
	if m["isWatched"] != 1.0 {
		t.Errorf("Expected display name to be '1'. Got '%v'", m["isWatched"])
	}
	if m["symbolID"] != 1.0 {
		t.Errorf("Expected symbol ID to be '1'. Got '%v'", m["symbolID"])
	}

}

func TestGetSymbols(t *testing.T) {
	clearTable()
	addUsers(1)
	addSymbols(1)

	req, _ := http.NewRequest("GET", "/symbols", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	testBody := response.Body.String()
	fmt.Printf("%v\n", testBody)

}

func TestGetPortfolio(t *testing.T) {
	clearTable()
	addUsers(1)
	addSymbols(1)
	addTwoPortfolio(1)

	req, _ := http.NewRequest("GET", "/portfolio/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	testBody := response.Body.String()
	fmt.Printf("%v\n", testBody)

}

func TestGetNews(t *testing.T) {
	clearTable()
	addUsers(1)
	addSymbols(1)
	addTwoPortfolio(1)
	addThreeNews(1)

	req, _ := http.NewRequest("GET", "/news/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	testBody := response.Body.String()
	fmt.Printf("News: %v\n", testBody)

}

func TestSetWatchedToFalse(t *testing.T) {
	clearTable()
	addUsers(1)
	addSymbols(1)
	addTwoPortfolio(1)

	req, _ := http.NewRequest("GET", "/portfolio/1", nil)
	response := executeRequest(req)

	originalPortfolio := response.Body.String()
	fmt.Printf("OP: %v\n", originalPortfolio)

	payload := []byte(`{"userID":1}`)

	req, _ = http.NewRequest("PUT", "/removewatch/portfolio/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req2, _ := http.NewRequest("GET", "/portfolio/1", nil)
	response2 := executeRequest(req2)

	newPortfolio := response2.Body.String()
	fmt.Printf("NP: %v\n", newPortfolio)

	if originalPortfolio == newPortfolio {
		t.Errorf("Expected different portfolios for userID 1")
	}

}

// func TestCreateAcount(t *testing.T) {
// 	clearTable()
// 	addUsers(1)

// 	payload := []byte(`{"checking":250000,"investment":0,"userid":1}`)
// 	req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusCreated, response.Code)

// 	var m map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &m)

// 	if m["accountID"] != 1.0 {
// 		t.Errorf("Expected accountID name to be '1'. Got '%v'", m["accountID"])
// 	}

// 	if m["checking"] != 250000.0 {
// 		t.Errorf("Expected checking to be '250000'. Got '%v'", m["checking"])
// 	}

// }
