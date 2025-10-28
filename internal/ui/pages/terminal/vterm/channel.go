package vterm

import (
	"bufio"
	"io"
)

// Channel handles communication between terminal and process
type Channel struct {
	reader *bufio.Reader
	writer io.Writer
}

// NewChannel creates a new communication channel
func NewChannel(reader io.Reader, writer io.Writer) *Channel {
	return &Channel{
		reader: bufio.NewReader(reader),
		writer: writer,
	}
}

// Read reads data from the channel
func (c *Channel) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

// Write writes data to the channel
func (c *Channel) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

// ReadByte reads a single byte
func (c *Channel) ReadByte() (byte, error) {
	return c.reader.ReadByte()
}

// WriteByte writes a single byte
func (c *Channel) WriteByte(b byte) error {
	_, err := c.writer.Write([]byte{b})
	return err
}

