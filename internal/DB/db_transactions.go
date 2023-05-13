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

func AddGroup(db *sql.DB, group string) {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

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
}

func ReadStudents(db *sql.DB) []model.Student {
	var id int
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
	query := "SELECT Id, FullName, Gender, StudentCard, Phone FROM Student"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone = id, name, gender, studentCard, phone
		ArrStudents = append(ArrStudents, st)
		count++
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
	defer rows.Close()

	return ArrStudents
}

func ReadSelectedGroup(db *sql.DB, groupId int) []model.Student {
	var id int
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
	query := "SELECT Student.Id, FullName, Gender, StudentCard, Phone FROM Student JOIN StudyGroup ON Student.GroupId = StudyGroup.Id WHERE StudyGroup.Id = @Group"
	rows, err := db.QueryContext(ctx, query, sql.Named("Group", groupId))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &gender, &studentCard, &phone); err != nil {
			log.Fatal(err)
		}
		st.Id, st.Name, st.Gender, st.StudentCard, st.Phone = id, name, gender, studentCard, phone
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
