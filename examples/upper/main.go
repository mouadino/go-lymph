package main

import (
	"strings"

	"github.com/mouadino/go-lymph"
)

type Upper struct{}

func (Upper) Upper(text string) string {
	return strings.ToUpper(text)
}

func main() {
	server := lymph.Server(Upper{})

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
