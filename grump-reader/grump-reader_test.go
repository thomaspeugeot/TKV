package grump_reader

import (
	"image"
	"image/color"
	"image/gif"

	"testing"
	"bufio"
	"fmt"
	"os"
	"log"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)


func BenchmarkReadGrumpNationalities(b * testing.B ) {

	var grumpNtlFile *os.File
	var err error
	grumpNtlFile, err = os.Open("/Users/thomaspeugeot/the-mapping-data/gl_grumpv1_ntlbndid_ascii_30/gluntlbnds.asc")

	popMapFile, _ := os.Create("popMapFile.gif")
	anim := gif.GIF{LoopCount: 1}
	rect := image.Rect(0, 0, 4320, 1720)
	img := image.NewPaletted(rect, palette)

	if err != nil {
		log.Fatal(err)
	}	

	var nbWords int
	scanner := bufio.NewScanner( grumpNtlFile)
	scanner.Split(bufio.ScanWords)
	
	// scan the header
	for nbWords < 12 {
		scanner.Scan()
		nbWords++		
		fmt.Println( fmt.Sprintf("item %d : %s", nbWords, scanner.Text()))
	}

	
	// Count the words, teh countries.
	countries := make(map[int]int)

	topLat := 85.0
	bottomLat := -85.0
	westLong := -180.0
	eastLong := 180.0

	lineLatWidth := (topLat - bottomLat) / 17160.0
	columnLongWdth := (eastLong - westLong) / 43200.0
	
	// for nbLines :=0; nbLines < 17160;  {
	parisIsMet := false
	for nbLines :=0; nbLines < 1710/5;  {
		lineLat := topLat - (float64(nbLines)*lineLatWidth)
		for nbWords =0; nbWords < 43200; scanner.Scan() {
			nbWords++
			columnLong := westLong + (float64(nbWords)*columnLongWdth)
			var value int
			fmt.Sscanf( scanner.Text(), "%d", &value)
			countries[value]++
			if( columnLong > 21.0 && lineLat < 48.5) {
				if parisIsMet == false {
					fmt.Printf("\nparis country code is %d\n", value)
				}
				parisIsMet = true
			}

			img.SetColorIndex(int(nbLines/10.0), int(nbWords/10.0),
				blackIndex)

			// fmt.Println(scanner.Text()) // Println will add back the final '\n'
		}
		nbLines++
		fmt.Printf("\rline %5d lat %2.3f", nbLines, lineLat)

	}
	for country, n := range countries {
		fmt.Printf("%d\t%d\n", country, n)
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}


	fmt.Println( fmt.Sprintf("nb lines %d", nbWords))


	gif.EncodeAll(popMapFile, &anim) // NOTE: ignoring encoding errors
	anim.Delay = append(anim.Delay, 8)
	anim.Image = append(anim.Image, img)
	popMapFile.Close()

}

