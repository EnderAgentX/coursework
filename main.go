package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
	"strconv"
)

type Student struct {
	id    int
	name  string
	phone string
}

type Group struct {
	id   int
	name string
}

var arrStudents []Student
var arrGroups []Group

var answerStudents = widget.NewLabel("")
var answerGroups = widget.NewLabel("")

var db, err = sql.Open("sqlserver", "server=localhost;user id=Artem;password=sql12345678;database=Students;encrypt=disable")

func main() {
	os.Setenv("FYNE_THEME", "light")
	w := App()
	if err != nil {
		log.Fatal("Ошибка при открытии соединения: ", err.Error())
		return
	}

	AddText(answerStudents, "Студенты")

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
	myApp := app.New()
	w := myApp.NewWindow("Курсовая работа")
	w.Resize(fyne.NewSize(1200, 600))
	w.CenterOnScreen()

	listStudents := widget.NewList(
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
	listGroups := widget.NewList(
		func() int {
			return len(arrGroups)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(idList widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(arrGroups[idList].name)
		},
	)

	ListId := widget.NewLabel("Id:")
	ListName := widget.NewLabel("Имя:")
	ListPhone := widget.NewLabel("Телефон:")
	var st Student
	var delStudentId int
	var delGroupId int
	listStudents.OnSelected = func(idList widget.ListItemID) {
		st.id, st.name, st.phone = arrStudents[idList].id, arrStudents[idList].name, arrStudents[idList].phone
		ListId.Text = "Id: " + strconv.Itoa(st.id)
		ListName.Text = "Имя: " + st.name
		ListPhone.Text = "Телефон: " + st.phone
		delStudentId = arrStudents[idList].id
		ListId.Refresh()
		ListName.Refresh()
		ListPhone.Refresh()
	}
	listGroups.OnSelected = func(idList widget.ListItemID) {
		delGroupId = arrGroups[idList].id
		ReadSelectedGroup(db, delGroupId)
		listStudents.UnselectAll()
		listStudents.Refresh()
	}

	scrStudents := container.NewVScroll(listStudents)
	scrGroups := container.NewVScroll(listGroups)
	scrStudents.SetMinSize(fyne.NewSize(300, 600))
	scrGroups.SetMinSize(fyne.NewSize(300, 600))
	cardStudents := widget.NewCard("Студенты", "", nil)
	cardGroups := widget.NewCard("Группы", "", nil)
	cardStudents.Resize(fyne.NewSize(300, 300))

	btnDelStudent := widget.NewButton("Удалить студента", func() {
		fmt.Println(delStudentId)
		DeleteStudent(db, delStudentId)
		if delGroupId == 0 {
			ReadGroup(db)
		} else {
			ReadSelectedGroup(db, delGroupId)
		}

		scrStudents.Refresh()
		listStudents.Refresh()
		ListId.Text = "Id: "
		ListName.Text = "Имя: "
		ListPhone.Text = "Телефон: "
		listStudents.UnselectAll()
		ListId.Refresh()
		ListName.Refresh()
		ListPhone.Refresh()
		//TODO fix len label

	})

	btnDelGroup := widget.NewButton("Удалить группу", func() {
		fmt.Println(delGroupId)
		DeleteGroup(db, w, delGroupId)
		scrGroups.Refresh()
		listStudents.Refresh()
		listGroups.UnselectAll()
		listStudents.UnselectAll()
	})

	entryName := widget.NewEntry()
	entryPhone := widget.NewEntry()
	buttonComfirm := widget.NewButton("Добавить ученика", func() {
		name := entryName.Text
		phone := entryPhone.Text
		AddStudent(db, name, phone)
		if delGroupId == 0 {
			ReadGroup(db)
		} else {
			ReadSelectedGroup(db, delGroupId)
		}
		scrStudents.Refresh()
	})

	btnAddStudent := widget.NewButton("Добавить", func() {
		dialog.ShowCustom("Добавить пользователя", "Закрыть",
			container.NewVBox(
				widget.NewLabel("Добавить ученика"),
				widget.NewLabel("ФИО"),
				entryName,
				widget.NewLabel("Телефон"),
				entryPhone,
				buttonComfirm,
			), w)

	})
	boxActions := container.NewVBox(
		widget.NewCard("Действия", "", nil),
		widget.NewCard("", "", container.NewHBox(
			container.NewVBox(
				btnDelStudent,
				btnAddStudent,
			),
			container.NewVBox(
				btnDelGroup,
			),
		)),
		widget.NewCard("", "", container.NewVBox(
			ListId,
			ListName,
			ListPhone,
		)),
	)

	w.SetContent(
		container.NewHBox(
			container.NewVBox(
				cardStudents,
				scrStudents,
			),
			container.NewVBox(
				cardGroups,
				scrGroups,
			),
			boxActions,
		))

	return w
}

func AddStudent(db *sql.DB, name, phone string) {
	fmt.Println(name)
	fmt.Println(phone)
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO Student (FullName, Phone, GroupId) VALUES (@Name, @Phone, @GroupId); select isNull(SCOPE_IDENTITY(), -1);"
	rows, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	rows.QueryRowContext(
		ctx,
		sql.Named("Name", name),
		sql.Named("Phone", phone),
		sql.Named("GroupId", 2))
	ReadStudents(db)
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
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &phone); err != nil {
			log.Fatal(err)
		}
		st.id, st.name, st.phone = id, name, phone
		arrStudents = append(arrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return arrStudents
}

func ReadSelectedGroup(db *sql.DB, groupId int) []Student {
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
	query := "SELECT Student.Id, FullName, Phone FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group"
	rows, err := db.QueryContext(ctx, query, sql.Named("Group", groupId))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &phone); err != nil {
			log.Fatal(err)
		}
		st.id, st.name, st.phone = id, name, phone
		arrStudents = append(arrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return arrStudents
}

func ReadGroup(db *sql.DB) []Group {
	var id int
	var group string
	var gr Group
	arrGroups = nil
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
		gr.id, gr.name = id, group
		arrGroups = append(arrGroups, gr)
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
	return arrGroups
}

func DeleteStudent(db *sql.DB, delId int) {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "DELETE FROM Student WHERE Id = @Id"
	_, err = db.ExecContext(ctx, query, sql.Named("Id", delId))
	if err != nil {
		panic(err)
	}
	ReadStudents(db)
}

func DeleteGroup(db *sql.DB, w fyne.Window, delId int) {
	ctx := context.Background()
	var cnt int
	fmt.Println(delId)
	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "SELECT COUNT(Student.Id) FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group"
	rows, err := db.QueryContext(ctx, query, sql.Named("Group", delId))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			log.Fatal(err)
		}
		fmt.Println(cnt)
	}
	if cnt == 0 {
		queryDel := "DELETE FROM StudyGroup WHERE Id = @Id"
		_, err = db.ExecContext(ctx, queryDel, sql.Named("Id", delId))
		if err != nil {
			panic(err)
		}
		ReadGroup(db)
		ReadStudents(db)
	} else {
		fmt.Println("Ошибка, удалите всех студентов!")
		dialog.ShowError(
			errors.New("Невозможно удалить группу. Необходимо удалить всех студентов из группы"),
			w,
		)

	}
}
