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

var palette = []color.Color{
	color.White, 
	color.Black, 
	color.RGBA{255,0,0,255}, // red
	color.RGBA{0,255,0,255}, // green
	color.RGBA{255, 255, 0, 255}, 
	color.RGBA{255, 118, 118, 255}, 
	color.RGBA{255, 38, 38, 255}, 
	color.RGBA{106, 106, 38, 255}, 
	color.RGBA{255, 95, 95, 255}, 
	color.RGBA{57, 57, 38, 255}, 
	color.RGBA{167, 167, 38, 255}, 
	}

const (
	whiteIndex = iota // first color in palette
	blackIndex // next color in palette
	redIndex // next color in palette
	greenIndex
	blueIndex // next color in palette
	index218
	index38
	index106
	index95
	index57
	index167
)


func BenchmarkReadGrumpNationalities(b * testing.B ) {

	var grumpNtlFile *os.File
	var err error
	grumpNtlFile, err = os.Open("/Users/thomaspeugeot/the-mapping-data/gl_grumpv1_ntlbndid_ascii_30/gluntlbnds.asc")
	if err != nil {
		log.Fatal(err)
	}	

	popMapFile, _ := os.Create("popMapFile.gif")
	anim := gif.GIF{LoopCount: 1}

	// for debug sake, introduce which ratio of the 
	// input date is processed
	ratioOfProcessedLines := 0.4
	maxNbLines := int( ratioOfProcessedLines * float64(17100))
	maxNbCols := 43200

	// prepare the output image
	displayRatioX := 0.1
	displayRatioY := 0.1
	maxImageX := int( float64( maxNbCols) * displayRatioX )
	maxImageY := int( float64( maxNbLines) * displayRatioY )
	rect := image.Rect(0, 0, maxImageX, maxImageY)
	img := image.NewPaletted(rect, palette)


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
	for nbLines :=0; nbLines < maxNbLines;  {
		lineLat := topLat - (float64(nbLines)*lineLatWidth)
		for nbWords =0; nbWords < maxNbCols; scanner.Scan() {
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

			var pixelIndex uint8
			pixelIndex = whiteIndex
			switch( value ) {
			case 35 : // canada
				pixelIndex = blackIndex
			case 177 : // russia
				pixelIndex = blueIndex
			case 185 :
				pixelIndex = redIndex
			case 155 :
				pixelIndex = greenIndex
			case 218 :
				pixelIndex = index218
			case 38 :
				pixelIndex = index38
			case 106 :
				pixelIndex = index106
			case 95 :
				pixelIndex = index95
			case 57 :
				pixelIndex = index57
			case 167 : // france
				pixelIndex = index167
			}

			if( value != -9999) {
				img.SetColorIndex(
					int( float64(nbWords) * displayRatioY), 
					int( float64(nbLines) * displayRatioY), 
					pixelIndex)
			}
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


	anim.Delay = append(anim.Delay, 8)
	anim.Image = append(anim.Image, img)
	gif.EncodeAll(popMapFile, &anim) // NOTE: ignoring encoding errors
	popMapFile.Close()

}

