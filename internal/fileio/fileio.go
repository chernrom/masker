package fileio

import (
	"bufio"
	"fmt"
	"os"
)

type FileProducer struct {
	path string
}

func NewFileProducer(path string) *FileProducer {
	return &FileProducer{path: path}
}

func (p *FileProducer) Produce() ([]string, error) {
	f, err := os.Open(p.path)
	if err != nil {
		return nil, fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan input: %w", err)
	}
	return lines, nil
}

type FilePresenter struct {
	path string
}

func NewFilePresenter(path string) *FilePresenter {
	if path == "" {
		path = "out.txt" // дефолт, если не задан
	}
	return &FilePresenter{path: path}
}

func (p *FilePresenter) Present(lines []string) error {
	f, err := os.Create(p.path) // создаёт или усекает существующий файл
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, ln := range lines {
		if _, err := w.WriteString(ln + "\n"); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}
