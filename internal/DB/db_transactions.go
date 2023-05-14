package DB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EnderAgentX/coursework/internal/model"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var ArrStudents []model.Student
var ArrGroups []model.Group

var AnswerStudents = widget.NewLabel("")
var AnswerGroups = widget.NewLabel("")

func AddStudent(db *sql.DB, w fyne.Window, name, gender, studentCard, phone string, groupId int) {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO Student (FullName, Gender, StudentCard, Phone, GroupId) VALUES (@Name, @Gender, @StudentCard, @Phone, @GroupId); select isNull(SCOPE_IDENTITY(), -1);"
	if groupId == 0 {
		dialog.ShowError(
			errors.New("Не выбрана группа!"),
			w,
		)
	} else {

		rows, err := db.Prepare(query)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		rows.QueryRowContext(
			ctx,
			sql.Named("Name", name),
			sql.Named("Gender", gender),
			sql.Named("StudentCard", studentCard),
			sql.Named("Phone", phone),
			sql.Named("GroupId", groupId))
		ReadStudents(db)
	}
}

func AddGroup(db *sql.DB, group string) int {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query2 := "SELECT COUNT(GroupName) FROM StudyGroup WHERE GroupName = @GroupName"
	rows2, err := db.QueryContext(ctx, query2, sql.Named("GroupName", group))
	nGroups := 0
	for rows2.Next() {

		if err := rows2.Scan(&nGroups); err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println(nGroups)

	if nGroups == 0 {
		query := "INSERT INTO StudyGroup (GroupName) VALUES (@GroupName); select isNull(SCOPE_IDENTITY(), -1);"
		rows, err := db.Prepare(query)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		rows.QueryRowContext(
			ctx,
			sql.Named("GroupName", group),
		)
		ReadGroup(db)
	} else {
		return 0
	}
	return 1
}

func StudentCardDuplicate(db *sql.DB, card string) int {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query2 := "SELECT COUNT(StudentCard) FROM Student WHERE StudentCard = @StudentCard"
	rows2, err := db.QueryContext(ctx, query2, sql.Named("StudentCard", card))
	nCard := 0
	for rows2.Next() {

		if err := rows2.Scan(&nCard); err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println(nCard)

	if nCard == 0 {
		return 1
	} else {
		return 0
	}

}

func ReadStudents(db *sql.DB) []model.Student {
	var id, groupId int
	var name, gender, studentCard, phone string
	var st model.Student
	ArrStudents = nil
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	AnswerStudents.Text = ""
	count := 0
	query := "SELECT Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone, &groupId); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId = id, name, gender, studentCard, phone, groupId
		ArrStudents = append(ArrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return ArrStudents
}

func ReadStudentsGender(db *sql.DB, genderSel string) []model.Student {
	var id, groupId int
	var name, gender, studentCard, phone string
	var st model.Student
	ArrStudents = nil
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	AnswerStudents.Text = ""
	count := 0
	query := "SELECT Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student WHERE Gender = @Gender"
	rows, err := db.QueryContext(ctx, query, sql.Named("Gender", genderSel))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone, &groupId); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId = id, name, gender, studentCard, phone, groupId
		ArrStudents = append(ArrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return ArrStudents
}

func ReadSelectedGroup(db *sql.DB, groupId int) []model.Student {
	var id, groupIdBD int
	var name, gender, studentCard, phone string
	var st model.Student
	ArrStudents = nil
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	AnswerStudents.Text = ""
	count := 0
	query := "SELECT Student.Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group"
	rows, err := db.QueryContext(ctx, query, sql.Named("Group", groupId))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone, &groupIdBD); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId = id, name, gender, studentCard, phone, groupIdBD
		ArrStudents = append(ArrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return ArrStudents
}

func CardSearch(db *sql.DB, cardSearch string) ([]model.Student, int) {
	var id, groupId int
	var name, gender, studentCard, phone string
	var st model.Student

	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	AnswerStudents.Text = ""
	count := 0
	query := "SELECT Student.Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudentCard = @StudentCard"
	rows, err := db.QueryContext(ctx, query, sql.Named("StudentCard", cardSearch))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if count == 0 {
			ArrStudents = nil
		}
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone, &groupId); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId = id, name, gender, studentCard, phone, groupId
		ArrStudents = append(ArrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
	return ArrStudents, count
}

func ReadSelectedGroupGender(db *sql.DB, groupId int, genderSel string) []model.Student {
	var id, groupIdBD int
	var name, gender, studentCard, phone string
	var st model.Student
	ArrStudents = nil
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	AnswerStudents.Text = ""
	count := 0
	query := "SELECT Student.Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group"
	if genderSel == "Мужской" {
		query = "SELECT Student.Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group AND Gender = @Gender"
	} else if genderSel == "Женский" {
		query = "SELECT Student.Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group AND Gender = @Gender"
	}

	rows, err := db.QueryContext(ctx, query, sql.Named("Group", groupId), sql.Named("Gender", genderSel))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone, &groupIdBD); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId = id, name, gender, studentCard, phone, groupIdBD
		ArrStudents = append(ArrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return ArrStudents
}

func ReadGroup(db *sql.DB) []model.Group {
	var id int
	var group string
	var gr model.Group
	ArrGroups = nil
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
		gr.Id, gr.Name = id, group
		ArrGroups = append(ArrGroups, gr)
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()
	return ArrGroups
}

func GetGroupName(db *sql.DB, groupId int) string {
	var groupName string

	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "SELECT GroupName FROM StudyGroup WHERE Id = @Id"

	rows, err := db.QueryContext(ctx, query, sql.Named("Id", groupId))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&groupName); err != nil {
			log.Fatal(err)
		}
	}

	return groupName
}

func GetGroupIdByName(db *sql.DB, groupName string) int {
	var groupId int

	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "SELECT Id FROM StudyGroup WHERE GroupName = @groupName"

	rows, err := db.QueryContext(ctx, query, sql.Named("groupName", groupName))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&groupId); err != nil {
			log.Fatal(err)
		}
	}

	return groupId
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

