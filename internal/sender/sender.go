// Package sender provides simple object for storing ticker prices in local file in csv-like format.
package sender

import (
	"bufio"
	"fmt"
	"time"
)

type FileSender struct {
	writer *bufio.Writer
}

func NewFileSender(writer *bufio.Writer) *FileSender {
	return &FileSender{
		writer: writer,
	}
}

func (fs *FileSender) Send(timestamp time.Time, price float64) error {
	_, err := fs.writer.WriteString(fmt.Sprintf("%d,%f\n", timestamp.Unix(), price))

	return err
}
