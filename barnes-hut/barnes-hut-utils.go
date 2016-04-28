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
    "github.com/ajstarks/svgo/float"
	)


func (r * Run) RenderGif(out io.Writer) {
	const (
		size    = 600   // image canvas 
		delay   = 4    // delay between frames in 10ms units
		nframes = 0
	)
	anim := gif.GIF{LoopCount: nframes}
	rect := image.Rect(0, 0, size+1, size+1)
	img := image.NewPaletted(rect, palette)
		
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
