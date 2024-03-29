package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lampctl/lampctl/db"
	"github.com/lampctl/lampctl/gpio"
	"github.com/lampctl/lampctl/hue"
	"github.com/lampctl/lampctl/registry"
	"github.com/lampctl/lampctl/sequencer"
	"github.com/lampctl/lampctl/server"
	"github.com/lampctl/lampctl/ws2811"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "lampctl",
		Usage: "HTTP interface for controlling lights",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				EnvVars: []string{"DEBUG"},
				Usage:   "enable debug mode",
			},
			&cli.StringFlag{
				Name:    "db-path",
				EnvVars: []string{"DB_PATH"},
				Usage:   "path to SQLite database",
			},
			&cli.StringFlag{
				Name:    "server-addr",
				Value:   ":http",
				EnvVars: []string{"SERVER_ADDR"},
				Usage:   "HTTP address to listen on",
			},
		},
		Commands: []*cli.Command{
			installCommand,
		},
		Action: func(c *cli.Context) error {

			// Create the database
			db, err := db.New(&db.Config{
				Path: c.String("db-path"),
			})
			if err != nil {
				return err
			}
			defer db.Close()

			// Create the registry
			r := registry.New()
			defer r.Close()

			// Add the currently-supported providers

			// GPIO
			g, err := gpio.New(&gpio.Config{
				DB: db,
			})
			if err != nil {
				return err
			}
			r.Register(g)

			// Hue
			h, err := hue.New(&hue.Config{
				DB: db,
			})
			if err != nil {
				return err
			}
			r.Register(h)

			// ws2811
			w, err := ws2811.New(&ws2811.Config{
				DB: db,
			})
			if err != nil {
				return err
			}
			r.Register(w)

			// Create the sequencer
			seq := sequencer.New(&sequencer.Config{
				Registry: r,
			})
			defer seq.Close()

			// Start up the server
			s, err := server.New(&server.Config{
				Addr:      c.String("server-addr"),
				Debug:     c.Bool("debug"),
				Registry:  r,
				Sequencer: seq,
			})
			if err != nil {
				return err
			}
			defer s.Close()

			// Wait for SIGINT or SIGTERM
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
	}
}
