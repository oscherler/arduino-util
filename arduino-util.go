package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const DefaultPortRegex = "cu\\.usb(serial|modem)"

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
	fmt.Println("Available commands:")
	fmt.Println("  makefile:     Output a Makefile suitable for building a project with local libraries.")
	fmt.Println("  find-board:   Find a connected Arduino based on a regexp.")
	os.Exit( 1 )
}

func makefile( executable string, args []string ) {
	makefileCommand := flag.NewFlagSet( "makefile", flag.ExitOnError )

	makefileCommand.Parse( args )
	tpl := template.Must( template.New("makefile").Parse( makefileTemplate ) )

	tpl.Execute( os.Stdout, struct { Executable string; DefaultRegex string }{ executable, DefaultPortRegex } )
}

func findBoard( args []string ) {
	findBoardCommand := flag.NewFlagSet( "find-board", flag.ExitOnError )
	findBoardRegex := findBoardCommand.String(
		"regex",
		DefaultPortRegex,
		"regex that matches the name the board appears under in /dev/" )

	findBoardCommand.Parse( args )

	re, err := regexp.Compile( *findBoardRegex )
	check( err )

	entries, err := os.ReadDir("/dev")
	check( err )

	board, err := findOneMatchingDevice( entries, *re )
	check( err )

	fmt.Println( board )
}

func check( err error ) {
	if err != nil {
		fmt.Fprintln( os.Stderr, err )
		os.Exit( 1 )
	}
}

func findOneMatchingDevice( entries []fs.DirEntry, re regexp.Regexp ) ( string, error ) {

	var boards []string
	
	for _, entry := range( entries ) {
		if isDevice( entry ) && re.MatchString( entry.Name() ) {
			boards = append( boards, entry.Name() )
		}
	}

	if len( boards ) == 0 {
		message := fmt.Sprintf( "No device matching '%v' found in /dev/.", re.String() )

		return "", errors.New( message )
	} else if len( boards ) > 1 {
		message := fmt.Sprintf( "More than one device matching '%v' found in /dev/:\n  %v", re.String(), strings.Join( boards, "\n  " ) )

		return "", errors.New( message )
	}
	
	return "/dev/" + boards[0], nil
}

func isDevice( entry fs.DirEntry ) bool {
	return entry.Type() & fs.ModeDevice != 0	
}
