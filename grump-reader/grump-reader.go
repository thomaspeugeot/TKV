//
// grump reader parses a country specific file
//
// usage grump-reader -country=xxx where xxx isthe 3 small letter ISO 3166 code for the country (for instance "fra")
// 
package main

import "flag"
import "math"
// import "math/rand"
import "fmt"
import "os"
import "log"
import "bufio"
import "github.com/thomaspeugeot/tkv/barnes-hut"
import "github.com/thomaspeugeot/tkv/quadtree"

// store country code
type country struct {
	Name string
	NCols, NRows, XllCorner, YllCorner int
}	



func main() {

	// flag "country"
	countryPtr := flag.String("country","fra","iso 3166 country code")

	var country country

	flag.Parse()
	fmt.Println( "country to parse", *countryPtr)
	country.Name = *countryPtr

	// create the path to the agragate country count
	grumpFilePath := fmt.Sprintf( "/Users/thomaspeugeot/the-mapping-data/%s_grumpv1_pcount_00_ascii_30/%sup00ag.asc", *countryPtr, *countryPtr )
	var grumpFile *os.File
	var err error
	grumpFile, err = os.Open( grumpFilePath)
	if err != nil {
		log.Fatal(err)
	}	

	// parse the grump
	var word int
	scanner := bufio.NewScanner( grumpFile)
	scanner.Split(bufio.ScanWords)

	// scan 8 first lines
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.NCols)
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.NRows)
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.XllCorner)
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.YllCorner)

	fmt.Println( country )

	// scan the reamining header
	for word < 4 {
		scanner.Scan()
		word++		
		fmt.Println( fmt.Sprintf("item %d : %s", word, scanner.Text()))
	}
	rowLatWidth := 0.0083333333333
	colLngWidth := 0.0083333333333

	// prepare the count matrix
	countMatrix := make([]float64, country.NRows * country.NCols)

	popTotal := 0.0
	// scan the file and store result in countMatrix
	for row :=0; row < country.NRows; row++ {
		lat := float64( country.YllCorner) + (float64( country.NRows - row)*rowLatWidth)
		for col :=0; col < country.NCols ; col++ {
			scanner.Scan()
			// lng := float64(country.XllCorner) + (float64(col)*colLngWidth)

			var count float64
			fmt.Sscanf( scanner.Text(), "%f", &count)
			popTotal += count

			countMatrix[ (country.NRows-row-1)*country.NCols + col ] = count
		}
		fmt.Printf("\rrow %5d lat %2.3f total %f", row, lat, popTotal)
	}
	fmt.Println("")

	// prepare the output density file
	targetMaxBodies := 200000
	var bodies []quadtree.Body

	cumulativePopTotal := 0.0
	for row :=0; row < country.NRows; row++ {
		lat := float64( country.YllCorner) + (float64( country.NRows - row)*rowLatWidth)
		for col :=0; col < country.NCols ; col++ {
			lng := float64(country.XllCorner) + (float64(col)*colLngWidth)

			// compute relative coordinate of the cell
			relX := (lng - float64(country.XllCorner)) / (float64(country.NCols) * colLngWidth)
			relY := (lat - float64(country.YllCorner)) / (float64(country.NRows) * rowLatWidth)

			// fetch count of the cell
			count := countMatrix[ row*country.NCols + col ]

			// how many bodies ? it is maxBodies *( count / country.PCount) 
			bodiesInCell := int( math.Floor( float64( targetMaxBodies) * (count/popTotal)))
			// newBodies := make( []quadtree.Body, bodiesInCell)

			// initiate the bodies
			for i :=0; i<bodiesInCell; i++ {
				var body quadtree.Body
				angle := float64(i) * 2.0 * math.Pi / float64(bodiesInCell)
				body.X = relX + (1.0/float64(country.NCols))*(0.5 + 0.3*math.Cos(angle))
				body.Y = relY + (1.0/float64(country.NRows))*(0.5 + 0.3*math.Sin(angle))
				body.M = count/float64(bodiesInCell)
				bodies = append( bodies,  body)
			}
			cumulativePopTotal += count

		}
	}

	// var quadtree quadtree.Quadtree
	// quadtree.Init( &bodies)

	var run barnes_hut.Run
	run.Init( & bodies)

	run.CaptureConfigCountry( country.Name)
}