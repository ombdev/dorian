package main

import (
	"net/http"

	"dorian/http-gui/src/api.jwt.auth/routers"
	"dorian/http-gui/src/api.jwt.auth/settings"
	"github.com/urfave/negroni"
)

func main() {
	settings.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":5000", n)
}
