package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run tools/hash_password.go <password>")
	}

	pw := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(hash))
}
