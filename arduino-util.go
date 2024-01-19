package main

import (
	"flag"
	"strings"
	"fmt"
	"regexp"
	"os"
	"io/fs"
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
	findBoardRegex := findBoardCommand.String(
		"regex",
		"cu\\.usb(serial|modem)",
		"regex that matches the name the board appears under in /dev/" )

	findBoardCommand.Parse( args )
	c, err := os.ReadDir("/dev")

	if err != nil {
		fmt.Println( err )
		os.Exit( 1 )
	}

	re, err := regexp.Compile( *findBoardRegex )

	if err != nil {
		fmt.Println( err )
		os.Exit( 1 )
	}

	var boards []string
	
	for _, entry := range( c ) {
		if isDevice( entry ) && re.MatchString( entry.Name() ) {
			boards = append( boards, entry.Name() )
		}
	}

	if len( boards ) == 0 {
		fmt.Fprintf( os.Stderr, "No device matching '%s' found in /dev/.\n", *findBoardRegex )
		os.Exit( 1 )
	} else if len( boards ) > 1 {
		fmt.Fprintf( os.Stderr, "More than one device matching '%s' found in /dev/:\n  %v\n", *findBoardRegex, strings.Join( boards, "\n  " ) )
		os.Exit( 1 )
	} else {
		fmt.Println( boards[0] )
	}
}

func isDevice( entry fs.DirEntry ) bool {
	return entry.Type() & fs.ModeDevice != 0	
}
