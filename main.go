package main

import (
	"context"
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
	ctx := context.Background()
	err = db.PingContext(ctx)
	count, err := ReadEmployees(db)
	fmt.Printf("Read %d row(s) successfully.\n", count)

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
	defer rows.Close()
	defer db.Close()
}

func ReadEmployees(db *sql.DB) (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("SELECT Id, FullName, Phone FROM Student;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {
		var name, phone string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &name, &phone)
		if err != nil {
			return -1, err
		}

		fmt.Printf("Id: %d, Name: %s, Phone: %s\n", id, name, phone)
		count++
	}

	return count, nil
}
