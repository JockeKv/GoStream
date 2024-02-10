package stream

import (
	"io"
	"os"
)

// A interface to represent a File.
// Needs to satisfy the ReadSeekCloser interface and Name() string
type File interface {
	io.ReadSeekCloser
	Name() string
}

// A implementation of File with an underlying buffer instead of file descriptor.
type FileBuffer struct {
	Reader   io.ReadSeeker
	Filename string

	CloseFunc func() error
}

// Pass the request to the ReaderSeeker
func (f *FileBuffer) Read(b []byte) (int, error) {
	return f.Reader.Read(b)
}

// Pass the request to the ReaderSeeker
func (f *FileBuffer) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}

// Name of the file
func (f *FileBuffer) Name() string {
	return f.Filename
}

// The function to be run when the file is no longer needed. Eg. return memory to a pool.
func (f *FileBuffer) Close() error {
	return f.CloseFunc()
}

// Return a Stream of *os.File's.
// If a file fails to open, the file is skipped silently.
func Files(filenames ...string) Stream[File] {
	c := make(chan File, len(filenames))
	go func() {
		defer close(c)
		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				continue
			}
			c <- file
		}
	}()
	return c
}

// Return a Stream of *os.File's from a directory.
// If a file fails to open or is a sub-directory, the file is skipped silently.
func Dir(directory string) Stream[File] {
	filenames, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}
	c := make(chan File, len(filenames))
	go func() {
		defer close(c)
		for _, filename := range filenames {
			if filename.IsDir() {
				continue
			}
			file, err := os.Open(directory + "/" + filename.Name())
			if err != nil {
				continue
			}
			c <- file
		}
	}()
	return c
}
