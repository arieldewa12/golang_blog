package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Produk struct {
	ID       string `db:"id"`
	SKU      string `db:"sku"`
	Name     string `db:"name"`
	Qty      int    `db:"qty"`
	Price    int    `db:"harga"`
	Unit     string `db:"unit"`
	IsActive bool   `db:"is_active"`
}

type ProdukModel interface {
	Upsert(db *sql.DB, input InputUpsert) error
	Delete(db *sql.DB, input InputDelete) error
	List(db *sql.DB, input InputList) ([]Produk, error)
	GetBySkus(db *sql.DB, input InputGetBySKU) ([]User, error)
}

type InputUpsert struct {
	Produk Produk `json:"produk"`
}

func Upsert(db *sql.DB, data InputUpsert) (err error) {
	var result sql.Result
	obj := data.Produk
	if obj.ID != "" {
		fmt.Println("sini1")
		sqlStatement := `UPDATE produk SET name = ?, qty = ?, harga = ?, unit =? , is_active = ? WHERE id = ?`
		result, err = db.Exec(sqlStatement, obj.Name, obj.Qty, obj.Price, obj.Unit, obj.IsActive, obj.ID)
		if err != nil {
			log.Fatalf("could not insert row: %v", err)
			return err
		}
	} else {
		fmt.Println("sini2")
		sqlStatement := `INSERT INTO produk (sku, name, qty, harga, unit, is_active) VALUES (?, ?, ?, ?, ?, ?)`
		result, err = db.Exec(sqlStatement, obj.SKU, obj.Name, obj.Qty, obj.Price, obj.Unit, obj.IsActive)
		if err != nil {
			log.Fatalf("could not insert row: %v", err)
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
		return err
	}
	// we can log how many rows were inserted
	fmt.Println("upsert", rowsAffected, "rows")
	return nil
}

type InputDelete struct {
	ID string `json:"id"`
}

func Delete(db *sql.DB, data InputDelete) (err error) {

	sqlStatement := `DELETE FROM produk WHERE id = ?`
	result, err := db.Exec(sqlStatement, data.ID)
	if err != nil {
		log.Fatalf("could not deleted row: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
		return err
	}
	// we can log how many rows were inserted
	fmt.Println("deleted", rowsAffected, "rows")
	return nil
}

type InputList struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func List(db *sql.DB, data InputList) (produks []Produk, err error) {
	sqlStatement := `SELECT id,sku,name,qty,harga,unit,is_active FROM produk order by created_at desc`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("could not find row: %v", err)
		return produks, err
	}
	defer rows.Close()
	for rows.Next() {
		var produk Produk
		err = rows.Scan(&produk.ID, &produk.SKU, &produk.Name, &produk.Qty, &produk.Price, &produk.Unit, &produk.IsActive)
		produks = append(produks, produk)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("could not find row: %v", err)
		return produks, err
	}

	// we can log how many rows were inserted
	fmt.Println("list", len(produks), "rows")
	return produks, nil
}

type InputGetBySKU struct {
	SKU string `json:"sku"`
}

func GetBySKU(db *sql.DB, data InputGetBySKU) (produks []Produk, err error) {
	sqlStatement := `SELECT id,sku,name,qty,harga,unit,is_active FROM produk WHERE sku =?`
	rows, err := db.Query(sqlStatement, data.SKU)
	if err != nil {
		log.Fatalf("could not find row: %v", err)
		return produks, err
	}
	defer rows.Close()
	for rows.Next() {
		var produk Produk
		err = rows.Scan(&produk.ID, &produk.SKU, &produk.Name, &produk.Qty, &produk.Price, &produk.Unit, &produk.IsActive)
		produks = append(produks, produk)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("could not find row: %v", err)
		return produks, err
	}

	// we can log how many rows were inserted
	fmt.Println("getBySku", len(produks), "rows")
	return produks, nil
}
