package internal

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var arrStudents []Student
var arrGroups []Group

var answerStudents = widget.NewLabel("")
var answerGroups = widget.NewLabel("")

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
			obj.(*widget.Label).SetText(arrStudents[idList].Name)
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
			obj.(*widget.Label).SetText(arrGroups[idList].Name)
		},
	)

	ListId := widget.NewLabel("Id:")
	ListName := widget.NewLabel("Имя:")
	ListPhone := widget.NewLabel("Телефон:")
	var st Student
	var delStudentId int
	var delGroupId int
	listStudents.OnSelected = func(idList widget.ListItemID) {
		st.Id, st.Name, st.Phone = arrStudents[idList].Id, arrStudents[idList].Name, arrStudents[idList].Phone
		ListId.Text = "Id: " + strconv.Itoa(st.Id)
		ListName.Text = "Имя: " + st.Name
		ListPhone.Text = "Телефон: " + st.Phone
		delStudentId = arrStudents[idList].Id
		ListId.Refresh()
		ListName.Refresh()
		ListPhone.Refresh()
	}
	listGroups.OnSelected = func(idList widget.ListItemID) {
		delGroupId = arrGroups[idList].Id
		ReadSelectedGroup(Db, delGroupId)
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
		DeleteStudent(Db, delStudentId)
		if delGroupId == 0 {
			ReadGroup(Db)
		} else {
			ReadSelectedGroup(Db, delGroupId)
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

	groups := func() []string {
		ReadGroup(Db)
		var s []string
		for i := 0; i <= len(arrGroups)-1; i++ {
			is := arrGroups[i].Name
			s = append(s, is)
		}
		return s
	}
	var selGroupArr []string
	var selectedGroupId int
	selGroupArr = groups()
	selectGroup := widget.NewSelect(selGroupArr, func(s string) {
		fmt.Println(s)
		selectedGroupId = GetGroupIdByName(Db, s)
	})

	selectGroup.PlaceHolder = "Группа"

	btnDelGroup := widget.NewButton("Удалить группу", func() {
		delGroupName := GetGroupName(Db, delGroupId)
		DeleteGroup(Db, w, delGroupId)
		selectGroup.ClearSelected()
		selectGroup.Refresh()
		scrGroups.Refresh()
		listStudents.Refresh()
		fmt.Println(delGroupId)
		fmt.Println(delGroupName)
		for i := 0; i < len(selectGroup.Options); i++ {
			if selectGroup.Options[i] == delGroupName {
				fmt.Println(GetGroupIdByName(Db, delGroupName))
				selectGroup.Options = append(selectGroup.Options[:i], selectGroup.Options[i+1:]...)
			}
		}
		//fmt.Println(selectGroup.Options[0])
		selectGroup.Refresh()
		listGroups.UnselectAll()
		listStudents.UnselectAll()
		delGroupId = 0
	})

	entryName := widget.NewEntry()
	entryPhone := widget.NewEntry()
	buttonComfirmStudent := widget.NewButton("Добавить", func() {
		name := entryName.Text
		phone := entryPhone.Text
		if name == "" || phone == "" {
			dialog.ShowError(
				errors.New("Не все данные введены"),
				w,
			)
		}
		if selectedGroupId == 0 {
			if name != "" && phone != "" {
				dialog.ShowError(
					errors.New("Не выбрана группа!"),
					w,
				)
			}
		} else {
			AddStudent(Db, w, name, phone, selectedGroupId)
		}
		if delGroupId == 0 {
			ReadGroup(Db)
			fmt.Println(selectedGroupId)
			selGroupArr = groups()
			selectGroup.Refresh()
		} else {
			ReadSelectedGroup(Db, delGroupId)
			fmt.Println(selectedGroupId)
			selGroupArr = groups()
			selectGroup.Refresh()
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
			selectGroup,
			buttonComfirmStudent,
		), w)

	WindowAddStudent.Resize(fyne.NewSize(300, 200))

	btnAddStudent := widget.NewButton("Добавить студента", func() {

		WindowAddStudent.Show()

	})

	entryGroup := widget.NewEntry()
	buttonComfirmGroup := widget.NewButton("Добавить", func() {
		name := entryGroup.Text
		AddGroup(Db, name)
		selectGroup.Options = append(selectGroup.Options, name)
		selectGroup.SetSelected(selectGroup.Options[0])
		selectGroup.ClearSelected()
		selectGroup.OnChanged(selectGroup.PlaceHolder)
		selectGroup.Show()
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
