package protocols

import (
	"context"
)

type StripProcotol struct{}

func (s *StripProcotol) Name() string {
	return "strip"
}

func (s *StripProcotol) Process(ctx context.Context, data []byte) ([]byte, error) {
	// fmt.Printf("Received: %X and stripping the first 4 bytes\n", data)
	return data[4:], nil
}
