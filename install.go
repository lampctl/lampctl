package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/urfave/cli/v2"
)

const systemdUnitFile = `[Unit]
Description=Lampctl
Requires=network.target

[Service]
ExecStart={{.path}} --db-path {{.db_path}}

[Install]
WantedBy=multi-user.target
`

var installCommand = &cli.Command{
	Name:  "install",
	Usage: "install the application as a local service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "db-path",
			Value:   "/etc/lampctl",
			EnvVars: []string{"DB_PATH"},
			Usage:   "path for storing SQLite database",
		},
	},
	Action: install,
}

func install(c *cli.Context) error {

	// Create the path that will be used for the database
	if err := os.MkdirAll(c.String("db-path"), 0755); err != nil {
		return err
	}

	// Determine the full path to the executable
	p, err := os.Executable()
	if err != nil {
		return err
	}

	// Compile the template
	t, err := template.New("").Parse(systemdUnitFile)
	if err != nil {
		return err
	}

	// Attempt to open the file
	f, err := os.Create("/etc/systemd/system/lampctl.service")
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the template
	t.Execute(f, map[string]interface{}{
		"path":    p,
		"db_path": c.String("db-path"),
	})

	fmt.Println("Service installed!")
	fmt.Println("")
	fmt.Println("To enable the service and start it, run:")
	fmt.Println("  systemctl enable lampctl")
	fmt.Println("  systemctl start lampctl")

	return nil
}
