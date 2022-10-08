package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetEmpById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error : %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(1, "Atharva", "j.atharva12@gmail.com", "Student")

	tcs := []struct {
		desc      string
		id        int
		employee  Employee
		mockQuery interface{}
		expectErr error
	}{

		{
			desc:      "success",
			id:        1,
			employee:  Employee{1, "Atharva", "j.atharva12@gmail.com", "Student"},
			mockQuery: mock.ExpectPrepare("SELECT * FROM employee2 where id = ?").ExpectQuery().WithArgs(1).WillReturnRows(rows),
		},
		{
			desc:      "err no rows",
			id:        2,
			employee:  Employee{},
			mockQuery: mock.ExpectPrepare("SELECT * FROM employee2 where id = ?").ExpectQuery().WithArgs(2).WillReturnError(sql.ErrNoRows),
			expectErr: sql.ErrNoRows,
		},
		{
			desc:      "err preparing",
			id:        2,
			employee:  Employee{},
			mockQuery: mock.ExpectPrepare("SELECT * FROM employee2 where id = ?").WillReturnError(errors.New("error preparing select query")),
			expectErr: errors.New("error preparing select query"),
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			emp, err := getEmpById(db, tc.id)
			fmt.Println("EMPLOYEE : ", emp, err)
			if err != nil && err.Error() != tc.expectErr.Error() {
				t.Errorf("expected error : %v , got: %v", tc.expectErr, err)
			}
			if !reflect.DeepEqual(emp, tc.employee) {
				t.Errorf("Expected value: %v , got: %v", tc.employee, emp)
			}
		})
	}
	fmt.Println("Complete GetById")
}

func TestInsertEmp(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("Error Occured: ", err)
	}
	defer db.Close()

	tcs := []struct {
		desc      string
		id        int
		employee  Employee
		result    sql.Result
		mockQuery interface{}
		expectErr error
	}{
		{
			desc:      "success case",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			mockQuery: mock.ExpectPrepare("INSERT INTO employee2 (name,email,role) VALUES(?,?,?)").ExpectExec().WithArgs("Punit", "punitj1221@gmail.com", "SDE Intern").WillReturnResult(sqlmock.NewResult(1, 1)),
			result:    sqlmock.NewResult(1, 1),
		},
		{
			desc:      "err while execution",
			id:        1,
			mockQuery: mock.ExpectPrepare("INSERT INTO employee2 (name,email,role) VALUES(?,?,?)").ExpectExec().WithArgs("Punit", "punitj1221@gmail.com", "SDE Intern").WillReturnError(errors.New("errors executing insert query")),
			expectErr: errors.New("error executing insert query"),
		},
		{
			desc:      "err while preparation",
			id:        1,
			mockQuery: mock.ExpectPrepare("INSERT INTO employee2 (name,email,role) VALUES(?,?,?)").WillReturnError(errors.New("error preparing insert query")),
			expectErr: errors.New("error preparing insert query"),
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			res, err := insertEmp(db, tc.employee.name, tc.employee.email, tc.employee.role)
			fmt.Println("Result : ", res, "Error: ", err)
			if err != nil && (err.Error() != tc.expectErr.Error()) {
				t.Errorf("desc: %v ----> expected error: %v , got: %v", tc.desc, tc.expectErr, err)
			}
			if err == nil {
				gotLastInsertId, _ := res.LastInsertId()
				gotRowsAffected, _ := res.RowsAffected()
				expLastInsertId, _ := tc.result.LastInsertId()
				expRowsAffected, _ := tc.result.RowsAffected()
				if gotLastInsertId != expLastInsertId || gotRowsAffected != expRowsAffected {

					t.Errorf("desc: %v ----> Expected Result: Result(%v,%v), got: Result(%v,%v)", tc.desc, expLastInsertId, expRowsAffected, gotLastInsertId, gotRowsAffected)
				}
			}
		})
	}

	fmt.Println("INSERT COMPLETE......")
}

