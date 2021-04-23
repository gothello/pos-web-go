package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gothello/pos-web-go/core/beer"
	"github.com/stretchr/testify/assert"
)

func Test_getAllBeer(t *testing.T) {
	b1 := &beer.Beer{
		ID:    10,
		Name:  "Heineken",
		Type:  beer.TypeAle,
		Style: beer.StyleRed,
	}

	b2 := &beer.Beer{
		ID:    20,
		Name:  "Skol",
		Type:  beer.StyleGolden,
		Style: beer.StyleRed,
	}

	db, err := sql.Open("sqlite3", "../../data/beer_test.db")
	assert.Nil(t, err)
	assert.Nil(t, clearDB(db))
	service := beer.NewService(db)
	assert.Nil(t, service.Store(b1))
	_ = service.Store(b2)

	service := beer.NewService(db)

	handler := getAllBeer(service)

	r := mux.NewRouter()
	r.Handle("/v1/beer", handler)
	req, err := http.NewRequest("GET", "/v1/beer", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []*beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, b1.ID, result[0].ID)
	assert.Equal(t, b2.ID, result[1].ID)
}

type BeerServiceMock struct{}

func (t BeerServiceMock) GetAll() ([]*beer.Beer, error) {
	b1 := &beer.Beer{
		ID:    10,
		Name:  "Heineken",
		Type:  beer.TypeAle,
		Style: beer.StyleRed,
	}

	b2 := &beer.Beer{
		ID:    20,
		Name:  "Skol",
		Type:  beer.TypeAle,
		Style: beer.StyleRed,
	}

	return []*beer.Beer{b1, b2}, nil
}

func (t BeerServiceMock) Get(ID int64) (*beer.Beer, error) {
	b1 := &beer.Beer{
		ID:    ID,
		Name:  "Heineken",
		Type:  beer.TypeAle,
		Style: beer.StylePale,
	}

	return b1, nil
}

func (t BeerServiceMock) Store(b *beer.Beer) error {
	return nil
}

func (t BeerServiceMock) Update(b *beer.Beer) error {
	return nil
}

func (t BeerServiceMock) Remove(ID int64) error {
	return nil
}

func Test_AllBeerWithMock(t *testing.T) {
	service := &BeerServiceMock{}
	handler := getAllBeer(service)

	r := mux.NewRouter()
	r.Handle("/v1/beer", handler)

	req, err := http.NewRequest("GET", "/v1/beer", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []*beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(10), result[0].ID)
	assert.Equal(t, int64(20), result[1].ID)
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
