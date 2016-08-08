package grump

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"testing"
	"bufio"
	"fmt"
	"os"
	"log"
	"github.com/thomaspeugeot/tkv/country"
	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/quadtree"
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

	// open grump file for ntl boundaries
	var grumpNtlFile *os.File
	var err error
	grumpNtlFile, err = os.Open("/Users/thomaspeugeot/the-mapping-data/gl_grumpv1_ntlbndid_ascii_30/gluntlbnds.asc")
	if err != nil {
		log.Fatal(err)
	}	

	// open grump files for densities
	var grumpNtlFileDensities *os.File
	var errDensities error
	grumpNtlFileDensities, errDensities = os.Open("/Users/thomaspeugeot/the-mapping-data/gl_grumpv1_pcount_00_ascii_30/glup00g.asc")
	if errDensities != nil {
		log.Fatal(err)
	}	


	// open output map
	popMapFile, _ := os.Create("popMapFile.gif")
	anim := gif.GIF{LoopCount: 1}

	// for debug sake, introduce which ratio of the 
	// input date is processed
	ratioOfProcessedLines := 0.45
	maxNbLines := int( ratioOfProcessedLines * float64(17100))
	maxNbCols := 43200

	// prepare the output image
	displayRatioX := 0.1
	displayRatioY := 0.1
	maxImageX := int( float64( maxNbCols) * displayRatioX )
	maxImageY := int( float64( maxNbLines) * displayRatioY )
	rect := image.Rect(0, 0, maxImageX, maxImageY)
	img := image.NewPaletted(rect, palette)

	// prepare borders of country
	type countryBoundary struct {
		topLat float64
		bottomLat float64
		westLng float64
		eastLng float64
	}

	// prepare the density file
	maxBodies := 200000
	bodyCountFrance := 0
	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, maxBodies)

	// prepare one file per country (for 1 to 229)
	countryGifFiles := make(map[int](* os.File))
	countryRectangle := make(map[int] image.Rectangle)
	countryPalette := make(map[int](* image.Paletted))
	countryAnim  := make(map[int]gif.GIF)
	countryName := make(map[int]string)
	countryBoundaries := make(map[int]*countryBoundary)
	countryFinalBoundaries := make(map[int]*countryBoundary)

	// parse countries and make a match between index and known country
	for index:= 0; index<=230; index++ {
		countryCodeGrump := country.CountryCodes[index]
		countryName[ countryCodeGrump.CodeGrump ] = countryCodeGrump.Name
				
		countryBorder := country.CountryBorders[index]
		countryFinalBoundaries[ countryBorder.Index] = & countryBoundary{ countryBorder.TopLat,
		countryBorder.BottomLat,
		countryBorder.EastLng,
		countryBorder.WestLng }
		
	}

	for index:= 0; index<=230; index++ {
	
		countryString, _ := countryName[ index]

		popMapCountryFile, _ := os.Create( fmt.Sprintf("popMapFile%03d-%s.gif", index, countryString))
		countryGifFiles[index] = popMapCountryFile
		countryRectangle[index] = image.Rect(0, 0, maxImageX, maxImageY)
		countryPalette[index] = image.NewPaletted(rect, palette)
		countryAnim[index] = gif.GIF{LoopCount: 1}
	
		countryBoundaries[index] = & countryBoundary{ -90.0,+90.0,180.0,-180.0}
		// prepare the output per country
	}

	// parse the grump
	var nbWords int
	scanner := bufio.NewScanner( grumpNtlFile)
	scanner.Split(bufio.ScanWords)
	
	// scan the header
	for nbWords < 12 {
		scanner.Scan()
		nbWords++		
		fmt.Println( fmt.Sprintf("item %d : %s", nbWords, scanner.Text()))
	}

	// parse the grump densities
	var nbWordsDensities int
	scannerD := bufio.NewScanner( grumpNtlFileDensities)
	scannerD.Split(bufio.ScanWords)
	
	// scan the header
	for nbWords < 12 {
		scanner.Scan()
		nbWords++		
		fmt.Println( fmt.Sprintf("item %d : %s", nbWords, scanner.Text()))
	}

	// scan the header
	for nbWordsDensities < 12 {
		scannerD.Scan()
		nbWordsDensities++		
		fmt.Println( fmt.Sprintf("densities item %d : %s", nbWordsDensities, scannerD.Text()))
	}
	
	// Count the words, teh countries.
	countries := make(map[int]int)
	
	// Count the words, teh countries.
	countriesDensities := make(map[int]int)

	topLat := 85.0
	bottomLat := -85.0
	westLong := -180.0
	eastLong := 180.0

	lineLatWidth := (topLat - bottomLat) / 17160.0
	columnLongWdth := (eastLong - westLong) / 43200.0
	
	// init france body count

	// for nbLines :=0; nbLines < 17160;  {
	parisIsMet := false
	for nbLines :=0; nbLines < maxNbLines;  {
		lineLat := topLat - (float64(nbLines)*lineLatWidth)
		for nbWords =0; nbWords < maxNbCols; scanner.Scan() {
			nbWords++
			columnLong := westLong + (float64(nbWords)*columnLongWdth)

			// get the scan value for boundary
			var value int
			fmt.Sscanf( scanner.Text(), "%d", &value)
			countries[value]++

			// get the scan value for count
			var count int
			scannerD.Scan()
			fmt.Sscanf( scanner.Text(), "%d", &count)
			countriesDensities[value] += count

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
				if( nbLines %100 == 0) {
					fmt.Printf("\rval %03d line %5d lat %2.3f ln %03.3f ", value, nbLines, lineLat, math.Abs(columnLong))
				}

				img.SetColorIndex(
					int( float64(nbWords) * displayRatioY), 
					int( float64(nbLines) * displayRatioY), 
					pixelIndex)
				if _, ok := countryPalette[value]; ok {
					countryPalette[value].SetColorIndex(
						int( float64(nbWords) * displayRatioY), 
						int( float64(nbLines) * displayRatioY), 
						blackIndex)
				}
				if boundary, ok := countryBoundaries[value]; ok {
					if boundary.topLat < lineLat {	boundary.topLat = lineLat }
					if boundary.bottomLat > lineLat { boundary.bottomLat = lineLat }
					if boundary.westLng > columnLong { boundary.westLng = columnLong}
					if boundary.eastLng < columnLong { boundary.eastLng = columnLong}
				}
			}
			// fmt.Println(scanner.Text()) // Println will add back the final '\n'


			// for france, fill up bodies count
			if( value == 67) {
				// compute relative position in the 1*1 sqaure
				border := country.CountryBorders[value]
				relX := (columnLong - border.WestLng) / (border.EastLng - border.WestLng)
				relY := (lineLat - border.BottomLat) / (border.TopLat - border.BottomLat)

				// how many bodies ? it is maxBodies *( count / border.PCount) 
				bodyNb := int( math.Floor( float64( maxBodies) * float64( count) / float64( border.PCount)))

				// initiate the bodies
				for i :=0; i<bodyNb && bodyCountFrance+i < maxBodies; i++ {
					bodies[bodyCountFrance+i].X = relX
					bodies[bodyCountFrance+i].Y = relY
				}
				bodyCountFrance += bodyNb
			}
		}
		nbLines++

	}

	// capture config
	var r barnes_hut.Run
	r.Init( & bodies)
	r.SetCountry("fra")
	r.CaptureConfig()

	// indicate the corners of each country
	for index, boundary := range countryBoundaries {

		nbLinesTop := math.Floor( (topLat - boundary.topLat) / lineLatWidth)
		nbLinesBottom := math.Floor( (topLat - boundary.bottomLat) / lineLatWidth)
		nbWordsWest := math.Floor( (boundary.westLng - westLong) / columnLongWdth)
		nbWordsEast := math.Floor( (boundary.eastLng - westLong) / columnLongWdth)

		countryPalette[index].SetColorIndex(
						int( float64(nbWordsWest) * displayRatioY), 
						int( float64(nbLinesTop) * displayRatioY), 
						redIndex)
		countryPalette[index].SetColorIndex(
						int( float64(nbWordsWest) * displayRatioY), 
						int( float64(nbLinesBottom) * displayRatioY), 
						redIndex)
		countryPalette[index].SetColorIndex(
						int( float64(nbWordsEast) * displayRatioY), 
						int( float64(nbLinesTop) * displayRatioY), 
						redIndex)
		countryPalette[index].SetColorIndex(
						int( float64(nbWordsEast) * displayRatioY), 
						int( float64(nbLinesBottom) * displayRatioY), 
						redIndex)
		fmt.Printf("{%d,\"%s\",%f,%f,%f,%f,%d},\n", 
			index, countryName[index], boundary.topLat, boundary.bottomLat, boundary.westLng, boundary.eastLng, countriesDensities[index])
	}

	// print the country
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

	for index:= 0; index<=230; index++ {
		gifImage := countryAnim[index]

		gifImage.Delay = append( gifImage.Delay, 8)
		gifImage.Image = append( gifImage.Image, countryPalette[index])
		gif.EncodeAll(countryGifFiles[index], &gifImage) // NOTE: ignoring encoding errors
		countryGifFiles[index].Close()
	}
}

