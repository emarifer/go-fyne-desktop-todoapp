package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"

	"github.com/ostafen/clover/v2/query"

	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
)

const (
	DB_NAME         = ".clover-db"
	TODO_COLLECTION = "todos"
)

type Db struct {
	db *c.DB
}

func MakeDb() Db {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("🔥 user directory not available: %s\n", err.Error())
	}
	dbAddress := filepath.Join(homeDir, DB_NAME)
	// It is necessary to create the folder that
	// will contain the DB storage files
	err = os.MkdirAll(dbAddress, os.ModePerm)
	if err != nil {
		log.Fatalf(
			"🔥 an error occurred while creating the DB directory: %s\n",
			err.Error(),
		)
	}
	db, err := c.Open(dbAddress)
	if err != nil {
		log.Fatalf("🔥 failed to connect to the DB: %s\n", err.Error())
	}

	db.CreateCollection(TODO_COLLECTION)

	return Db{db}
}

func (db *Db) Close() {
	err := db.db.Close()
	if err != nil {
		log.Fatalf("🔥 failed to close the connection DB: %s\n", err.Error())
	}
}

func formatTodo(doc *d.Document) models.Todo {
	t := models.Todo{}
	doc.Unmarshal(&t)
	t.Id = doc.ObjectId()

	return t
}

func (db *Db) InsertTodo(todo *models.Todo) bool {
	doc := d.NewDocumentOf(todo.ToMap())
	id, err := db.db.InsertOne(TODO_COLLECTION, doc)
	todo.Id = id

	return err == nil
}

func (db *Db) GetAllTodos() []models.Todo {
	docs, err := db.db.FindAll(query.NewQuery(TODO_COLLECTION).Sort(query.SortOption{Field: "created_at", Direction: 1}))
	if err != nil {
		log.Fatalln(err)
	}

	result := []models.Todo{}
	for _, doc := range docs {
		result = append(result, formatTodo(doc))
	}

	return result
}

func (db *Db) UpdateTodo(todo *models.Todo) bool {
	updates := todo.ToMap()

	// creating the query when the Id field (whose default name is "_id")
	// has the value of the todo that we pass to it
	q := query.NewQuery(TODO_COLLECTION).Where(query.Field("_id").Eq(todo.Id))
	err := db.db.Update(q, updates)

	return err == nil
}

func (db *Db) Save(todos []*models.Todo) bool {
	for _, t := range todos {
		if t.Id == "" {
			db.InsertTodo(t)
		} else {
			db.UpdateTodo(t)
		}
	}

	return true
}

func (db *Db) DeleteTodo(todo *models.Todo) bool {
	err := db.db.DeleteById(TODO_COLLECTION, todo.Id)

	return err == nil
}

func (db *Db) Drop() bool {
	err := db.db.DropCollection(TODO_COLLECTION)

	return err == nil
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
		homeDir, fmt.Sprintf("%s.json", TODO_COLLECTION),
	)

	err = db.db.ExportCollection(TODO_COLLECTION, exportResult)

	return err == nil
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
	importJSON := filepath.Join(
		homeDir, fmt.Sprintf("%s.json", TODO_COLLECTION),
	)

	err = db.db.ImportCollection(TODO_COLLECTION, importJSON)

	return err == nil
}

/*
UPDATE A COLLECTION ITEM GIVEN ITS ID:
https://github.com/ostafen/clover/blob/v2/examples/update/main.go#L32
*/
