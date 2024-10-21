package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/earthboundkid/csv/v2"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"
	"github.com/google/uuid"

	"github.com/joho/sqltocsv"

	_ "github.com/mattn/go-sqlite3"
)

const FTODO_TABLE_NAME = "ftodos"

type Db struct {
	db *sql.DB
}

func MakeDb(DbName string) Db {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("🔥 user directory not available: %s\n", err.Error())
	}
	// creating the address of the file where the DB
	// will be saved as a hidden file in the user folder
	dDAddress := filepath.Join(homeDir, fmt.Sprintf(".%s", DbName))

	// Init SQLite3 database
	db, err := sql.Open("sqlite3", dDAddress)
	if err != nil {
		log.Fatalf("🔥 failed to connect to the DB: %s\n", err)
	}

	sqlStr := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id TEXT NOT NULL PRIMARY KEY,
			description TEXT NOT NULL,
			done BOOLEAN DEFAULT(FALSE),
			created_at DATE DEFAULT (datetime('now','localtime'))
			);`, FTODO_TABLE_NAME,
	)

	_, err = db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("🔥 failed to create table: %s\n", err)
	}

	return Db{db}
}

func (db *Db) Close() {
	err := db.db.Close()
	if err != nil {
		log.Fatalf("🔥 failed to close the connection DB: %s\n", err)
	}
}

func (db *Db) InsertTodo(todo *models.Todo) (*time.Time, bool) {
	sqlStr := fmt.Sprintf(`INSERT INTO %s (id, description)
		VALUES(?, ?) RETURNING created_at`, FTODO_TABLE_NAME)

	stmt, err := db.db.Prepare(sqlStr)
	if err != nil {
		return nil, err == nil
	}

	defer stmt.Close()

	t := models.Todo{}

	err = stmt.QueryRow(todo.Id, todo.Description).Scan(&t.CreatedAt)
	if err != nil {
		return nil, err == nil
	}

	return &t.CreatedAt, true
}

func (db *Db) GetAllTodos() []models.Todo {
	todos := []models.Todo{}

	query := fmt.Sprintf("SELECT * FROM %s", FTODO_TABLE_NAME)

	rows, err := db.db.Query(query)
	if err != nil {
		return todos
	}
	// We close the resource
	defer rows.Close()

	t := models.Todo{}

	for rows.Next() {
		rows.Scan(&t.Id, &t.Description, &t.Done, &t.CreatedAt)

		todos = append(todos, t)
	}

	return todos
}

func (db *Db) UpdateTodo(todo *models.Todo) bool {
	query := fmt.Sprintf(`UPDATE %s SET done = ?
		WHERE id=?`, FTODO_TABLE_NAME)

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err == nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(todo.Done, todo.Id)
	if err != nil {
		return err == nil
	}

	return true
}

func (db *Db) DeleteTodo(todo *models.Todo) bool {

	query := fmt.Sprintf(`DELETE FROM %s
		WHERE id=?`, FTODO_TABLE_NAME)

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err == nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(todo.Id)
	if err != nil {
		return err == nil
	}

	return true
}

func (db *Db) Drop() bool {
	sqlStr := fmt.Sprintf(`DELETE FROM %s;`, FTODO_TABLE_NAME)

	_, err := db.db.Exec(sqlStr)
	if err != nil {
		return err == nil
	}

	return true
}

func (db *Db) ExportData() bool {
	var (
		homeDir string
		err     error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		return err == nil
	}
	exportResult := filepath.Join(
		homeDir, fmt.Sprintf("%s.csv", FTODO_TABLE_NAME),
	)

	query := fmt.Sprintf("SELECT * FROM %s", FTODO_TABLE_NAME)

	rows, err := db.db.Query(query)
	if err != nil {
		return err == nil
	}
	// We close the resource
	defer rows.Close()

	err = sqltocsv.WriteFile(exportResult, rows)
	if err != nil {
		return err == nil
	}

	return true
}

func (db *Db) ImportData() bool {
	var (
		homeDir string
		err     error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		return err == nil
	}
	importCSV := filepath.Join(
		homeDir, fmt.Sprintf("%s.csv", FTODO_TABLE_NAME),
	)

	// We delete the old table so that
	// the import generates a new `clean` table
	if ok := db.Drop(); !ok {
		return ok
	}

	data, err := os.ReadFile(importCSV)
	if err != nil {
		return err == nil
	}

	csvOpt := csv.Options{Reader: strings.NewReader(string(data))}
	rowsMap, err := csvOpt.ReadAll()
	if err != nil {
		return err == nil
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s(id, description, done, created_at) VALUES ", FTODO_TABLE_NAME)
	vals := []interface{}{}

	for _, row := range rowsMap {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(
			vals,
			uuid.NewString(),
			row["description"],
			convertToBool(row["done"]),
			convertToDatetime(row["created_at"]),
		)
	}
	// trim the last `,`
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	// prepare the statement
	stmt, err := db.db.Prepare(sqlStr)
	if err != nil {
		return err == nil
	}

	defer stmt.Close()

	// format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		return err == nil
	}

	return true
}

func convertToBool(s string) bool {
	boolValue, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal(err)
	}
	return boolValue
}

func convertToDatetime(s string) time.Time {
	layout := "2006-01-02 15:04:05 -0700 MST"
	date, err := time.Parse(layout, s)
	if err != nil {
		log.Fatal(err)
	}
	return date
}

/* REFERENCES:
How to create timestamp column with default value 'now'?:
https://stackoverflow.com/questions/200309/how-to-create-timestamp-column-with-default-value-now
Convert UTC string to time object:
https://stackoverflow.com/questions/38798043/convert-utc-string-to-time-object
https://dev.to/luthfisauqi17/golangs-unique-way-to-parse-string-to-time-2jmk
Querying for data:
https://go.dev/doc/database/querying
Go CSV reader like Python's DictReader:
https://github.com/earthboundkid/csv/blob/master/example_test.go
https://pkg.go.dev/github.com/carlmjohnson/csv#FieldReader.ReadAll

How to insert multiple data at once:
https://stackoverflow.com/questions/21108084/how-to-insert-multiple-data-at-once

https://stackoverflow.com/questions/69217606/bulk-insert-rows-from-an-array-to-an-sql-server-with-golang
https://golangbot.com/mysql-create-table-insert-row/

https://github.com/joho/sqltocsv

https://universalglue.dev/posts/csv-to-sqlite/
https://ahdeyy.hashnode.dev/converting-csv-data-to-sqlite-in-golang
https://github.com/Ahdeyyy/SpotifyTop2018

UPDATE A COLLECTION ITEM GIVEN ITS ID:
https://github.com/ostafen/clover/blob/v2/examples/update/main.go#L32
*/
