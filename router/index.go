package router

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	idxPage []byte
)

func init() {
	var err error
	if idxPage, err = ioutil.ReadFile("./templates/index.html"); err != nil {
		log.Fatalf("router: failed to read index page template: %s", err)
	}
}

func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	w.Write(idxPage)
}
