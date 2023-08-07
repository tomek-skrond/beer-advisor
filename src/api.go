package main

import (
	"beercli"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenPort string
}

func NewAPIServer(port string) *APIServer {
	return &APIServer{
		ListenPort: port,
	}
}

func (s *APIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/", MakeHTTPHandler(s.Home)).Methods("GET")
	r.HandleFunc("/randombeer", MakeHTTPHandler(s.RandomBeer)).Methods("GET")

	fmt.Printf("Server listening at port %s \n", s.ListenPort)

	if err := http.ListenAndServe(s.ListenPort, r); err != nil {
		panic(err)
	}

}

func (s *APIServer) Home(w http.ResponseWriter, r *http.Request) error {
	path := os.Getenv("TEMPLATE_PATH")
	t, err := template.ParseFiles(path + "index.html")
	if err != nil {
		fmt.Println("error in home template")
		return err
	}

	if err := HandleTemplate(w, http.StatusOK, t, nil); err != nil {
		return err
	}

	return nil
}

func (s *APIServer) RandomBeer(w http.ResponseWriter, r *http.Request) error {
	path := os.Getenv("TEMPLATE_PATH")
	t, err := template.ParseFiles(path + "randombeer.html")
	if err != nil {
		fmt.Println("error in home template")
		return err
	}

	randomBeer, err := beercli.GetRandomBeer()
	if err != nil {
		return err
	}

	if err := HandleTemplate(w, http.StatusOK, t, randomBeer); err != nil {
		return err
	}

	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(apiError); ok {
				WriteJSON(w, e.Status, e)
				return
			}
			WriteJSON(w, http.StatusInternalServerError, "internal server error")
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	fmt.Printf("%v, Status: %v\n", time.Now(), status)
	return json.NewEncoder(w).Encode(v)
}

func HandleTemplate(w http.ResponseWriter, status int, t *template.Template, data any) error {
	fmt.Printf("%v, Status: %v\n", time.Now(), status)
	if err := t.Execute(w, data); err != nil {
		return err
	}
	return nil
}
