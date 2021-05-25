package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/urfave/cli/v2"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func main() {
	app := &cli.App{
		Name:  "sshdocker",
		Usage: "interactive ssh connection to container",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "container-name",
				Aliases:  []string{"c"},
				Usage:    "Target container",
				EnvVars:  []string{"CONTAINER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "User for authentication",
				EnvVars:  []string{"SSH_USER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Password for authentication",
				EnvVars:  []string{"SSH_PASSWORD"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "host-key",
				Aliases: []string{"k"},
				Usage:   "Host key `FILE`",
				EnvVars: []string{"HOST_KEY_FILE"},
			},
			&cli.StringFlag{
				Name:    "shell",
				Aliases: []string{"s"},
				Usage:   "Default shell",
				Value:   "/bin/sh",
				EnvVars: []string{"CONTAINER_SHELL"},
			},
			&cli.StringFlag{
				Name:    "port",
				Usage:   "Binding port",
				Value:   "2222",
				EnvVars: []string{"PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			containerName := c.String("container-name")
			user := c.String("user")
			password := c.String("password")
			hostKeyFile := c.String("host-key")
			shell := strings.Split(c.String("shell"), " ")
			port := c.String("port")

			options := []ssh.Option{
				ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
					return pass == password && ctx.User() == user
				}),
			}
			if hostKeyFile != "" {
				options = append(options, ssh.HostKeyFile(hostKeyFile))
			}

			ssh.Handle(func(s ssh.Session) {
				args := []string{"exec", "-it", containerName}
				args = append(args, shell...)

				cmd := exec.Command("docker", args...)

				ptyReq, winCh, isPty := s.Pty()
				if isPty {
					home, _ := os.UserHomeDir() // Workaround until docker client 20.10.5 is available for stable alpine image (https://github.com/docker/cli/pull/2934)
					cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term), fmt.Sprintf("HOME=%s", home))
					f, err := pty.Start(cmd)
					if err != nil {
						panic(err)
					}
					go func() {
						for win := range winCh {
							setWinsize(f, win.Width, win.Height)
						}
					}()
					go func() {
						io.Copy(f, s) // stdin
					}()
					io.Copy(s, f) // stdout
					cmd.Wait()
				} else {
					io.WriteString(s, "No PTY requested.\n")
					s.Exit(1)
				}
			})

			log.Printf("starting ssh server on port %v", port)
			err := ssh.ListenAndServe(":"+port, nil, options...)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
