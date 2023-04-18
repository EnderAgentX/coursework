package main

import (
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
		delStudentId = 0
		ListId.Refresh()
		ListName.Refresh()
		ListPhone.Refresh()

	})

	btnDelGroup := widget.NewButton("Удалить группу", func() {
		fmt.Println(delGroupId)
		DeleteGroup(db, w, delGroupId)
		scrGroups.Refresh()
		listStudents.Refresh()
		listGroups.UnselectAll()
		listStudents.UnselectAll()
		delGroupId = 0
	})

	entryName := widget.NewEntry()
	entryPhone := widget.NewEntry()
	buttonComfirmStudent := widget.NewButton("Добавить", func() {
		name := entryName.Text
		phone := entryPhone.Text
		AddStudent(db, w, name, phone, delGroupId)
		if delGroupId == 0 {
			ReadGroup(db)
		} else {
			ReadSelectedGroup(db, delGroupId)
		}
		scrStudents.Refresh()
	})
	WindowAddStudent := dialog.NewCustom("Добавить студента", "Закрыть",
		container.NewVBox(
			widget.NewLabel("Добавить ученика"),
			widget.NewLabel("ФИО"),
			entryName,
			widget.NewLabel("Номер телефона"),
			entryPhone,
			buttonComfirmStudent,
		), w)

	WindowAddStudent.Resize(fyne.NewSize(300, 200))

	btnAddStudent := widget.NewButton("Добавить студента", func() {
		if delGroupId == 0 {
			dialog.ShowError(
				errors.New("Не выбрана группа!"),
				w,
			)
		} else {
			WindowAddStudent.Show()
		}
	})

	entryGroup := widget.NewEntry()
	buttonComfirmGroup := widget.NewButton("Добавить", func() {
		name := entryGroup.Text
		AddGroup(db, name)
		scrGroups.Refresh()
	})

	WindowAddGroup := dialog.NewCustom("Добавить группу", "Закрыть",
		container.NewVBox(
			widget.NewLabel("Группа"),
			entryGroup,
			buttonComfirmGroup,
		), w)
	WindowAddGroup.Resize(fyne.NewSize(300, 200))
	btnAddGroup := widget.NewButton("Добавить группу", func() {
		WindowAddGroup.Show()
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
				btnAddGroup,
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
