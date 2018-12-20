package grump

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/gyuho/goraph"
	"github.com/thomaspeugeot/tkv/quadtree"
)

func AddBodiesOfParselyPopulatedCells(
	startRow, endRow int,
	country *Country,
	parselyPopulatedCellCoords [][]bool,
	inputPopulationMatrix [][]float64,
	colLngWidth float64,
	cutoff float64,
	sampleRatio float64,
	bodies []quadtree.Body,
	popInParselyPopulatedCells, notAccountedForPop *float64) {

	Info.Printf("Construct the nodes of the graph of parsely populated cells")
	graph := goraph.NewGraph()
	Info.Printf("Nb of nodes at the start\t%10d\n", graph.GetNodeCount())

	for row := startRow; row < endRow; row++ {
		for col := 0; col < country.NCols; col++ {
			if parselyPopulatedCellCoords[row][col] == true {
				nodeID := fmt.Sprintf("%d-%d", row, col)
				n := goraph.NewNode(nodeID)
				graph.AddNode(n)
			}
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("\rrow %5d total %10d, %v MiB", row, graph.GetNodeCount(), bToMb(m.Sys))
	}
	Info.Println()

	// construct edges of the graph of parsely populated cells
	// an edge is present if the next cell to the right or below is also
	// parsely populated
	Info.Printf("Construct edges of the graph of parsely populated cells")
	for row := startRow; row < endRow-1; row++ {
		for col := 0; col < country.NCols-1; col++ {
			if parselyPopulatedCellCoords[row][col] == true {

				var from goraph.StringID
				from = goraph.StringID(fmt.Sprintf("%d-%d", row, col))

				// node to the right
				if parselyPopulatedCellCoords[row][col+1] == true {
					var to goraph.StringID = goraph.StringID(fmt.Sprintf("%d-%d", row, col+1))
					graph.AddEdge(from, to, 1.0)
					graph.AddEdge(to, from, 1.0)
				}

				// node below
				if parselyPopulatedCellCoords[row+1][col] == true {
					var to goraph.StringID = goraph.StringID(fmt.Sprintf("%d-%d", row+1, col))
					graph.AddEdge(from, to, 1.0)
					graph.AddEdge(to, from, 1.0)
				}
			}
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("\rrow %5d %v MiB", row, bToMb(m.Sys))
	}
	fmt.Printf("Graph nb of nodes\t\t%10d\n", graph.GetNodeCount())
	setOfSets := goraph.Tarjan(graph)
	fmt.Printf("Graph number of connected graph\t%10d\n", len(setOfSets))

	// parse the connected set
	// population that is not accounted for in the graph
	for setId := 0; setId < len(setOfSets); setId++ {
		popInGraph := 0.0
		for nodeRank := 0; nodeRank < len(setOfSets[setId]); nodeRank++ {
			nodeID := setOfSets[setId][nodeRank]
			var row, col int
			fmt.Sscanf(nodeID.String(), "%d-%d", &row, &col)
			*popInParselyPopulatedCells += inputPopulationMatrix[row][col]
			popInGraph += inputPopulationMatrix[row][col]

			if inputPopulationMatrix[row][col] > cutoff {
				Error.Printf("Too much pop ! %f row %d col %d", inputPopulationMatrix[row][col], row, col)
			}

			// generates a body if popInGraph above cutoff
			if popInGraph > cutoff {
				popInGraph -= cutoff

				// get lat/lng
				lat := country.Row2Lat(row)
				lng := float64(country.XllCorner) + (float64(col) * colLngWidth)
				Trace.Printf("%f %f", lat, lng)

				// compute relative coordinate of the cell
				relX, relY := country.LatLng2XY(lat, lng)

				var body quadtree.Body
				// angle := float64(i) * 2.0 * math.Pi / float64(nbBodiesInCell)
				body.X = relX + (1.0/float64(country.NCols))*0.5
				body.Y = relY + (1.0/float64(country.NRows))*0.5
				body.M = cutoff

				// sample bodies
				sample := rand.Float64() * 100.0
				if sample < sampleRatio {
					bodies = append(bodies, body)
				}
			}

			// get remainder
			if nodeRank == (len(setOfSets[setId]) - 1) {
				*notAccountedForPop += popInGraph
			}
		}
	}

	graph = nil
	runtime.GC()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
