package router

import (
	"github.com/go-chi/chi"
	"github.com/lnquy/fugu/languages"
	"github.com/lnquy/fugu/modules/global"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/lnquy/fugu/languages/base"
)

func CalcSizeOfStruct(w http.ResponseWriter, r *http.Request) {
	lang, arch, data, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := languages.CalcSizeOf(string(data), lang, arch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}

func OptimizeStruct(w http.ResponseWriter, r *http.Request) {
	lang, arch, data, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := &base.Struct{}
	if err = json.Unmarshal(data, s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := languages.OptimizeMem(s, lang, arch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}

func getParams(r *http.Request) (global.Language, global.Architecture, []byte, error) {
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
	} else {
		arch = global.Amd64
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, 0, nil, err
	}
	return lang, arch, data, nil
}
