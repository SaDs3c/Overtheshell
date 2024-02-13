package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/gliderlabs/ssh"
)

func Handler() func(s ssh.Session) {
	return func(s ssh.Session) {
		ptyWindows(s)
	}

}

func PasswordHandler() ssh.PasswordHandler {
	return func(ctx ssh.Context, password string) bool {
		return password == "bandit0"
	}
}

func SessionRequestCallback() ssh.SessionRequestCallback {
	return func(sess ssh.Session, requestType string) bool {
		return sess.User() == "bandit0"
	}
}

func ptyWindows(s ssh.Session) {
	_, _, pty := s.Pty()

	if pty {
		os.Chdir("flag")
		cmd := exec.Command("powershell")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			return
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			io.Copy(stdin, s)
		}()

		go func() {
			io.Copy(s, stdout)
		}()

		go func() {
			io.Copy(s, stderr)
		}()

		err = cmd.Run()
		if err == nil {
			log.Println("session ended normally")
			s.Exit(0)
		} else {
			log.Println("session ended with an error:", err)
			exitCode := 1
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
				log.Println("Exit Code:", exitCode)
			}

			s.Exit(exitCode)
		}
	} else {
		io.WriteString(s, "No PTY requested. \n")
		s.Exit(1)
	}
}
