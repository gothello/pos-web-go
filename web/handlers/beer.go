package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gothello/pos-web-go/core/beer"
	"github.com/urfave/negroni"
)

func GetBeerHandlers(r *mux.Router, n *negroni.Negroni, service beer.UseCase) {
	r.Handle("/v1/beer", n.With(
		negroni.Wrap(getAllBeer(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(getBeer(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/beer", n.With(
		negroni.Wrap(storeBeer(service)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(updateBeer(service)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(removeBeer(service)),
	)).Methods("DELETE", "OPTIONS")
}

func formatJson(w http.ResponseWriter, message interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(message)
}

func formatError(w http.ResponseWriter, messageErr string, code int) {
	formatJson(w, map[string]string{"error": messageErr}, code)
}

func getBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r.Body)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			formatError(w, err.Error(), http.StatusBadRequest)
			return
		}

		b, err := service.Get(id)
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formatJson(w, b, http.StatusOK)
	})
}

func getAllBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Accept") {
		case "application/json":
			getAllBeerJSON(w, service)
		default:
			getAllBeerHTML(w, service)
		}
	})
}

func getAllBeerJSON(w http.ResponseWriter, service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		all, err := service.GetAll()
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formatJson(w, all, http.StatusOK)
	})
}

func getAllBeerHTML(w http.ResponseWriter, service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		ts, err := template.ParseFiles(
			"./web/templates/header.html",
			"/web/templates/index.html",
			"/web/templates/footer.html",
		)

		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		all, err := service.GetAll()
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Beers []*beer.Beer
		}{
			Title: "Beers",
			Beers: all,
		}

		err = ts.Lookup("index.html").ExecuteTemplate(w, "index", data)
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func storeBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b beer.Beer

		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			formatError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.Store(&b)
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formatError(w, "ok", http.StatusCreated)
	})
}

func updateBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b beer.Beer

		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			formatError(w, err.Error(), http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r.Body)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			formatError(w, err.Error(), http.StatusBadRequest)
			return
		}

		b.ID = id

		err = service.Update(&b)
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func removeBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r.Body)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			formatError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.Remove(id)
		if err != nil {
			formatError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
