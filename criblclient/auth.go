package criblclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthLogin(auth *Auth, authUrl string) (string, error) {
	client := &http.Client{}
	body, err := json.Marshal(auth)
	if err != nil {
		log.Fatalln("Error Occured while Marshalling the body : ", err)
	}
	req, err := http.NewRequest("GET", authUrl+"/api/v1/auth/login", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resbody), nil
}