func TestUpdateEmp(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("Error Occured: ", err)
	}
	defer db.Close()
	tcs := []struct {
		desc      string
		id        int
		employee  Employee
		mockQuery interface{}
		expectErr error
		result    sql.Result
		updates   []updates
	}{
		{
			desc:      "success",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			updates:   []updates{{"name", "Punit Jain"}, {"email", "punitjain1221@gmail.com"}, {"role", "SDE Intern"}},
			mockQuery: mock.ExpectPrepare("UPDATE employee2 SET name = ?, email = ? , role = ? where id = ?").ExpectExec().WithArgs("Punit Jain", "punitjain1221@gmail.com", "SDE Intern", 1).WillReturnResult(sqlmock.NewResult(1, 1)),
			result:    sqlmock.NewResult(1, 1),
		},
		{
			desc:      "err while preparing",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			mockQuery: mock.ExpectPrepare("UPDATE employee2 SET name = ?, email = ? , role = ? where id = ?").WillReturnError(errors.New("failed preparing update query")),
			result:    sqlmock.NewResult(1, 1),
			expectErr: errors.New("failed preparing update query"),
		},
		{
			desc:      "err while exec",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			mockQuery: mock.ExpectPrepare("UPDATE employee2 SET name = ?, email = ? , role = ? where id = ?").ExpectExec().WillReturnError(errors.New("failed executing update query")),
			expectErr: errors.New("failed executing update query"),
		},
		{
			desc:      "err updating id",
			id:        1,
			updates:   []updates{{"id", "3"}, {"name", "Punit Jain"}, {"email", "punitjain1221@gmail.com"}, {"role", "SDE Intern"}},
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			expectErr: errors.New("cannot update id: PRIMARY KEY"),
		},
	}

	for i, tc := range tcs {

		t.Run("", func(t *testing.T) {
			res, err := tc.employee.updateEmp(db, tc.updates)
			fmt.Println("Result : ", res, "Error: ", err)
			if err != nil && (err.Error() != tc.expectErr.Error()) {
				t.Errorf("desc: %v ---> expected error: %v , got: %v", tc.desc, tc.expectErr, err)
			}
			if err == nil {
				gotLastInsertId, _ := res.LastInsertId()
				gotRowsAffected, _ := res.RowsAffected()
				expLastInsertId, _ := tc.result.LastInsertId()
				expRowsAffected, _ := tc.result.RowsAffected()
				if gotLastInsertId != expLastInsertId || gotRowsAffected != expRowsAffected {

					t.Errorf("desc: %v ----> Expected Result: Result(%v,%v), got: Result(%v,%v)", tc.desc, expLastInsertId, expRowsAffected, gotLastInsertId, gotRowsAffected)
				}
			}
		})
		fmt.Printf("TEST CASE %v completed\n", i+1)
	}
}

func TestDeleteEmp(t *testing.T) {
	fmt.Println("DELETE TESTING STARTED")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("Error Occured: ", err)
	}
	defer db.Close()

	tcs := []struct {
		desc      string
		id        int
		employee  Employee
		mockQuery interface{}
		expectErr error
		result    sql.Result
	}{
		{
			desc:      "success",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			mockQuery: mock.ExpectPrepare("DELETE FROM employee2 where id = ?").ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1)),
			result:    sqlmock.NewResult(0, 1),
		},
		{
			desc:      "invalid employee",
			id:        1,
			employee:  Employee{0, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			result:    sqlmock.NewResult(0, 0),
			expectErr: errors.New("employee not found"),
		},
		{
			desc:      "err while query preparation",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			mockQuery: mock.ExpectPrepare("DELETE FROM employee2 where id = ?").WillReturnError(errors.New("error preparing delete query")),
			expectErr: errors.New("error preparing delete query"),
		},
		{
			desc:      "err while query preparation",
			id:        1,
			employee:  Employee{1, "Punit", "punitj1221@gmail.com", "SDE Intern"},
			mockQuery: mock.ExpectPrepare("DELETE FROM employee2 where id = ?").ExpectExec().WillReturnError(errors.New("failed executing delete query")),
			expectErr: errors.New("failed executing delete query"),
		},
	}

	for i, tc := range tcs {

		t.Run("", func(t *testing.T) {
			res, err := tc.employee.deleteEmp(db)
			fmt.Println("Result : ", res, "Error: ", err)
			if err != nil && (err.Error() != tc.expectErr.Error()) {
				t.Errorf("desc: %v ----> expected error: %v , got: %v", tc.desc, tc.expectErr, err)
				return
			}
			if err == nil {
				gotLastInsertId, _ := res.LastInsertId()
				gotRowsAffected, _ := res.RowsAffected()
				expLastInsertId, _ := tc.result.LastInsertId()
				expRowsAffected, _ := tc.result.RowsAffected()
				if gotLastInsertId != expLastInsertId || gotRowsAffected != expRowsAffected {

					t.Errorf("desc: %v ----> Expected Result: Result(%v,%v), got: Result(%v,%v)", tc.desc, expLastInsertId, expRowsAffected, gotLastInsertId, gotRowsAffected)
				}
			}
		})
		fmt.Printf("TEST CASE %v completed\n", i+1)
	}
	fmt.Println("DELETE TESTING FINISHED.....")
}
