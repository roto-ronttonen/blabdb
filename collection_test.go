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
	Key   string `json:"key"`
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

	i := testDocument(item)
	i.Key = key

	assert.EqualValues(t, i, val, "Saved values not match")

	// Add second document

	secondDocKey, err := collection.Insert(item)

	if err != nil {
		t.Logf("Error inserting second document")
	}

	// Update first document

	newItem := testDocument(val)
	newItem.Age = 16

	newItem.Name = "Matti"

	updateKey, err := collection.Update(key, newItem)

	assert.Equal(t, key, updateKey, "Inserted and updated kesy dont match")

	if err != nil {
		t.Logf("Error updating item")
	}

	// Query only first

	exp1 := Expression{a: "name", operator: "==", b: "Teppo"}
	exp2 := Expression{a: "name", operator: "==", b: "Matti"}

	q1 := Query{
		blocks: [][]Expression{{exp1}},
	}
	items, err := collection.Find(q1, &val)

	if err != nil {
		t.Logf("Error querying results")
	}

	assert.Len(t, items, 1, "Length should be 1 is %d", len(items))

	items, err = collection.Find(Query{
		skip: 1,
	}, &val)

	if err != nil {
		t.Logf("Error querying results")
	}

	assert.Len(t, items, 1, "Length should be 1 is %d", len(items))

	items, err = collection.Find(Query{
		limit: 1,
	}, &val)

	if err != nil {
		t.Logf("Error querying results")
	}

	assert.Len(t, items, 1, "Length should be 1 is %d", len(items))

	// Query both

	q2 := Query{
		blocks: [][]Expression{{exp1}, {exp2}},
	}

	items, err = collection.Find(q2, &val)

	if err != nil {
		t.Logf("Error querying results")
	}

	assert.Len(t, items, 2, "Length should be 2 is %d", len(items))

	items, err = collection.Find(Query{}, &val)

	if err != nil {
		t.Logf("Error querying results")
	}

	assert.Len(t, items, 2, "Length should be 2 is %d", len(items))

	// Check val is correct
	err = collection.GetByKey(key, &val)

	if err != nil {
		t.Logf("Error getting document (should exist)")
	}

	i = testDocument(newItem)
	i.Key = updateKey

	assert.EqualValues(t, newItem, val, "Saved values not match")

	// Check keys

	keys, err := collection.GetAllKeys()

	if err != nil {
		t.Logf("Error getting collection keys")
	}

	assert.Equal(t, keys[0], key, "Listed key not equal to item key")

	// Delete keys

	err = collection.DeleteByKey(key)
	if err != nil {
		t.Logf("Error deleting document (should exist)")
	}

	err = collection.DeleteByKey(secondDocKey)
	if err != nil {
		t.Logf("Erorr deleting second document")
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
