package protocols

import "context"

type StripProcotol struct{}

func (s *StripProcotol) Name() string {
	return "strip"
}

func (s *StripProcotol) Process(ctx context.Context, data []byte) ([]byte, error) {
	return data[4:], nil
}
