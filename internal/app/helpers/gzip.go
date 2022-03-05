package helpers

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

type GzipHelper struct {
	b bytes.Buffer
}

func (gz *GzipHelper) Compress(data []byte) ([]byte, error) {
	gz.b = bytes.Buffer{}

	w, err := gzip.NewWriterLevel(&gz.b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	// запись данных
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	// переменная b содержит сжатые данные
	return gz.b.Bytes(), nil
}

func (gz *GzipHelper) Decompress(data []byte) ([]byte, error) {
	// переменная r будет читать входящие данные и распаковывать их
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	defer r.Close()

	var b bytes.Buffer
	// в переменную b записываются распакованные данные
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
