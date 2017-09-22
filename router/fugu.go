package router

import (
	"github.com/go-chi/chi"
	"github.com/lnquy/fugu/languages"
	"github.com/lnquy/fugu/modules/global"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func CalcSizeOfStruct(w http.ResponseWriter, r *http.Request) {
	var lang global.Language
	parm := chi.URLParam(r, "lang")
	if parm != "" {
		lang = global.LanguageEnum(parm)
	} else {
		lang = global.Go
	}

	var arch global.Architecture
	parm = chi.URLParam(r, "arch")
	if parm != "" {
		arch = global.ArchitectureEnum(parm)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Info(string(data))
	ret, err := languages.CalcSizeOf(string(data), lang, arch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}
