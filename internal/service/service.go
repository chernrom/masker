package service

import (
	"fmt"
	"sync"
)

const maxGoroutines = 10

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

	type job struct {
		idx  int
		line string
	}

	type result struct {
		idx  int
		line string
	}

	jobs := make(chan job)
	results := make(chan result)

	sem := make(chan struct{}, maxGoroutines) //ставим лимит, вынес отдельно в константу
	var wg sync.WaitGroup

	//worker
	wg.Add(maxGoroutines)
	go func() {
		for j := range jobs {
			sem <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() {
					<-sem
				}()

				masked := s.maskHttpLinks(j.line)
				results <- result{
					idx:  j.idx,
					line: masked,
				}
			}()
		}

	}()

	//producer
	go func() {
		for i, line := range in {
			jobs <- job{
				idx:  i,
				line: line,
			}
		}
		close(jobs)
	}()

	//watcher
	go func() {
		wg.Wait()
		close(results)
	}()

	//collector
	for r := range results {
		out[r.idx] = r.line
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
