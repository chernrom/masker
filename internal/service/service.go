package service

import "fmt"

type Producer interface {
	Produce() ([]string, error)
}
type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service {
	if prod == nil {
		panic("nil Producer")
	}
	if pres == nil {
		panic("nil Presenter")
	}
	return &Service{prod: prod, pres: pres}
}

func (s *Service) Run() error {
	in, err := s.prod.Produce()
	if err != nil {
		return fmt.Errorf("produce: %w", err)
	}

	out := make([]string, len(in))
	for i, line := range in {
		out[i] = s.maskHttpLinks(line)
	}

	if err := s.pres.Present(out); err != nil {
		return fmt.Errorf("present: %w", err)
	}
	return nil
}

func (s *Service) maskHttpLinks(str string) string {
	b := []byte(str)
	pattern := []byte("http://")

	for i := 0; i+len(pattern) <= len(b); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if b[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			for k := i + len(pattern); k < len(b) && b[k] != ' ' && b[k] != '\n'; k++ {
				b[k] = '*'
			}
		}
	}
	return string(b)
}
