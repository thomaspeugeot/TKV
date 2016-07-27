// this package provides the function for retriving villages locations, borders as well as 
// translation 
package translation

import (
)

type Translation struct {


}


func (t * Translation) Init(country Country) {

	// get country coordinates
	country.Unserialize()
	country.LoadConfig( true )
	country.LoadConfig( false )

	Info.Printf("Country is %s with step %d", country.Name, country.Step)

	// rowLatWidth := 0.0083333333333
	// colLngWidth := 0.0083333333333

	// load final config
	// load initial config
}



