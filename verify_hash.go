package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
    hash := "$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO"
    password := "password"
    
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
        fmt.Println("No Match:", err)
    } else {
        fmt.Println("Match")
    }
}
