package main

import (
	"Golang-WebAuthN/Models"
	"Golang-WebAuthN/Redis"
	"Golang-WebAuthN/Utils"
	"encoding/json"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

var web *webauthn.WebAuthn
var userMapping map[string]*Models.User
var redis Redis.Redis

func main() {
	var err error
	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "James",            // Display Name for your site
		RPID:          "localhost",        // name of website
		RPOrigin:      "http://localhost", // The origin URL for WebAuthn requests
	})
	userMapping = map[string]*Models.User{}

	if err != nil {
		log.Fatal(err)
	}

	redis.Client = Redis.NewClient()

	router := httprouter.New()

	router.GET("/index", index)
	router.GET("/beginRegister", BeginRegistration)
	router.POST("/finishRegister", FinishRegistration)

	http.ListenAndServe(":8080", router)
}

func index(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	contents, err := os.ReadFile("Views/index.html")
	if err != nil {
		panic(err) // or do something useful
	}
	writer.WriteHeader(200)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Write(contents)
}

func BeginRegistration(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	name := query.Get("name")
	user := Models.NewUser(name)

	options, sessionData, err := web.BeginRegistration(user)
	Utils.ErrorHandle(err)
	user.Session = *sessionData

	redis.Set(name, user)

	// return options
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(options)
	Utils.ErrorHandle(err)
}

func FinishRegistration(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	name := query.Get("name")

	var user Models.User
	err := redis.Get(name, &user)
	Utils.ErrorHandle(err)

	session := user.Session
	response, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	Utils.ErrorHandle(err)

	credential, err := web.CreateCredential(user, session, response)
	user.AddCredential(*credential)
	// If login was successful, handle next steps
	log.Println(credential.ID)

	w.WriteHeader(http.StatusOK)
}
