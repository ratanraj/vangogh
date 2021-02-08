package main

import (
	"flag"
	"fmt"
	"github.com/ratanraj/vangogh/cmd/vangoghcli/api"
)

func printUsage(command ...string) {
	fmt.Println("usage: vangoghcli [options] <command> [<subcommand>....]")
	fmt.Println("To see help text, you can run:")
	fmt.Println()
	fmt.Println("  vangoghcli help")
	fmt.Println("  vangoghcli <command> help")

	if len(command) == 0 {
		fmt.Println("Error: the following arguments are required: command")
	}
}

func main() {
	username := flag.String("user", "username", "username")
	password := flag.String("pass", "password", "password")
	albumName := flag.String("album", "", "New Album Name")
	albumID := flag.Uint("album-id", 0, "Album ID")
	//photoID := flag.Uint("photo-id", 0, "Photo ID")

	flag.Parse()

	client := api.NewAPI("http://127.0.0.1:8080/")

	if len(flag.Args()) < 1 {
		printUsage(flag.Args()...)
	}

	command := flag.Args()[0]

	switch command {
	case "login":
		_ = client.DoLogin(*username, *password)
		break
	case "album":
		subcommand := flag.Args()[1]
		switch subcommand {
		case "list":
			err := client.ListAlbums()
			if err != nil {
				panic(err)
			}
			break
		case "create":
			if len(*albumName) == 0 {
				panic(fmt.Errorf("album name required"))
			}
			err := client.CreateAlbum(*albumName)
			if err != nil {
				panic(err)
			}
			break
		case "delete":
			if *albumID == 0 {
				panic(fmt.Errorf("album ID required"))
			}
			err := client.DeleteAlbum(*albumID)
			if err != nil {
				panic(err)
			}
			break
		}
		break
	case "photo":
		subcommand := flag.Args()[1]
		switch subcommand {
		case "list":
			fmt.Println("")
			if *albumID == 0 {
				panic(fmt.Errorf("album ID required"))
			}
			err := client.ListPhotos(*albumID)
			if err != nil {
				panic(err)
			}
			break
		case "upload":
			fileName := flag.Args()[2]
			if *albumID == 0 {
				panic(fmt.Errorf("album ID required"))
			}
			err := client.UploadPhoto(*albumID, fileName)
			if err != nil {
				panic(err)
			}
			break
		}
		break
	default:
		printUsage(flag.Args()...)
	}
}
