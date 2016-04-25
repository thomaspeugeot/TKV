//
// grump reader parses a country specific file
//
// usage grump-reader -country=xxx where xxx isthe 3 small letter ISO 3166 code for the country (for instance "fra")
// 
package main

import "flag"
import "fmt"
import "os"
import "log"
import "bufio"
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
	lineLatWidth := 0.0083333333333
	// columnLngWidth := 0.0083333333333

	// prepare the count matrix
	countMatrix := make([]float64, country.NRows * country.NCols)

	// prepare the density file
	maxBodies := 200000
	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, maxBodies)

	popTotal := 0.0
	// scan the file
	for row :=0; row < country.NRows; row++ {
		lat := float64( country.YllCorner) + (float64( country.NRows - row)*lineLatWidth)
		for col :=0; col < country.NCols ; col++ {
			scanner.Scan()
			// lng := float64(country.XllCorner) + (float64(col)*columnLngWidth)

			var count float64
			fmt.Sscanf( scanner.Text(), "%f", &count)
			popTotal += count

			countMatrix[ row*country.NCols + col ] = count
		}
		fmt.Printf("\rrow %5d lat %2.3f total %f", row, lat, popTotal)
	}
	fmt.Println("")
}