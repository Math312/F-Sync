package common

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strconv"
)

func GetFileSize(fileName string) (int32, error) {
	stat, err := os.Stat(fileName)
	if err != nil {
		log.Fatalf("File ERROR: %s", err)
		return -1, err
	}
	return int32(stat.Size()), nil
}

func FileDepart(fileName string, size int32) []string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	idx := 0
	reader := bufio.NewReader(file)
	buf := make([]byte, size)
	var result []string
	for true {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		} else {
			temp := fileName + strconv.Itoa(idx) + ".fsync"
			os.WriteFile(temp, buf[:n], 0666)
			idx++
			result = append(result, temp)
		}
	}
	return result
}

func GenMd5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
