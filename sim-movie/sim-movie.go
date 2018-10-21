/*
Contains the main package for the sim-movie programm. sim-movie creates a movie from snapshots
from the the barnes hut simulation
*/
package main

import "image"
import "image/gif"
import "os"
import "strings"
import "fmt"

func main() {

	// open the current working directory
	cwd, error := os.Open("../sim_server")

	if error != nil {
		panic("not able to open current working directory")
	}

	// get files with their names
	names, err := cwd.Readdirnames(0)

	if err != nil {
		panic("cannot read names in current working directory")
	}

	// parse the list of names and pick the ones that match the
	var files []string

	for _, dirname := range names {
		if strings.Contains(dirname, "hti") && strings.Contains(dirname, "gif") {
			fmt.Printf("dirname %s\n", dirname)
			files = append(files, dirname)
		}
	}

	// load static image and construct outGif
	outGif := &gif.GIF{}
	for _, name := range files {
		f, _ := os.Open("../sim_server/" + name)
		inGif, _ := gif.Decode(f)
		f.Close()

		outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
		outGif.Delay = append(outGif.Delay, 100)
	}

	// save to out.gif
	f, _ := os.OpenFile("hti.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