func DeleteGroup(db *sql.DB, w fyne.Window, delId int) bool {
	ctx := context.Background()
	var cnt int
	del := false
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
		del = true
	} else {
		fmt.Println("Ошибка, удалите всех студентов!")
		dialog.ShowError(
			errors.New("Невозможно удалить группу. Необходимо удалить всех студентов из группы"),
			w,
		)

	}
	return del
}

func UpdateGroup(db *sql.DB, groupId int, newName string) {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "UPDATE StudyGroup SET GroupName = @GroupName WHERE Id = @Id"
	_, err = db.ExecContext(ctx, query, sql.Named("GroupName", newName), sql.Named("Id", groupId))
	if err != nil {
		panic(err)
	}
	ReadGroup(db)
}

func IdSearch(db *sql.DB, idSearch int) (string, string, string, string, int) {
	var id, groupId int
	var name, gender, studentCard, phone string
	var st model.Student
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	count := 0
	query := "SELECT Student.Id, FullName, Gender, StudentCard, Phone, GroupId FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE Student.Id = @IdSearch"
	rows, err := db.QueryContext(ctx, query, sql.Named("IdSearch", idSearch))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone, &groupId); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId = id, name, gender, studentCard, phone, groupId
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return st.Name, st.Gender, st.StudentCard, st.Phone, st.GroupId
}

func UpdateStudent(db *sql.DB, studentId int, fullName, gender, studentCard, phone string, groupId int) []model.Student {
	var fullNameDB, genderDB, studentCardDB, phoneDB string
	fullNameDB, genderDB, studentCardDB, phoneDB = "", "", "", ""
	fullNameDB, genderDB, studentCardDB, phoneDB = fullName, gender, studentCard, phone
	fmt.Println(groupId)

	fullNameS, genderS, studentCardS, phoneS, groupIdS := IdSearch(db, studentId)
	fmt.Println(groupIdS)
	if fullNameDB == "" {
		fullNameDB = fullNameS
	}
	if genderDB == "" {
		genderDB = genderS
	}
	if studentCardDB == "" {
		studentCardDB = studentCardS
	}
	if phoneDB == "" {
		phoneDB = phoneS
	}

	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	AnswerStudents.Text = ""
	if groupId == 0 {
		query := "UPDATE Student SET FullName = @FullName, Gender = @Gender, StudentCard = @StudentCard, Phone = @Phone WHERE Id = @Id"
		_, err = db.ExecContext(ctx, query,
			sql.Named("FullName", fullNameDB), sql.Named("Gender", genderDB), sql.Named("StudentCard", studentCardDB),
			sql.Named("Phone", phoneDB), sql.Named("Id", studentId))
	} else {
		query := "UPDATE Student SET FullName = @FullName, Gender = @Gender, StudentCard = @StudentCard, Phone = @Phone, GroupId = @GroupId WHERE Id = @Id"
		_, err = db.ExecContext(ctx, query,
			sql.Named("FullName", fullNameDB), sql.Named("Gender", genderDB), sql.Named("StudentCard", studentCardDB),
			sql.Named("Phone", phoneDB), sql.Named("GroupId", groupId), sql.Named("Id", studentId))
	}

	if err != nil {
		log.Fatal(err)
	}

	return ArrStudents
}
