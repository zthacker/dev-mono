package protocols

import (
	"fmt"
	"testing"
)

func TestStrip(t *testing.T) {
	s := &StripProcotol{}
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res, err := s.Process(nil, data)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%d", len(res))
	fmt.Printf("%b", res)
}
