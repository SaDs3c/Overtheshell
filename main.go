package main

import (
	"fmt"

	"github.com/gliderlabs/ssh"
	"github.com/sads3c/overtheshell/server"
)

func main() {

	srv := ssh.Server{
		Addr:                   ":2220",
		Handler:                server.Handler(),
		PasswordHandler:        server.PasswordHandler(),
		SessionRequestCallback: server.SessionRequestCallback(),
	}

	fmt.Println("SSH server listening on port 2220")
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}

}
