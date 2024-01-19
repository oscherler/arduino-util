package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	_ "embed"
	"text/template"
)

//go:embed Makefile.tmpl
var makefileTemplate string

func main() {

	if len( os.Args ) < 2 {
		usage()
	}
	
	switch os.Args[1] {
		case "makefile":
			makefile( filepath.Base( os.Args[0] ), os.Args[2:] )
		case "find-board":
			findBoard( os.Args[2:] )
		default:
			usage()
	}
}

func usage() {
	fmt.Println("Usage")
	os.Exit( 1 )
}

func makefile( executable string, args []string ) {
	makefileCommand := flag.NewFlagSet( "makefile", flag.ExitOnError )

	makefileCommand.Parse( args )
	tpl := template.Must( template.New("makefile").Parse( makefileTemplate ) )

	tpl.Execute( os.Stdout, struct { Executable string }{ executable } )
}

func findBoard( args []string ) {
	findBoardCommand := flag.NewFlagSet( "find-board", flag.ExitOnError )
	findBoardRegex := findBoardCommand.String( "regex", "cu\\.usb(serial|modem)", "regex that matches the name the board appears under in /dev/" )

	findBoardCommand.Parse( args )
	fmt.Printf( "find-board %v\n", *findBoardRegex )
}
