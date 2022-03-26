package app

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/vvkosty/go_sprint_1/internal/app/helpers"
)

const delimiter = "|"

type FileDatabase struct {
	urls    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func NewFileStorage(path string) *FileDatabase {
	newFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	fd := FileDatabase{
		urls:    newFile,
		writer:  bufio.NewWriter(newFile),
		scanner: bufio.NewScanner(newFile),
	}
	return &fd
}

func (fd *FileDatabase) Find(id string) (string, error) {
	result := ""
	fd.urls.Seek(0, io.SeekStart)
	for fd.scanner.Scan() {
		url := strings.Split(fd.scanner.Text(), delimiter)
		if id == url[0] {
			result = url[1]
			break
		}
	}

	if err := fd.scanner.Err(); err != nil {
		return "", err
	}

	return result, nil
}

func (fd *FileDatabase) Save(url string, userID string) (string, error) {
	checksum := helpers.GenerateChecksum(url)
	if _, err := fd.writer.WriteString(checksum + delimiter + url + delimiter + userID + "\n"); err != nil {
		return "", err
	}

	return checksum, fd.writer.Flush()
}

func (fd *FileDatabase) List(userID string) map[string]string {
	result := make(map[string]string)
	fd.urls.Seek(0, io.SeekStart)
	for fd.scanner.Scan() {
		url := strings.Split(fd.scanner.Text(), delimiter)
		if userID == url[2] {
			result[url[0]] = url[1]
		}
	}

	return result
}

func (fd *FileDatabase) Close() error {
	return fd.urls.Close()
}
