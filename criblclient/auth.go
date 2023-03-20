package criblclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type AuthResponse struct {
	Token               string `json:"token"`
	ForcePasswordChange bool   `json:"forcePasswordChange`
}
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) AuthLogin() (*AuthResponse, error) {
	body, err := json.Marshal(c.Auth)
	if err != nil {
		log.Fatalln("Error Occured while Marshalling the body : ", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/auth/login", c.Host), strings.NewReader(string(body)))
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	ar := AuthResponse{}
	err = json.Unmarshal(res, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
