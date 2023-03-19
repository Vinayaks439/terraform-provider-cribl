package criblclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

var authUrl = DefaultRestURL + "/auth"

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) AuthLogin(auth *Auth) (string, error) {
	body, err := json.Marshal(auth)
	if err != nil {
		log.Fatalln("Error Occured while Marshalling the body : ", err)
	}
	req, err := http.NewRequest("GET", authUrl+"/login", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	return base64.StdEncoding.EncodeToString(res), nil
}
