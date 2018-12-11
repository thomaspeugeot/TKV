// Package grumpreader is the main package for the "extractor program" that parses a country file in the GRUMP format and generates a config file of bodies
//
// For each cell of the country specifc file, this program generate bodies per cells according to
// the population count in the cell
//
// The arrangement of circle in each cell is taken from a outsite source (csq something) up to 200 circles
//
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/grump"
	"github.com/thomaspeugeot/tkv/quadtree"
)

// coordinates of arrangement of circle packing in a square
type circleCoord struct {
	x, y float64
}

// storage of circle arrangement per number of circle in the square
type arrangementsStore [][]circleCoord

// coordinates of cell
type cellCoord struct {
	x, y int
}

// var targetMaxBodies = 400000
//   var targetMaxBodies = 10000000
// var targetMaxBodies = 40000
var targetMaxBodies = 100000

var maxCirclePerCell = 750

// on the PC
// go run grump-reader.go -tkvdata="C:\Users\peugeot\tkv-data"
// usage grump-reader -country=xxx where xxx is the 3 small letter ISO 3166 code for the country (for instance "fra")
func main() {

	// flag "country"
	countryPtr := flag.String("country", "fra", "iso 3166 country code")

	// flag "targetMaxBodies"
	targetMaxBodiesPtr := flag.String("targetMaxBodies", "100000", "target nb of bodies")

	// flag "sampleRatio"
	sampleRatioPtr := flag.String("sampleRatio", "100", "Ratio (in %) of output bodies, default is 100%")

	// get the directory containing tkv data through the flag "tkvdata"
	dirTKVDataPtr := flag.String("tkvdata", "/Users/thomaspeugeot/the-mapping-data/", "directory containing input tkv data")

	// use fibonacci packing, not the optimal packing
	fiboPtr := flag.Bool("fibo", true, "if true, uses fibonacci packing")

	var country grump.Country
	var sampleRatio float64

	flag.Parse()

	{
		_, errScan := fmt.Sscanf(*targetMaxBodiesPtr, "%d", &targetMaxBodies)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}

	grump.Info.Printf("country to parse %s", *countryPtr)
	country.Name = *countryPtr
	grump.Info.Printf("directory containing tkv data %s", *dirTKVDataPtr)
	dirTKVData := *dirTKVDataPtr

	// create the path to the agragate country count
	grumpFilePath := fmt.Sprintf("%s/%s_grumpv1_pcount_00_ascii_30/%sup00ag.asc", dirTKVData, *countryPtr, *countryPtr)
	grump.Info.Printf("relative path %s", filepath.Clean(grumpFilePath))
	var grumpFile *os.File
	var err error
	grumpFile, err = os.Open(filepath.Clean(grumpFilePath))
	if err != nil {
		log.Fatal(err)
	}

	// get sample ratio
	{
		_, errScan := fmt.Sscanf(*sampleRatioPtr, "%f", &sampleRatio)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}

	// parse the grump
	var word int
	scanner := bufio.NewScanner(grumpFile)
	scanner.Split(bufio.ScanWords)

	// scan 8 first lines
	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &country.NCols)
	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &country.NRows)
	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%f", &country.XllCorner)
	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%f", &country.YllCorner)

	country.Serialize()
	grump.Info.Println("country struct content is ", country)

	// scan the reamining header
	for word < 4 {
		scanner.Scan()
		word++
		fmt.Println(fmt.Sprintf("item %d : %s", word, scanner.Text()))
	}
	colLngWidth := 0.0083333333333

	// prepare the input population matrix
	inputPopulationMatrix := make([][]float64, country.NRows)

	popTotal := 0.0
	// scan the file and store result in inputPopulationMatrix
	for row := 0; row < country.NRows; row++ {
		lat := country.Row2Lat(row)
		inputPopulationMatrix[(country.NRows - row - 1)] = make([]float64, country.NCols)
		for col := 0; col < country.NCols; col++ {
			scanner.Scan()
			// lng := float64(country.XllCorner) + (float64(col)*colLngWidth)

			var nbIndividualsInCell float64
			fmt.Sscanf(scanner.Text(), "%f", &nbIndividualsInCell)

			if -2147483647 == nbIndividualsInCell {
				nbIndividualsInCell = 0
			}

			popTotal += nbIndividualsInCell

			inputPopulationMatrix[(country.NRows - row - 1)][col] = nbIndividualsInCell
		}
		fmt.Printf("\rrow %5d lat %2.3f total %f", row, lat, popTotal)
	}
	fmt.Printf("\n")
	grump.Info.Printf("reading grump file is over, closing")
	grumpFile.Close()
	fmt.Printf("pop total\t\t\t%10.0f\n", popTotal)
	cutoff := popTotal / float64(targetMaxBodies)
	fmt.Printf("pop cutoff per cell\t%10.0f\n", cutoff)

	// get the arrangement
	var arrangements arrangementsStore
	if !*fiboPtr {
		arrangements = make(arrangementsStore, maxCirclePerCell+1)
		for nbCircles := 1; nbCircles <= maxCirclePerCell; nbCircles++ {

			// fmt.Printf("\rgetting arrangement for %3d circles", nbCircles)

			arrangements[nbCircles] = make([]circleCoord, nbCircles)

			// open the reference file
			circlePackingFilePath := fmt.Sprintf("%s/csq_coords/csq%d.txt", dirTKVData, nbCircles)
			var circlePackingFile *os.File
			var errCirclePackingFile error
			circlePackingFile, errCirclePackingFile = os.Open(filepath.Clean(circlePackingFilePath))
			if errCirclePackingFile != nil {
				log.Fatal(err)
			}

			// prepare scanner
			scannerCircle := bufio.NewScanner(circlePackingFile)
			scannerCircle.Split(bufio.ScanWords)

			// one line per circle
			for circle := 0; circle < nbCircles; circle++ {

				// scan the id of the circle
				scannerCircle.Scan()

				// scan X coordinate
				scannerCircle.Scan()
				fmt.Sscanf(scannerCircle.Text(), "%f", &(arrangements[nbCircles][circle].x))
				// scan Y coordinate
				scannerCircle.Scan()
				fmt.Sscanf(scannerCircle.Text(), "%f", &(arrangements[nbCircles][circle].y))
				// fmt.Printf("getting arrangement for %d circle %3d, coord %f %f\n", nbCircles, circle, arrangements[nbCircles][circle].x, arrangements[nbCircles][circle].y)
			}
			circlePackingFile.Close()
		}
		grump.Info.Printf("reading circle packing files is over")
	} else {
		maxCirclePerCell = 10000
		arrangements = make(arrangementsStore, maxCirclePerCell+1)
		goldentRatio := 1.0 + math.Sqrt(5.0)
		for nbCircles := 1; nbCircles <= maxCirclePerCell; nbCircles++ {

			// coef is the spacing at the end and the beginning
			// of each row
			coef := math.Sqrt(float64(nbCircles)) / (math.Sqrt(float64(nbCircles)) + 1.0)

			arrangements[nbCircles] = make([]circleCoord, nbCircles)
			for circle := 0; circle < nbCircles; circle++ {
				x := (float64(circle) + 0.5) / float64(nbCircles)
				_, y := math.Modf(((float64(circle) + 0.5) * goldentRatio))

				// shrink by coef at the center
				x = 0.5 + (x-0.5)*coef
				y = 0.5 + (y-0.5)*coef

				arrangements[nbCircles][circle].x = x
				arrangements[nbCircles][circle].y = y
			}
		}
	}

	// prepare the output density file
	var bodies []quadtree.Body
	bodiesInCellMax := 0

	grump.Info.Printf("Preparing the ouput")
	cumulativePopTotal := 0.0
	nbCellsWithZeroBodies := 0
	nbCellsWithPopButWithZeroBodies := 0
	missedPopulationTotal := 0.0

	// 2D array to store wether the cell has no bodies but some pop
	parselyPopulatedCellCoords := make([][]bool, country.NRows)

	grump.Info.Printf("Parsing the pop cells and generating bodies")
	for row := 0; row < country.NRows; row++ {
		lat := country.Row2Lat(row)

		// allocate for col
		parselyPopulatedCellCoords[row] = make([]bool, country.NCols)
		for col := 0; col < country.NCols; col++ {
			lng := float64(country.XllCorner) + (float64(col) * colLngWidth)

			// compute relative coordinate of the cell
			relX, relY := country.LatLng2XY(lat, lng)

			// fetch count of the cell
			nbIndividualsInCell := inputPopulationMatrix[row][col]

			// if cell is -2147483647, then set it to 0
			if -2147483647 == nbIndividualsInCell {
				nbIndividualsInCell = 0
				inputPopulationMatrix[row][col] = 0
			}

			// how many bodies ? it is maxBodies *( nbIndividualsInCell / country.PCount)
			nbBodiesInCell := int(math.Floor(float64(targetMaxBodies) * nbIndividualsInCell / popTotal))

			massPerBody := cutoff

			if nbBodiesInCell == 0 {
				nbCellsWithZeroBodies++
			}
			if nbBodiesInCell == 0 && nbIndividualsInCell > 0 {
				nbCellsWithPopButWithZeroBodies++
				missedPopulationTotal += nbIndividualsInCell
				parselyPopulatedCellCoords[row][col] = true
				// grump.Info.Printf("Fist cell %d, %d", row, col)
			}
			if nbBodiesInCell > maxCirclePerCell {
				grump.Error.Printf("nbBodiesInCell %d superior to maxCirclePerCell %d", nbBodiesInCell, maxCirclePerCell)

				nbBodiesInCell = maxCirclePerCell
				massPerBody = float64(nbIndividualsInCell) / float64(nbBodiesInCell)
			}

			// initiate the bodies in cell
			nbBodiesInCellAfterSamplingRatio := 0
			for i := 0; i < nbBodiesInCell; i++ {
				var body quadtree.Body
				// angle := float64(i) * 2.0 * math.Pi / float64(nbBodiesInCell)
				body.X = relX + (1.0/float64(country.NCols))*(0.5+arrangements[nbBodiesInCell][i].x)
				body.Y = relY + (1.0/float64(country.NRows))*(0.5+arrangements[nbBodiesInCell][i].y)
				body.M = massPerBody

				// sample bodies
				sample := rand.Float64() * 100.0
				if sample < sampleRatio {
					bodies = append(bodies, body)
					nbBodiesInCellAfterSamplingRatio++
				}
			}
			cumulativePopTotal += nbIndividualsInCell
		}
	}

	arrangements = nil
	var popInParselyPopulatedCells, notAccountedForPop float64

	// since this is a memory hungry operation
	// the following operation is split among set of rows
	nbChunk := 20
	for chunk := 0; chunk < nbChunk; chunk++ {

		grump.Info.Printf("%d/%d to %d/%d", chunk, nbChunk, chunk+1, nbChunk)
		grump.AddBodiesOfParselyPopulatedCells(
			chunk*(country.NRows/nbChunk),
			(chunk+1)*(country.NRows/nbChunk),
			&country,
			parselyPopulatedCellCoords,
			inputPopulationMatrix,
			colLngWidth,
			cutoff,
			sampleRatio,
			bodies,
			&popInParselyPopulatedCells,
			&notAccountedForPop)
	}

	fmt.Printf("Total pop in graph cells\t%10.0f\n", popInParselyPopulatedCells)

	// var quadtree quadtree.Quadtree
	// quadtree.Init( &bodies)
	// fmt.Println(" ", )
	fmt.Printf("bodies in cell max\t\t%10d\n", bodiesInCellMax)
	fmt.Printf("cumulative pop\t\t\t%10.0f\n", cumulativePopTotal)
	fmt.Printf("nb of bodies\t\t\t%10d\n", len(bodies))
	fmt.Printf("nb of cells \t\t\t%10d\n", country.NRows*country.NCols)
	fmt.Printf("nb of cells with bodies\t\t%10d\n", country.NRows*country.NCols-nbCellsWithZeroBodies)
	fmt.Printf("nb of cells without bodies\t%10d\n", nbCellsWithZeroBodies)
	fmt.Printf("nb of cells with pop w/o bodies\t%10d\n", nbCellsWithPopButWithZeroBodies)
	fmt.Printf("graph pop of cells\t\t\t%10.0f\n", missedPopulationTotal)
	fmt.Printf("missed pop of cells w/o bodies\t%10.0f\n", notAccountedForPop)

	var run barneshut.Run
	run.Init(&bodies)
	run.OutputDir = "."
	run.SetCountry(country.Name)

	run.CaptureConfig()
}
