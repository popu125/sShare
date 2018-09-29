package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/popu125/sShare/web/safe"
)

func GetRouter(api *ApiServe) *mux.Router {
	r := mux.NewRouter()

	if api.conf.Safe.CityCheck {
		cityMgr := safe.NewCityLimit(api.conf.Safe.CityFile)
		r.Use(cityMgr.Middleware)
	}

	if api.conf.Safe.AntiCC {
		anticcMgr := safe.NewAntiCC(128)
		r.Use(anticcMgr.Middleware)
		r.HandleFunc("/get_cctoken", anticcMgr.Redirect).Methods("GET")
	}

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/count", api.serveCount).Methods("POST")
	s.HandleFunc("/new", api.newProc).Methods("POST")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	return r
}
