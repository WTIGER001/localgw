package main

import (
	"os"
	// "github.com/urfave/cli"
	// cli "github.com/urfave/cli/v2"
)

var swaggers []string

func main() {

	cmd := os.Args[1]
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-s" && len(os.Args) > i+1 {
			index := i + 1
			f := os.Args[index]
			swaggers = append(swaggers, f)
			i++
		}

	}

	switch cmd {
	case "generate":
		GenerateVSCodeLaunch(swaggers[0])
	case "serve":
		Serve(swaggers[0])
	}

	// commands = append(commands, generateCmd())
	// commands = append(commands, serveCmd())

	// app := &cli.App{
	// 	Commands: commands,
	// }

	// app.EnableBashCompletion = true
	// err := app.Run(os.Args)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// func serveCmd() *cli.Command {
// 	return &cli.Command{
// 		Name:  "serve",
// 		Usage: "Serves http proxy",
// 		Action: func(c *cli.Context) error {
// 			swagger := c.StringSlice("swagger")
// 			Serve(swagger)
// 			return nil
// 		},
// 		Flags: []cli.Flag{
// 			&cli.StringSliceFlag{
// 				Name:    "swagger",
// 				Aliases: []string{"s"},
// 			},
// 		},
// 	}
// }

// func generateCmd() *cli.Command {
// 	return &cli.Command{
// 		Name:  "generate",
// 		Usage: "List all themes",
// 		Action: func(c *cli.Context) error {
// 			swagger := c.StringSlice("swagger")

// 			for _, s := range swagger {
// 				GenerateVSCodeLaunch(s)
// 			}

// 			return nil
// 		},
// 		Flags: []cli.Flag{
// 			&cli.StringFlag{
// 				Name:    "out",
// 				Aliases: []string{"o"},
// 			},
// 			&cli.StringSliceFlag{
// 				Name:    "swagger",
// 				Aliases: []string{"s"},
// 			},
// 		},
// 	}
// }
