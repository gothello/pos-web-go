package beer_test

import (
	"database/sql"
	"testing"

	"github.com/gothello/pos-web-go/beer"
)

func TestStore(t *testing.T) {

	b := &beer.Beer{
		ID:    1,
		Name:  "cerveja 1",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}

	db, err := sql.Open("sqlite3", "../data/beer.db")
	if err != nil {
		t.Fatalf("Erro conectando ao banco de dados: %s", err.Error())
	}

	defer db.Close()

	service := beer.NewService(db)

	err = service.Store(b)
	if err != nil {
		t.Fatalf("Erro salvando no banco de dados: %s", err.Error())
	}

	saved, err := service.Get(1)
	if err != nil {
		t.Fatalf("Erro buscando no banco de dados: %s", err.Error())
	}

	if saved.ID != 1 {
		t.Fatalf("Dados invalidos. Esperado %d, recebido %d", 1, saved.ID)
	}
}
