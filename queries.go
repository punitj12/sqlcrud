package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/punitj1221/connecting-sql/conn"
)

type Employee struct {
	id    int
	name  string
	email string
	role  string
}

var (
	name  string
	email string
	role  string
)

/*CREATE TABLE FUNCTION */
// func create(db *sql.DB) {
// 	create := "CREATE TABLE if NOT EXISTS employee2(id int PRIMARY KEY AUTO_INCREMENT, name varchar(20), email varchar(25), role varchar(25));"
// 	_, er := db.Exec(create)
// 	if er != nil {
// 		fmt.Println(er)
// 		panic(er)
// 	}
// 	fmt.Println("Table Created")
// }

type updates struct {
	key   string
	value string
}

func (e *Employee) updateEmp(db *sql.DB, up []updates) (sql.Result, error) {
	var res sql.Result
	var er error
	empCopy := e
	for i := range up {
		if up[i].key == "id" {
			return res, errors.New("cannot update id: PRIMARY KEY")
		}
		if up[i].key == "name" {
			empCopy.name = up[i].value
		}
		if up[i].key == "email" {
			empCopy.email = up[i].value
		}
		if up[i].key == "role" {
			empCopy.role = up[i].value
		}
	}

	stmp, err := db.Prepare("UPDATE employee2 SET name = ?, email = ? , role = ? where id = ?")
	if err == nil {
		defer stmp.Close()

		res, er = stmp.Exec(empCopy.name, empCopy.email, empCopy.role, empCopy.id)
		if er == nil {
			return res, nil
		} else {
			return res, errors.New("failed executing update query")
		}
	} else {
		return res, errors.New("failed preparing update query")
	}
}

func insertEmp(db *sql.DB, name string, email string, role string) (sql.Result, error) {

	/* WITH PREPARED STATEMENT */
	var res sql.Result
	var err error
	stmp, er := db.Prepare("INSERT INTO employee2 (name,email,role) VALUES(?,?,?)")
	if er == nil {
		defer stmp.Close()
		res, err = stmp.Exec(name, email, role)
		if err != nil {
			return res, errors.New("error executing insert query")
		}
	} else {
		return res, errors.New("error preparing insert query")
	}
	return res, err
	/* W/O PREPARED STATEMENT */
	// insert := `INSERT INTO employee2(name,email,role) VALUES( "` + name + `","` + email + `","` + role + `");`
	// fmt.Println(insert)
	// _, err := db.Exec(insert)
	// if err == nil {
	// 	fmt.Println("Insert Query Executed")
	// } else {
	// 	fmt.Println(err)
	// }
}

func getEmpById(db *sql.DB, id int) (Employee, error) {
	/* USING PREPARED STATEMENTS */
	stmp, er := db.Prepare("SELECT * FROM employee2 where id = ?")
	var em Employee
	var err error
	if er == nil {
		defer stmp.Close()

		row := stmp.QueryRow(id)
		fmt.Println("GET")
		e := row.Scan(&id, &name, &email, &role)
		fmt.Println(e)
		if e != nil {
			return em, e
		}
		em = Employee{id, name, email, role}
	} else {
		return em, errors.New("error preparing select query")
	}
	return em, err
	/*  WITHOUT USING PREPARED STATEMENTS */
	// sel := `SELECT  * FROM employee2 where id = ` + fmt.Sprint(id)
	// fmt.Println(sel)
	// rows, err := db.Query(sel)
	// var em Employee
	// if err == nil {
	// 	for rows.Next() {
	// 		er := rows.Scan(&id, &name, &role, &email)
	// 		if er == nil {
	// 			em = Employee{id, name, role, email}
	// 		}
	// 	}
	// 	defer rows.Close()
	// 	return em, nil
	// } else {
	// 	fmt.Println(err)
	// 	return em, err
	// }
}

func (e *Employee) deleteEmp(db *sql.DB) (sql.Result, error) {

	var res sql.Result
	var er error
	if e.id == 0 {
		return res, errors.New("employee not found")
	}
	stmp, err := db.Prepare("DELETE FROM employee2 where id = ?")
	if err == nil {
		defer stmp.Close()
		res, er = stmp.Exec(e.id)
		if er != nil {
			return res, errors.New("failed executing delete query")
		}
	} else {
		return res, errors.New("error preparing delete query")
	}
	return res, er
}

func main() {
	// users := []user{{1, "Punit", "punitj1221@gmail.com", "Sde intern"}, {2, "Naman", "naman.garg@zopsmart.com", "sde intern"}}
	db, err := conn.Connect()
	name := "Punit"
	email := "punitjain1221@gmail.com"
	role := "SDE Intern"
	em := Employee{1, "Punit Jain", "punitjain1221@gmail.com", "SDE Intern"}
	if err == nil {
		// create(db)
		insertEmp(db, name, email, role)
		getEmpById(db, 1)
		em.updateEmp(db, []updates{{"email", "punitj1221@gmail.com"}, {"name", "Punit"}})
		_, er := em.deleteEmp(db)
		if er == nil {
			fmt.Println("Deleted Successfully")
		}
	}

}
