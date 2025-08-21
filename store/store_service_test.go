package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStoreService = &StorageService{}

func init() {
	_store := InitializeStore()
	testStoreService = _store
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	initialLink := "https://docs.google.com/spreadsheets/d/16VDUCGNobAeqdh8JvwBLNwseyKfjKwOBze10LKqyZOo/edit?gid=0#gid=0"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"

	SaveUrlMapping(shortURL, initialLink, userUUId)

	assert.Equal(t, initialLink, RetrieveInitialUrl(shortURL))
}
