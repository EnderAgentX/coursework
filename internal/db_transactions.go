package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

func AddStudent(db *sql.DB, w fyne.Window, name, phone string, groupId int) {
	ctx := context.Background()

	// Проверка работает ли база
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO Student (FullName, Phone, GroupId) VALUES (@Name, @Phone, @GroupId); select isNull(SCOPE_IDENTITY(), -1);"
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
		st.Id, st.Name, st.Phone = id, name, phone
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
		st.Id, st.Name, st.Phone = id, name, phone
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
		gr.Id, gr.Name = id, group
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
