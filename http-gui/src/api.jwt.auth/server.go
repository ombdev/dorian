package main

import (
    "dorian/http-gui/src/api.jwt.auth/routers"
    "dorian/http-gui/src/api.jwt.auth/settings"
    "github.com/codegangsta/negroni"
    "net/http"
)

func main() {
	settings.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":5000", n)
//    http.ListenAndServe(":5000", nil)
}
