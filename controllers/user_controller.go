package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"mvc-inventary/models"
)

type UserController struct {
	Users []models.User
}

func (c *UserController) Load() {
	file, err := os.Open(models.UsersFileName)
	if os.IsNotExist(err) {
		c.Users = []models.User{}
		return
	}
	if err != nil {
		fmt.Println("Error opening users file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c.Users); err != nil {
		fmt.Println("Error decoding users file:", err)
		c.Users = []models.User{}
	}
}

func (c *UserController) Save() {
	file, err := os.Create(models.UsersFileName)
	if err != nil {
		fmt.Println("Error creating users file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c.Users); err != nil {
		fmt.Println("Error encoding users file:", err)
	}
}

func (c *UserController) FindByUsername(username string) *models.User {
	for i := range c.Users {
		if c.Users[i].Username == username {
			return &c.Users[i]
		}
	}
	return nil
}
