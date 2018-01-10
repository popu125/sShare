package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter(api *ApiServe) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/count", api.serveCount).Methods("POST")
	s.HandleFunc("/new", api.newProc).Methods("POST")
	//s.HandleFunc("/check/{procid}", api.procCheck)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	return r
}
