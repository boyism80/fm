package logfile

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

type Daily struct {
	service string
	path    string
	file    *os.File
}

func New(service string) *Daily {
	return &Daily{service: service}
}

func (d *Daily) Write(p []byte) (int, error) {
	path := filepath.Join("logs", fmt.Sprintf("%s-%s.log", time.Now().Format("2006-01-02"), d.service))
	if path != d.path {
		err := d.open(path)
		if err != nil {
			return 0, err
		}
	}

	return d.file.Write(p)
}

func (d *Daily) open(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_APPEND|os.O_WRONLY, 0o644)
	switch {
	case err == nil:
		file.Write(utf8BOM)
	case os.IsExist(err):
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
	default:
		return err
	}

	if d.file != nil {
		d.file.Close()
	}
	d.path = path
	d.file = file
	return nil
}
