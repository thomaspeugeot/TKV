//
// contains utilies functions for the run
//
package barnes_hut

import (
	"image"
	"image/gif"
	"io"
	"log"
	"fmt"
	"math"
	"time"
	"bytes"
	"encoding/base64"
	"image/color"
    "github.com/ajstarks/svgo/float"
	)


// var palette = []color.Color{color.White, color.Black}
var palette = []color.Color{color.White, color.Black, color.RGBA{255,0,0,255}, color.RGBA{0,255,0,255}, }
const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
	redIndex = 2 // next color in palette
	blueIndex = 3 // next color in palette
)

const ( NbPaletteGrays = 100 
		Padding = 4
		)
// init the palette with the gray depth
func init() {
	for gray := 0; gray < NbPaletteGrays; gray++ {
		palette = append( palette, color.RGBA{uint8(250-gray), uint8(250-gray), uint8(250-gray), 255})
	}
}


func (r * Run) RenderGif(out io.Writer) {

	renderingMutex.Lock()
	t0 := time.Now()

	Trace.Printf("RenderGif begin with r.gridFieldNb %d", r.gridFieldNb)
	
	const (
		size    = 600   // image canvas 
		delay   = 4    // delay between frames in 10ms units
		nframes = 0
	)
	anim := gif.GIF{LoopCount: nframes}
	rect := image.Rect(0, 0, size+1, size+1)
	img := image.NewPaletted(rect, palette)
	
	// compute the field
	if r.fieldRendering {
		f := NewRepulsionField( r.xMin, r.yMin, 
							r.xMax, r.yMax, 
							r.gridFieldNb,
							&(r.q),
							r.minInterBodyDistance/2.0) // quadtree
		f.ComputeField()

		// parse the image 
		for i:=0; i<size+1;i++ {
			for j:=0;j<size+1;j++ {
				fx := int(math.Floor( (float64(i)/float64(size+1)) * float64( r.gridFieldNb ) ))
				fy := int(math.Floor( (float64(j)/float64(size+1)) * float64( r.gridFieldNb ) ))
				
				field := f.values[fx][fy]
				indexPalette := uint8( Padding + math.Floor( (field/f.maxValue)*(NbPaletteGrays-1)))
				if ( i % (size/r.gridFieldNb) ==0) {
					if ( j % (size/r.gridFieldNb) ==0) {
				 		Trace.Printf("RenderGif pixel %3d %3d, grid coord %3d %3d f %e, index %d", i,j, fx, fy, field, indexPalette)
					}
				}
				
				img.SetColorIndex( 
					i, 
					j,
					indexPalette)
			}
		}
	} 


	for idx, _ := range (*r.bodies) {
	
		body := (*r.bodies)[idx]
		bodyOrig := (*r.bodiesOrig)[idx]
	
		if false { fmt.Printf("Encoding body %d %f %f\n", idx, body.X, body.Y) }
	
		// take into account rendering window
		// in gif, A Point is an (x, y) co-ordinate on the integer grid, with axes increasing right and down.
		var imX, imY float64
		if( r.renderChoice == RUNNING_CONFIGURATION) {
			imX = (body.X - r.xMin)/(r.xMax-r.xMin)
			imY = (r.yMax - body.Y)/(r.yMax-r.yMin) // coordinates in y are down
		} else { 
			// we display the original
			imX = (bodyOrig.X - r.xMin)/(r.xMax-r.xMin)
			imY = (r.yMax - bodyOrig.Y)/(r.yMax-r.yMin)	// coordinates in y are down
		}
		// if( (body.X > r.xMin) && (body.X < r.xMax) && (body.Y > r.yMin) && (body.Y < r.yMax) ) {
		if( (imX > 0.0) && (imX < 1.0) && (imY > 0.0) && (imY < 1.0) ) {

			// check wether body is on a border
			isOnBorder := false
			coordX := body.X * float64(nbVillagePerAxe)
			distanceToBorderX := coordX - math.Floor( coordX)
			if( distanceToBorderX < ratioOfBorderVillages /2.0) { isOnBorder = true }
			if( distanceToBorderX > 1.0 -  ratioOfBorderVillages /2.0) { isOnBorder = true }

			coordY := body.Y * float64(nbVillagePerAxe)
			distanceToBorderY := coordY - math.Floor( coordY)
			if( distanceToBorderY < ratioOfBorderVillages / 2.0) { isOnBorder = true }
			if( distanceToBorderY > 1.0 -  ratioOfBorderVillages /2.0) { isOnBorder = true }

			// compute village coordinate (from 0 to nbVillagePerAxe-1)
			x := int( math.Floor(float64( nbVillagePerAxe) * body.X))
			y := int( math.Floor(float64( nbVillagePerAxe) * body.Y))

			// we want to alternate red and blue
			var borderIndex uint8
			borderIndex = redIndex
			if( (x+y)%2 ==0) { borderIndex = blueIndex }

			if( isOnBorder && r.renderState == WITH_BORDERS) {
				img.SetColorIndex(
					int(imX*size+0.5), 
					int(imY*size+0.5),
					borderIndex)
				img.SetColorIndex(
					int(imX*size+0.5)+1, 
					int(imY*size+0.5),
					borderIndex)
				img.SetColorIndex(
					int(imX*size+0.5)-1, 
					int(imY*size+0.5),
					borderIndex)
				img.SetColorIndex(
					int(imX*size+0.5), 
					int(imY*size+0.5)+1,
					borderIndex)
				img.SetColorIndex(
					int(imX*size+0.5), 
					int(imY*size+0.5)-1,
					borderIndex)
			} else {
				img.SetColorIndex(
					int(imX*size+0.5), 
					int(imY*size+0.5),
					blackIndex)				
			}
		}
	}
	anim.Delay = append(anim.Delay, delay)
	anim.Image = append(anim.Image, img)
	var b bytes.Buffer
	gif.EncodeAll(&b, &anim)
	encodedB64 := base64.StdEncoding.EncodeToString([]byte(b.Bytes()))
	out.Write( []byte(encodedB64))

	
	t1 := time.Now()
	StepDuration = float64((t1.Sub(t0)).Nanoseconds())
	
	Trace.Printf("RenderGif %d dur %e", r.gridFieldNb, StepDuration/1000000000)
	renderingMutex.Unlock()

}


