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

const TODO_COLLECTION = "todos"

type Db struct {
	db *c.DB
}

func MakeDb(dbFiles string) Db {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("🔥 user directory not available: %s\n", err.Error())
	}
	// creating the address of the folder where the DB files
	// will be saved as a hidden folder in the user folder
	dbAddress := filepath.Join(homeDir, fmt.Sprintf(".%s", dbFiles))
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

	collectionExists, err := db.HasCollection(TODO_COLLECTION)
	if err != nil {
		log.Fatalf("🔥 failed to check collection: %s\n", err.Error())
	}

	if !collectionExists {
		db.CreateCollection(TODO_COLLECTION)
	}

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
		log.Printf("something went wrong: %s", err)
	}

	result := []models.Todo{}
	for _, doc := range docs {
		result = append(result, formatTodo(doc))
	}

	return result
}

func (db *Db) UpdateTodo(todo *models.Todo) bool {
	updates := todo.ToMap()
	// ↓ We delete the field that we do not want to update ↓
	delete(updates, "created_at")
	// creating the query when the Id field (whose default name is "_id")
	// has the value of the todo that we pass to it
	q := query.NewQuery(TODO_COLLECTION).Where(query.Field("_id").Eq(todo.Id))
	err := db.db.Update(q, updates)

	return err == nil
}

func (db *Db) DeleteTodo(todo *models.Todo) bool {
	err := db.db.DeleteById(TODO_COLLECTION, todo.Id)

	return err == nil
}

func (db *Db) Drop() bool {
	err := db.db.DropCollection(TODO_COLLECTION)
	if err != nil {
		return false
	}

	// After deleting the collection, if we want to save data again,
	// we have to create a new empty collection
	err = db.db.CreateCollection(TODO_COLLECTION)

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

	// We delete the old collection so that
	// the import generates a new `clean` collection
	err = db.db.DropCollection(TODO_COLLECTION)
	if err != nil {
		return err == nil
	}

	err = db.db.ImportCollection(TODO_COLLECTION, importJSON)

	return err == nil
}

/*
UPDATE A COLLECTION ITEM GIVEN ITS ID:
https://github.com/ostafen/clover/blob/v2/examples/update/main.go#L32
*/
