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

func formatJson(w hhtp.ResponseWriter, message interface{}, status int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(message)
}

func formatError(w http.ResponseWriter, err error, code int){
	formatJson(w, map[string]string{"error": err.String()})
}

func getAllBeer(service beer.UseCase) http.Handler{
	return http.HandlerFunc(w http.ResponseWriter, r *http.Request){
		all, err := service.GetBeer()
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formatJson(w, all, http.StatusOK)
	}
}
