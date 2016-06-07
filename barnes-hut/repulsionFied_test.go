package barnes_hut

import (
	"testing"
)

func TestRepulsionFieldInit(t *testing.T) {

	f := NewRepulsionField( 100)
	f.ComputeField()
}
