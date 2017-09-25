package router

import (
	"github.com/lnquy/fugu/modules/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	idxData         []byte
	idxStatics      []string
	staticLinkRegex = `/static/[a-zA-Z0-9./_\-?]+`
)

func init() {
	var err error
	if idxData, err = ioutil.ReadFile("./templates/index.html"); err != nil {
		log.Fatalf("router: failed to read index page template: %s", err)
	}
	idxStatics = util.GetStaticLinks(idxData, staticLinkRegex)
}

func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	w.Write(idxData)
	p, ok := w.(http.Pusher) // HTTP/2 Push
	if ok {
		for _, v := range idxStatics {
			p.Push(v, nil)
		}
	}
}
