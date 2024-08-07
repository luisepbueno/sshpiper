package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tg123/sshpiper/libplugin"
	"github.com/urfave/cli/v2"
)

func main() {

	var hostSeparator string
	var portSeparator string

	libplugin.CreateAndRunPluginTemplate(&libplugin.PluginTemplate{
		Name:  "hostinusername",
		Usage: "sshpiperd hostinusername plugin, only password auth is supported",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host-separator",
				Usage:       "username and host separator",
				Required:    false,
				Value:       "+",
				Destination: &hostSeparator,
			},
			&cli.StringFlag{
				Name:        "port-separator",
				Usage:       "host and port separator",
				Required:    false,
				Value:       ":",
				Destination: &portSeparator,
			},
		},
		CreateConfig: func(c *cli.Context) (*libplugin.SshPiperPluginConfig, error) {
			return &libplugin.SshPiperPluginConfig{
				PasswordCallback: func(conn libplugin.ConnMetadata, password []byte) (*libplugin.Upstream, error) {
					// username must be in format user<host-separator>host
					u := strings.Split(conn.User(), hostSeparator)
					if len(u) != 2 {
						return nil, fmt.Errorf("invalid username: %s", conn.User())
					}
					username := u[0]
					host := u[1]

					// host might be in format host<port-separator>port
					port := int32(22)
					h := strings.Split(host, portSeparator)
					if len(h) == 2 {
						host = h[0]
						_port, err := strconv.ParseInt(h[1], 10, 16)
						port = int32(_port)
						if err != nil {
							return nil, err
						}
					}

					if username == "" || host == "" {
						return nil, fmt.Errorf("invalid username: %s", conn.User())
					}

					return &libplugin.Upstream{
						Host:          host,
						Port:          port,
						UserName:      username,
						IgnoreHostKey: true,
						Auth:          libplugin.CreatePasswordAuth(password),
					}, nil
				},
			}, nil
		},
	})
}
