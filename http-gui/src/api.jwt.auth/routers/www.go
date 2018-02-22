package routers

import (
	"github.com/gorilla/mux"
        "net/http"
)

func SetPublicWWW(router *mux.Router) *mux.Router {
    
    fs := http.FileServer(http.Dir("/home/j2eeserver/maxima/res/public"))
    
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))    
    
    return router
}
