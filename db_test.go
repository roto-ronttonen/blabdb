package blabdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	db, err := Open("test.db")
	if err != nil {
		t.Logf("Error opening database")

	}

	// Add couple of collections
	data := map[string]string{
		"asd": "asd",
	}
	col1 := db.Collection("col1")
	col2 := db.Collection("col2")

	col1Key, err := col1.Insert(data)
	if err != nil {
		t.Logf("Failed creating collection1")
	}
	col2Key, err := col2.Insert(data)

	if err != nil {
		t.Logf("Failed creating collection2")
	}

	// List collections
	collections, err := db.GetAllCollections()
	if err != nil {
		t.Logf("Failed to list collections")
	}

	expected := []string{"col1", "col2"}

	assert.EqualValues(t, expected, collections, "Collection names don't match")

	// Delete collections

	col1.DeleteByKey(col1Key)
	col2.DeleteByKey(col2Key)

	// List collections (should be empty)
	collections, err = db.GetAllCollections()
	if err != nil {
		t.Logf("Failed to list collections")
	}

	expected = []string{}

	assert.EqualValues(t, expected, collections, "Collection names don't match")

	err = db.Close()
	if err != nil {
		t.Logf("Error closing database")
	}
}
