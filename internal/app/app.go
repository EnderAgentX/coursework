package app

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EnderAgentX/coursework/internal/DB"
	"github.com/EnderAgentX/coursework/internal/model"
)

func App() fyne.Window {
	myApp := app.New()
	w := myApp.NewWindow("Курсовая работа")
	w.Resize(fyne.NewSize(1200, 600))
	w.CenterOnScreen()

	listStudents := widget.NewList(
		func() int {
			return len(DB.ArrStudents)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(idList widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(DB.ArrStudents[idList].Name)
		},
	)
	listGroups := widget.NewList(
		func() int {
			return len(DB.ArrGroups)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(idList widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(DB.ArrGroups[idList].Name)
		},
	)

	ListName := widget.NewLabel("Имя:")
	ListGender := widget.NewLabel("Пол:")
	ListStudentCard := widget.NewLabel("Студенческий билет:")
	ListPhone := widget.NewLabel("Телефон:")
	ListGroup := widget.NewLabel("Группа:")
	var st model.Student
	var selectedListStudentId int
	var selectedListGroupId int
	listStudents.OnSelected = func(idList widget.ListItemID) {
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId =
			DB.ArrStudents[idList].Id, DB.ArrStudents[idList].Name, DB.ArrStudents[idList].Gender,
			DB.ArrStudents[idList].StudentCard, DB.ArrStudents[idList].Phone, DB.ArrStudents[idList].GroupId
		ListName.Text = "Имя: " + st.Name
		ListGender.Text = "Пол: " + st.Gender
		ListStudentCard.Text = "Студенческий билет: " + st.StudentCard
		ListPhone.Text = "Телефон: " + st.Phone
		ListGroup.Text = "Группа: " + DB.GetGroupName(DB.Db, st.GroupId)
		selectedListStudentId = DB.ArrStudents[idList].Id
		ListName.Refresh()
		ListGender.Refresh()
		ListStudentCard.Refresh()
		ListPhone.Refresh()
		ListGroup.Refresh()
	}
	listGroups.OnSelected = func(idList widget.ListItemID) {
		selectedListGroupId = DB.ArrGroups[idList].Id
		DB.ReadSelectedGroup(DB.Db, selectedListGroupId)
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
		fmt.Println(selectedListStudentId)
		DB.DeleteStudent(DB.Db, selectedListStudentId)
		if selectedListGroupId == 0 {
			DB.ReadGroup(DB.Db)
		} else {
			DB.ReadSelectedGroup(DB.Db, selectedListGroupId)
		}

		scrStudents.Refresh()
		listStudents.Refresh()
		ListName.Text = "Имя: "
		ListGender.Text = "Пол: "
		ListStudentCard.Text = "Студенческий билет: "
		ListPhone.Text = "Телефон: "
		ListPhone.Text = "Группа: "
		listStudents.UnselectAll()
		selectedListStudentId = 0
		ListName.Refresh()
		ListGender.Refresh()
		ListStudentCard.Refresh()
		ListPhone.Refresh()
		ListGroup.Refresh()

	})

	groups := func() []string {
		DB.ReadGroup(DB.Db)
		var s []string
		for i := 0; i <= len(DB.ArrGroups)-1; i++ {
			is := DB.ArrGroups[i].Name
			s = append(s, is)
		}
		return s
	}
	var selGroupArr []string
	var selectedGroupId int
	var selectedGender string
	selGroupArr = groups()
	selectGroup := widget.NewSelect(selGroupArr, func(s string) {
		fmt.Println(s)
		selectedGroupId = DB.GetGroupIdByName(DB.Db, s)
	})

	selectGroup.PlaceHolder = "Группа"

	selGenderArr := []string{"Мужской", "Женский"}
	selectGender := widget.NewSelect(selGenderArr, func(s string) {
		fmt.Println(s)
		selectedGender = s
	})

	selectGender.PlaceHolder = "Пол"

	btnDelGroup := widget.NewButton("Удалить группу", func() {
		delGroupName := DB.GetGroupName(DB.Db, selectedListGroupId)
		del := DB.DeleteGroup(DB.Db, w, selectedListGroupId)
		scrGroups.Refresh()
		listStudents.Refresh()
		if del == true {
			fmt.Println(selectedListGroupId)
			fmt.Println(delGroupName)
			for i := 0; i < len(selectGroup.Options); i++ {
				if selectGroup.Options[i] == delGroupName {
					fmt.Println(DB.GetGroupIdByName(DB.Db, delGroupName))
					selectGroup.Options = append(selectGroup.Options[:i], selectGroup.Options[i+1:]...)
				}
			}
			selectGroup.Refresh()
			listGroups.UnselectAll()
			listStudents.UnselectAll()
			selectedListGroupId = 0
		}
	})

	entryName := widget.NewEntry()
	entryStudentCard := widget.NewEntry()
	entryPhone := widget.NewEntry()
	buttonConfirmStudent := widget.NewButton("Добавить", func() {
		name := entryName.Text
		gender := selectedGender
		studentCard := entryStudentCard.Text
		phone := entryPhone.Text
		if name == "" || phone == "" || gender == "" || studentCard == "" {
			dialog.ShowError(
				errors.New("Не все данные введены"),
				w,
			)
		} else if selectedGroupId == 0 {
			if name != "" && phone != "" && gender != "" && studentCard != "" {
				dialog.ShowError(
					errors.New("Не выбрана группа!"),
					w,
				)
			}
		} else {
			DB.AddStudent(DB.Db, w, name, gender, studentCard, phone, selectedGroupId)
		}
		if selectedListGroupId == 0 {
			DB.ReadGroup(DB.Db)
			fmt.Println(selectedGroupId)
			selGroupArr = groups()
			selectGroup.Refresh()
		} else {
			DB.ReadSelectedGroup(DB.Db, selectedListGroupId)
			fmt.Println(selectedGroupId)
			selGroupArr = groups()
			selectGroup.Refresh()
		}
		scrStudents.Refresh()
		selectedListGroupId = 0
	})

	WindowAddStudent := dialog.NewCustom("Добавить студента", "Закрыть",
		container.NewVBox(
			widget.NewLabel("Добавить ученика"),
			widget.NewLabel("ФИО"),
			entryName,
			widget.NewLabel("Студенческий билет"),
			entryStudentCard,
			widget.NewLabel("Номер телефона"),
			entryPhone,
			selectGender,
			selectGroup,
			buttonConfirmStudent,
		), w)

	WindowAddStudent.Resize(fyne.NewSize(300, 200))

	btnAddStudent := widget.NewButton("Добавить студента", func() {
		entryName.Text = ""
		entryPhone.Text = ""
		selectGroup.ClearSelected()
		selectGender.ClearSelected()
		for i := 0; i < len(selectGroup.Options); i++ {
			if DB.GetGroupName(DB.Db, selectedListGroupId) == selectGroup.Options[i] {
				selectGroup.SetSelected(selectGroup.Options[i])
			}
		}

		WindowAddStudent.Show()

	})

	entryGroup := widget.NewEntry()
	buttonConfirmAddGroup := widget.NewButton("Добавить", func() {
		name := entryGroup.Text
		if name == "" {
			dialog.ShowError(
				errors.New("Не все данные введены!"),
				w,
			)
		} else {
			DB.AddGroup(DB.Db, name)
			selectGroup.Options = append(selectGroup.Options, name)
			selectGroup.SetSelected(selectGroup.Options[0])
			selectGroup.ClearSelected()
			selectGroup.OnChanged(selectGroup.PlaceHolder)
			selectGroup.Show()
			scrGroups.Refresh()
		}
	})

	WindowAddGroup := dialog.NewCustom("Добавить группу", "Закрыть",
		container.NewVBox(
			widget.NewLabel("Группа"),
			entryGroup,
			buttonConfirmAddGroup,
		), w)
	WindowAddGroup.Resize(fyne.NewSize(250, 200))
	btnAddGroup := widget.NewButton("Добавить группу", func() {
		entryGroup.Text = ""
		WindowAddGroup.Show()

	})

	buttonConfirmEditStudent := widget.NewButton("Изменить", func() {
		DB.UpdateStudent(DB.Db, selectedListStudentId, entryName.Text, selectedGender, entryStudentCard.Text,
			entryPhone.Text, selectedGroupId)
		if selectedListGroupId != 0 {
			fmt.Println("selected")
			fmt.Println(selectedListGroupId)
			DB.ReadSelectedGroup(DB.Db, selectedListGroupId)
		} else {
			fmt.Println("all")
			fmt.Println(selectedListGroupId)
			DB.ReadStudents(DB.Db)
		}
		listStudents.Refresh()
		scrStudents.Refresh()
		selectedGroupId = 0
		selectedListGroupId = 0
	})

	WindowEditStudent := dialog.NewCustom("Изменить студента", "Закрыть",
		container.NewVBox(
			widget.NewLabel("Изменить ученика"),
			widget.NewLabel("ФИО"),
			entryName,
			widget.NewLabel("Студенческий билет"),
			entryStudentCard,
			widget.NewLabel("Номер телефона"),
			entryPhone,
			selectGender,
			selectGroup,
			buttonConfirmEditStudent,
		), w)

	WindowEditStudent.Resize(fyne.NewSize(300, 200))

	btnEditStudent := widget.NewButton("Изменить студента", func() {
		fmt.Println(DB.IdSearch(DB.Db, selectedListStudentId))
		entryName.Text = ""
		entryStudentCard.Text = ""
		entryPhone.Text = ""
		selectGender.ClearSelected()
		selectGroup.ClearSelected()
		WindowEditStudent.Show()

	})

	buttonConfirmEditGroup := widget.NewButton("Изменить", func() {
		editGroupName := DB.GetGroupName(DB.Db, selectedGroupId)
		fmt.Println(selectedGroupId)
		fmt.Println(editGroupName)
		fmt.Println(entryGroup.Text)
		DB.UpdateGroup(DB.Db, selectedGroupId, entryGroup.Text)
		for i := 0; i < len(selectGroup.Options); i++ {
			if selectGroup.Options[i] == editGroupName {
				selectGroup.Options[i] = entryGroup.Text
			}
		}
		for i := 0; i < len(selectGroup.Options); i++ {
			if DB.GetGroupName(DB.Db, selectedListGroupId) == selectGroup.Options[i] {
				selectGroup.SetSelected(selectGroup.Options[i])
			}
		}
		selectGroup.Refresh()
		listGroups.Refresh()
		selectedListGroupId = 0
		selectedGroupId = 0
	})

	WindowEditGroup := dialog.NewCustom("Изменить группу", "Закрыть",
		container.NewVBox(
			widget.NewLabel("Группа"),
			selectGroup,
			entryGroup,
			buttonConfirmEditGroup,
		), w)
	WindowEditGroup.Resize(fyne.NewSize(250, 200))
	btnEditGroup := widget.NewButton("Изменить группу", func() {
		entryGroup.Text = ""
		selectGroup.ClearSelected()
		for i := 0; i < len(selectGroup.Options); i++ {
			if DB.GetGroupName(DB.Db, selectedListGroupId) == selectGroup.Options[i] {
				selectGroup.SetSelected(selectGroup.Options[i])
			}
		}

		WindowEditGroup.Show()
	})

	btnMale := widget.NewButton("Показать студентов мужского пола", func() {
		fmt.Println(selectedGroupId)
		if selectedListGroupId == 0 {
			DB.ReadStudentsGender(DB.Db, "Мужской")
		} else {
			DB.ReadSelectedGroupGender(DB.Db, selectedListGroupId, "Мужской")
		}
		listStudents.Refresh()
	})
	btnFemale := widget.NewButton("Показать студентов женского пола", func() {
		if selectedListGroupId == 0 {
			DB.ReadStudentsGender(DB.Db, "Женский")
		} else {
			DB.ReadSelectedGroupGender(DB.Db, selectedListGroupId, "Женский")
		}
		listStudents.Refresh()
	})

	labelStudentCardSearch := widget.NewLabel("Поиск по студенческому билету")
	entryStudentCardSearch := widget.NewEntry()
	btnSearch := widget.NewButton("Поиск", func() {
		if entryStudentCardSearch.Text == "" {
			dialog.ShowError(
				errors.New("Не все данные введены!"),
				w,
			)
		} else {

			_, count := DB.CardSearch(DB.Db, entryStudentCardSearch.Text)
			if count == 0 {
				dialog.ShowError(
					errors.New("Студенты не найдены!"),
					w,
				)
			} else {
				listStudents.Select(0)
			}
			listStudents.Refresh()
		}

	})
	btnCancel := widget.NewButton("Отмена", func() {
		listGroups.UnselectAll()
		DB.ReadStudents(DB.Db)
		ListName.Text = "Имя: "
		ListGender.Text = "Пол: "
		ListStudentCard.Text = "Студенческий билет: "
		ListPhone.Text = "Телефон: "
		ListGroup.Text = "Группа: "
		entryStudentCardSearch.Text = ""
		selectedListStudentId = 0
		selectGroup.ClearSelected()
		entryStudentCardSearch.Refresh()
		ListName.Refresh()
		ListGender.Refresh()
		ListStudentCard.Refresh()
		ListPhone.Refresh()
		ListGroup.Refresh()
		listStudents.Refresh()
	})

	boxActions := container.NewVBox(
		widget.NewCard("Действия", "", nil),

		widget.NewCard("", "", container.NewHBox(
			container.NewVBox(
				container.NewHBox(
					container.NewVBox(
						btnAddStudent,
						btnEditStudent,
						btnDelStudent,
					),
					container.NewVBox(
						btnAddGroup,
						btnEditGroup,
						btnDelGroup,
					),
				),
				container.NewVBox(
					btnMale,
					btnFemale,
					labelStudentCardSearch,
				),
				entryStudentCardSearch,
				container.NewHBox(
					btnSearch,
					btnCancel,
				),
			),
		)),
		widget.NewCard("", "", container.NewVBox(
			ListName,
			ListGender,
			ListStudentCard,
			ListPhone,
			ListGroup,
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
