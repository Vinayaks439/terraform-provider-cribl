package cribl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const UsersURL = DefaultRestURL + "/system/users"

type GetUser struct {
	Username string   `json:"username"`
	First    string   `json:"first"`
	Last     string   `json:"last"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Id       string   `json:"id"`
	Disabled bool     `json:"disabled,omitempty"`
}

type CreateUser struct {
	Username string   `json:"username"`
	First    string   `json:"first"`
	Last     string   `json:"last"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Id       string   `json:"id"`
	Disabled bool     `json:"disabled"`
	Password string   `json:"password"`
}

type PatchUser struct {
	Username string   `json:"username"`
	First    string   `json:"first"`
	Last     string   `json:"last"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles,omitempty"`
	Id       string   `json:"id"`
	Disabled bool     `json:"disabled"`
	Password string   `json:"password,omitempty"`
}

type Users struct {
	Items []GetUser `json:"items"`
	Count int       `json:"count"`
}

func (c *Client) GetUsers() (Users, error) {
	req, err := http.NewRequest("GET", UsersURL, nil)
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	var users Users
	err = json.Unmarshal(res, &users)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling: ", err)
	}

	return users, nil
}

func (c *Client) CreateUser(user *CreateUser) (Users, error) {
	body, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Error Occured while Marshalling the body : ", err)
	}
	req, err := http.NewRequest("POST", UsersURL, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	var users Users
	err = json.Unmarshal(res, &users)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling: ", err)
	}

	return users, nil
}

func (c *Client) GetUserByID(id string) (Users, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(UsersURL+"/%s", id), nil)
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	var users Users
	err = json.Unmarshal(res, &users)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling: ", err)
	}

	return users, nil
}

func (c *Client) PatchUserByID(id string, user *PatchUser) (Users, error) {
	body, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Error Occured while Marshalling the body : ", err)
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf(UsersURL+"/%s", id), bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	var users Users
	err = json.Unmarshal(res, &users)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling: ", err)
	}

	return users, nil
}

func (c *Client) DeleteUserbyID(id string) (Users, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf(UsersURL+"/%s/info", id), nil)
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	var users Users
	err = json.Unmarshal(res, &users)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling: ", err)
	}

	return users, nil
}

func (c *Client) PatchUserInfo(id string, user *PatchUser) (Users, error) {
	body, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Error Occured while Marshalling the body : ", err)
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf(UsersURL+"/%s", id), bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Error Occured while initializing new request : ", err)
	}
	res, err := c.doRequest(req)
	if err != nil {
		log.Fatalln("Error Occured while doing the POST call : ", err)
	}
	var users Users
	err = json.Unmarshal(res, &users)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling: ", err)
	}

	return users, nil
}
