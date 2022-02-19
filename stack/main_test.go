package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	db := NewDB()

	db.Set("key", "value")
	db.Set("key1", "value")
	assert.Equal(t, "value", db.Get("key"))
	assert.Equal(t, 2, db.Count("value"))
	db.Delete("key")
	assert.Equal(t, "", db.Get("key"))
	assert.Equal(t, 1, db.Count("value"))
}

func TestCRUDSingleTransRollback(t *testing.T) {
	db := NewDB()

	db.Set("key", "value")
	db.Set("key1", "value")
	db.Begin()
	assert.Equal(t, "value", db.Get("key"))
	assert.Equal(t, 2, db.Count("value"))
	db.Delete("key")
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	assert.Equal(t, "value", db.Get("key1"))
	db.Set("key-1", "value1")
	assert.Equal(t, "value1", db.Get("key-1"))
	db.Rollback()
	assert.Equal(t, 2, db.Count("value"))
	assert.Equal(t, "", db.Get("key-1"))
}

func TestCRUDSingleTransCommit(t *testing.T) {
	db := NewDB()

	db.Set("key", "value")
	db.Set("key1", "value")
	db.Begin()
	assert.Equal(t, "value", db.Get("key"))
	assert.Equal(t, 2, db.Count("value"))
	db.Delete("key")
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	assert.Equal(t, "value", db.Get("key1"))
	db.Set("key-1", "value1")
	assert.Equal(t, "value1", db.Get("key-1"))
	db.Commit()
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "value1", db.Get("key-1"))
}

func TestCRUDChainTransactionCommit(t *testing.T) {
	db := NewDB()

	db.Set("key", "value")
	db.Set("key1", "value")
	db.Begin()
	assert.Equal(t, "value", db.Get("key"))
	assert.Equal(t, 2, db.Count("value"))
	db.Delete("key")
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	assert.Equal(t, "value", db.Get("key1"))
	db.Begin()
	db.Set("key-1", "value1")
	assert.Equal(t, "value1", db.Get("key-1"))
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	db.Commit()
	assert.Equal(t, "value1", db.Get("key-1"))
	db.Commit()
	assert.Equal(t, 1, db.Count("value"))
}

func TestCRUDChainTransactionRollback(t *testing.T) {
	db := NewDB()

	db.Set("key", "value")
	db.Set("key1", "value")
	db.Begin()
	assert.Equal(t, "value", db.Get("key"))
	assert.Equal(t, 2, db.Count("value"))
	db.Delete("key")
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	assert.Equal(t, "value", db.Get("key1"))
	db.Begin()
	db.Set("key-1", "value1")
	assert.Equal(t, "value1", db.Get("key-1"))
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	db.Rollback()
	assert.Equal(t, "", db.Get("key-1"))
	db.Commit()
	assert.Equal(t, 1, db.Count("value"))
	assert.Equal(t, "", db.Get("key"))
	assert.Equal(t, "value", db.Get("key1"))
}