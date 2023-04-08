package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

func main() {

	db, err := sql.Open("sqlserver", "server=localhost;user id=Artem;password=sql12345678;database=Students;encrypt=disable")
	if err != nil {
		log.Fatal("Ошибка при открытии соединения: ", err.Error())
		return
	}
	ReadPhone(db)
	defer db.Close()
}

func ReadPhone(db *sql.DB) {
	var phone string
	query := "SELECT Phone FROM Student"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&phone); err != nil {
			log.Fatal(err)
		}
		fmt.Println(phone)
	}
	defer rows.Close()
}
