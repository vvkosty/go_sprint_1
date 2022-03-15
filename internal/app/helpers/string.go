package helpers

import (
	"hash/crc32"
	"strconv"
)

func GenerateChecksum(url string) string {
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(url))))
	return checksum
}
