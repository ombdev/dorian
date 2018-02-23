package routers

import (
	"github.com/gorilla/mux"
        "net/http"
)

func SetPublicWWW(router *mux.Router) *mux.Router {
    
    fs := http.FileServer(http.Dir("/home/dorian/go/src/dorian/resources/public"))
    
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))    
    
    return router
}
