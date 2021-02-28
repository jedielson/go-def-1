package uweb

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func FromQuery(r *http.Request, key string, defaultValue int, f func(int) bool) int {

	off := r.URL.Query().Get(key)
	i, err := strconv.Atoi(off)
	if err == nil && f(i) {
		return i
	}

	return defaultValue
}

func FromPath(r *http.Request, key string, f func(int) bool, erro error) (int, error) {
	params := mux.Vars(r)
	value := params[key]
	p, err := strconv.Atoi(value)
	if err == nil && f(p) {
		return p, nil
	}

	return 0, erro
}

func ToJson(w http.ResponseWriter, i interface{}, e ...error) {

	w.Header().Set("Content-Type", "application/json")
	if len(e) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(i)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(bytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
