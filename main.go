package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var answer = widget.NewLabel("")

func main() {
	w := App()
	db, err := sql.Open("sqlserver", "server=localhost;user id=Artem;password=sql12345678;database=Students;encrypt=disable")
	if err != nil {
		log.Fatal("Ошибка при открытии соединения: ", err.Error())
		return
	}
	AddText(answer, "Ученики")
	ReadPhone(db)
	AddText(answer, "Группы")
	ReadGroup(db)
	w.ShowAndRun()
	defer db.Close()
}

func AddText(ans *widget.Label, text string) {
	ans.Text = ans.Text + text + "\n"
	ans.SetText(ans.Text)
}

func App() fyne.Window {
	newApp := app.New()
	w := newApp.NewWindow("Курсовая работа")
	w.Resize(fyne.NewSize(1200, 600))
	w.CenterOnScreen()
	scr := container.NewVScroll(answer)
	scr.SetMinSize(fyne.NewSize(300, 600))
	w.SetContent(container.NewVBox(
		scr,
	))

	return w
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
		AddText(answer, fmt.Sprintf("Id: %d, Name: %s, Phone: %s\n", id, name, phone))
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
		AddText(answer, fmt.Sprintf("Id: %d, Group: %s\n", id, group))
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
}
