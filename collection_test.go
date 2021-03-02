package blabdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type address struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

type testDocument struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Age   int    `json:"age"`
	address
	LuckyNumbers []int `json:"luckyNumbers"`
}

func TestCollection(t *testing.T) {
	db, err := Open("test.db")
	if err != nil {
		t.Logf("Error opening database")

	}

	defer db.Close()

	collection := db.Collection("test")

	address := address{
		City:   "Turku",
		Street: "A street",
	}

	nums := []int{1, 2, 3}

	item := testDocument{
		Name:         "Teppo",
		Phone:        "0401234567",
		Age:          15,
		address:      address,
		LuckyNumbers: nums,
	}

	// Insert item
	key, err := collection.Insert(item)
	if err != nil {
		t.Logf("Error inserting document")
	}
	// Check val is correct
	var val testDocument

	err = collection.GetByKey(key, &val)

	if err != nil {
		t.Logf("Error getting document (should exist)")
	}

	assert.EqualValues(t, item, val, "Saved values not match")

	// Update item

	newItem := testDocument(val)
	newItem.Age = 16

	newItem.Name = "Matti"

	updateKey, err := collection.Update(key, newItem)

	assert.Equal(t, key, updateKey, "Inserted and updated kesy dont match")

	if err != nil {
		t.Logf("Error updating item")
	}

	// Check val is correct
	err = collection.GetByKey(key, &val)

	if err != nil {
		t.Logf("Error getting document (should exist)")
	}

	assert.EqualValues(t, newItem, val, "Saved values not match")

	// Check keys

	keys, err := collection.GetAllKeys()

	if err != nil {
		t.Logf("Error getting collection keys")
	}

	assert.Equal(t, keys[0], key, "Listed key not equal to item key")

	// Delete key

	err = collection.DeleteByKey(key)
	if err != nil {
		t.Logf("Error deleting document (should exist)")
	}

	// Check not found
	err = collection.GetByKey(key, &val)

	if err == nil {
		t.Logf("Should return not found error")
	}

	// Check update not found
	_, err = collection.Update(key, newItem)

	if err == nil {
		t.Logf("Should return not found error")
	}
}
