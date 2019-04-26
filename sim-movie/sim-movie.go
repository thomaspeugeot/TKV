/*
Contains the main package for the sim-movie programm. sim-movie creates a movie from snapshots
from the the barnes hut simulation
*/
package main

import "flag"
import "image"
import "image/gif"
import "os"
import "strings"
import "fmt"

func main() {

	dirPtr := flag.String("dir", "../sim_server", "directory where the snapshop are located and the movie will be generated")

	flag.Parse()

	// open the current working directory
	cwd, error := os.Open(*dirPtr)

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
	var movieFileName string

	for _, dirname := range names {
		if strings.Contains(dirname, ".gif") && strings.Contains(dirname, "conf-") && len(dirname) > 15 {

			// get 17 first caracters of the dirname
			sliptedName := strings.Split(dirname, "-")
			movieFileName = "movie-" + sliptedName[1] + "-" + sliptedName[2]

			fmt.Printf("dirname %s\n", dirname)
			files = append(files, dirname)
		}
	}

	fmt.Printf("movie name %s\n", movieFileName)

	// load static image and construct outGif
	outGif := &gif.GIF{}
	for _, name := range files {
		fileName := cwd.Name() + "/" + name
		f, _ := os.Open(fileName)
		inGif, _ := gif.Decode(f)
		f.Close()

		outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
		outGif.Delay = append(outGif.Delay, 5)
	}
	outGif.LoopCount = -1

	// save to out.gif
	outputFilename := cwd.Name() + "/" + movieFileName + ".gif"
	f, _ := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
