package quadtree

import (
	"testing"
)


func TestLevel(t *testing.T) {

	cases := []struct {
		in Coord
		want int
	}{
		{0, 0},
		{1<<24, 1},
		{1<<24 + 34, 1},
		{21<<24 + 1 <<23, 21},
	}
		
	for _, c := range cases {
		got := Level(c.in)
		if( got != c.want) {
			t.Errorf("Level(%b) == %d, want %d", c.in, got, c.want)
		}	
	}
}