package main

import (
	"strings"

	"github.com/mouadino/go-lymph"
	"github.com/mouadino/go-nano/discovery"
)

type Upper struct{}

func (Upper) Upper(text string) string {
	return strings.ToUpper(text)
}

func main() {
	server := lymph.Server(Upper{})
	// TODO: Repeating identity :(
	server.Announce("echo", discovery.ServiceMetadata{"identity": "..."}, lymph.ZookeeperAnnouncer)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
