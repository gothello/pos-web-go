package handlers

import (
	"database/sql"
	"testing"

	"github.com/gothello/pos-web-go/core/beer"
	"github.com/stretchr/testify/assert"
)

func Test_getAllBeer(t *testing.T) {
	b1 := &beer.Beer{
		ID:    10,
		Name:  "Heineken",
		Type:  beer.StyleGolden,
		Style: beer.StyleRed,
	}

	b2 := &beer.Beer{
		ID:    20,
		Name:  "Skol",
		Type:  beer.StyleBrown,
		Style: beer.StyleLime,
	}

	db, err := sql.Open("sqlite3", "../../data/beer_test.db")
	assert.Nil(t, err)
	assert.Nil(t, clearDB(db))
	service := beer.NewService(db)
	assert.Nil(t, service.Store(b1))
	assert.Nil(t, service.Store(b2))
}

func clearDB(db *sql.DB) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from beer")
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
