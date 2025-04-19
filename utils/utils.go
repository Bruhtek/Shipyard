package utils

import (
	"log"
	"math/rand"
)

func IFErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
		panic(err)
	}
}

const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
