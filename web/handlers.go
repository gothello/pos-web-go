package handlers

import (

	"net/http"
	"github.com/gorilla/mux"
	"github.com/gothello/pos-web-go/beer"
	"github.com/urfave/negroni"
)

func GetBeerHandlers(r *mux.Router, n *negroni.Negroni, s beer.UseCase) {
	n.Handle("/v1/beer", n.With(
		negroni.Wrap(getAllBeer(service)),
	)).Methods("GET", "OPTIONS")

	n.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(getBeer(service)),
	)).Methods("GET", "OPTIONS")

	n.Handle("/v1/beer", n.With(
		negroni.Wrap(storeBeer(service)),
	)).Methods("POST", "OPTIONS")

	n.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(updateBeer(service)),
	)).Methods("PUT", "OPTIONS")

	n.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(removeBeer(service)),
	)).Methods("DELETE", "OPTIONS")
}

func getAllBeer(service beer.UseCase) http.Handler{
	return http.HandlerFunc(w http.ResponseWriter, r *http.Request){
		
	}
}
