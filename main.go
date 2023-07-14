package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func main() {
	home := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		data := map[string]string{
			"wow": "mega",
		}
		tmpl.Execute(w, data)
	}
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleSignUp(w http.ResponseWriter, r *http.Request) error {
	var registerReq *SingUpReq
	json.NewDecoder(r.Body).Decode(&registerReq)

	return WriteResponse(w, 200, "succesfully register user")
}

func WriteResponse(w http.ResponseWriter, code int, data ...any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(&data)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteResponse(w, http.StatusBadRequest, "err: "+err.Error())
		}
	}
}
