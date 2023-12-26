package utils

import (
	"bufio"
	"fmt"
	"os"
)

func Save2File(filePath string, contents []string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("failed to open file", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, content := range contents {
		_, err = write.WriteString(content + "\n")
	}
	return write.Flush()
}
