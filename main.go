package main

import (
	"context"
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"strconv"
)

type Student struct {
	id    int
	name  string
	phone string
}

var arrStudents []Student

var answerStudents = widget.NewLabel("")
var answerGroups = widget.NewLabel("")
var db, err = sql.Open("sqlserver", "server=localhost;user id=Artem;password=sql12345678;database=Students;encrypt=disable")

func main() {
	w := App()
	if err != nil {
		log.Fatal("Ошибка при открытии соединения: ", err.Error())
		return
	}
	AddText(answerStudents, "Ученики")
	ReadStudents(db)
	AddText(answerGroups, "Группы")
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

	list := widget.NewList(
		func() int {
			return len(arrStudents)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(idList widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(arrStudents[idList].name)
		},
	)

	scrStudents := container.NewVScroll(list)
	scrGroups := container.NewVScroll(answerGroups)
	scrStudents.SetMinSize(fyne.NewSize(300, 600))
	scrGroups.SetMinSize(fyne.NewSize(300, 600))

	label1 := widget.NewLabel("Удалить ученика")
	entry1 := widget.NewEntry()
	btn1 := widget.NewButton("Удалить", func() {
		//TODO list.OnSelected(id)
		n, err := strconv.Atoi(entry1.Text)
		if err != nil {
			panic(err)
		}
		DeleteStudent(db, n)
	})
	label2 := widget.NewLabel("Добавить ученика")
	entry2 := widget.NewEntry()

	w.SetContent(container.NewHBox(
		scrStudents,
		container.NewVBox(
			label1,
			entry1,
			btn1,
			label2,
			entry2,
		),
	))

	return w
}

func ReadStudents(db *sql.DB) []Student {
	var id int
	var name, phone string
	var st Student
	arrStudents = nil
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	answerStudents.Text = ""
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
		st.id, st.name, st.phone = id, name, phone
		arrStudents = append(arrStudents, st)
		AddText(answerStudents, fmt.Sprintf("Id: %d, Name: %s, Phone: %s\n", st.id, st.name, st.phone))
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
	return arrStudents
}

func ReadGroup(db *sql.DB) {
	var id int
	var group string
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	count := 0
	query := "SELECT Id, GroupName FROM StudyGroup"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &group); err != nil {
			log.Fatal(err)
		}
		AddText(answerGroups, fmt.Sprintf("Id: %d, Group: %s\n", id, group))
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
}

func DeleteStudent(db *sql.DB, id int) {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "DELETE FROM Student WHERE Id = @Id"
	_, err = db.ExecContext(ctx, query, sql.Named("Id", id))
	if err != nil {
		panic(err)
	}
	ReadStudents(db)
}
