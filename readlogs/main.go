package readlogs

import (
	"bufio"
	"os"
)

// Streamfile content
func StreamFile(filePath string, seekByte int64, maxLineCount int64) ([]byte, int64, int64, int64) {

	file, err := os.Open(filePath)

	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}

	file.Seek(seekByte, 0)
	if err != nil {
		panic(err)

	}

	bytesCount := 0
	block := []byte{}
	scanner := bufio.NewScanner(file)

	var currentLine int64 = 0
	for scanner.Scan() {

		l := scanner.Bytes()

		bytesCount += len(l)
		block = append(block, l...)
		currentLine++

		if currentLine >= maxLineCount {
			break
		} else {
			block = append(block, "\n"...)
			bytesCount++
		}

	}

	return block, seekByte, int64(bytesCount) + seekByte, fi.Size()

}
