package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

func main() {
	App()
	db, err := sql.Open("sqlserver", "server=localhost;user id=Artem;password=sql12345678;database=Students;encrypt=disable")
	if err != nil {
		log.Fatal("Ошибка при открытии соединения: ", err.Error())
		return
	}
	ReadPhone(db)
	ReadGroup(db)
	defer db.Close()
}

func App() {
	newApp := app.New()
	w := newApp.NewWindow("Метод Гаусса")
	w.Resize(fyne.NewSize(300, 600))
	w.CenterOnScreen()
}

func ReadPhone(db *sql.DB) {
	var id int
	var name, phone string
	count := 0
	query := "SELECT Id, FullName, Phone FROM Student"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &phone); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Id: %d, Name: %s, Phone: %s\n", id, name, phone)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
}

func ReadGroup(db *sql.DB) {
	var id int
	var group string
	count := 0
	query := "SELECT Id, GroupName FROM StudyGroup"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &group); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Id: %d, Group: %s\n", id, group)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
}