func (r * Run) RenderSVG(out io.Writer) {
	const (
		size    = 600   // image canvas 
	)
	s := svg.New(out)
	s.Start(size, size)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	
	for idx, _ := range (*r.bodies) {
	
		body := (*r.bodies)[idx]
		bodyOrig := (*r.bodiesOrig)[idx]
	
		if false { fmt.Printf("Encoding body %d %f %f\n", idx, body.X, body.Y) }
	
		// take into account rendering window
		var imX, imY float64
		if( r.renderChoice == RUNNING_CONFIGURATION) {
			imX = (body.X - r.xMin)/(r.xMax-r.xMin)
			imY = (body.Y - r.yMin)/(r.yMax-r.yMin)
		} else { 
			// we display the original
			imX = (bodyOrig.X - r.xMin)/(r.xMax-r.xMin)
			imY = (bodyOrig.Y - r.yMin)/(r.yMax-r.yMin)				
		}
		// if( (body.X > r.xMin) && (body.X < r.xMax) && (body.Y > r.yMin) && (body.Y < r.yMax) ) {
		if( (imX > 0.0) && (imX < 1.0) && (imY > 0.0) && (imY < 1.0) ) {


			// check wether body is on a border
			isOnBorder := false
			coordX := body.X * float64(nbVillagePerAxe)
			distanceToBorderX := coordX - math.Floor( coordX)
			if( distanceToBorderX < ratioOfBorderVillages /2.0) { isOnBorder = true }
			if( distanceToBorderX > 1.0 -  ratioOfBorderVillages /2.0) { isOnBorder = true }

			coordY := body.Y * float64(nbVillagePerAxe)
			distanceToBorderY := coordY - math.Floor( coordY)
			if( distanceToBorderY < ratioOfBorderVillages / 2.0) { isOnBorder = true }
			if( distanceToBorderY > 1.0 -  ratioOfBorderVillages /2.0) { isOnBorder = true }

			if( isOnBorder && r.renderState == WITH_BORDERS) {
				s.Circle(imX*size, imY*size, 0.1, "fill:none;stroke:red")
			} else {
				s.Circle(imX*size, imY*size, 0.1, "fill:none;stroke:black")
			}
		}
	}
	s.End()
	log.Output( 1, fmt.Sprintf( "end of render SVG"))
}

// output position of bodies of the Run into a GIF representation
func (r * Run) OutputGif(out io.Writer, nbStep int) {

	for r.step < nbStep  {

		// if state is STOPPED, pause
		for r.state == STOPPED {
			time.Sleep(100 * time.Millisecond)
		}
		r.q.ComputeQuadtreeGini()

		// append the new gini elements
		// create the array
		giniArray := make( []float64, 10)
		copy( giniArray, r.q.BodyCountGini[8][:])
		r.giniOverTime = append( r.giniOverTime, giniArray)

		r.OneStep()
	}
}
