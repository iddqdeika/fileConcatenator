package core

import (
	"fileConcatenator/definition"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const fileSizeThresholdInBytes int64 = 1073741824

var filesCreated = make(map[string]struct{})
var currentFileSize map[string]int64 = make(map[string]int64)
var currentFileIndex map[string]int = make(map[string]int)

func FileConcatenatorHandler(outputFileName string) definition.Handler {
	return func(path string) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("cant read file %v, err: %v\r\n", path, err)
			return
		}

		if currentFileSize[outputFileName] > fileSizeThresholdInBytes {
			currentFileIndex[outputFileName]++
			currentFileSize[outputFileName] = 0
		}
		fname := path[strings.LastIndex(path, string(os.PathSeparator)):]
		currentOutputFileName := outputFileName[:strings.LastIndex(outputFileName, ".")] +
			strconv.Itoa(currentFileIndex[outputFileName]) +
			outputFileName[strings.LastIndex(outputFileName, "."):]

		i, err := appendFile(currentOutputFileName, fname, data)
		if err != nil {
			fmt.Printf("cant write to file %v, err: %v\r\n", outputFileName, err)
			return
		}
		currentFileSize[outputFileName] += i
	}
}

//writes data to given filename and returns amoun t of bytes written into file
func appendFile(outputFileName string, filename string, data []byte) (int64, error) {
	if _, ok := filesCreated[outputFileName]; !ok {
		_, err := os.Create(outputFileName)
		if err != nil {
			return 0, fmt.Errorf("cant create file %v, err: %v", outputFileName, err)
		}
		filesCreated[outputFileName] = struct{}{}
	}
	f, err := os.OpenFile(outputFileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return 0, fmt.Errorf("cant open file %v, err: %v\r\n", outputFileName, err)
	}
	defer f.Close()

	_, err = f.Write([]byte("\r\n" + filename + ":\r\n{\r\n"))
	if err != nil {
		return 0, fmt.Errorf("cant write file %v, err: %v\r\n", outputFileName, err)
	}

	_, err = f.Write(data)
	if err != nil {
		return 0, fmt.Errorf("cant write file %v, err: %v\r\n", outputFileName, err)
	}

	_, err = f.Write([]byte("\r\n}\r\n"))
	if err != nil {
		return int64(len(data)), fmt.Errorf("cant write file %v, err: %v\r\n", outputFileName, err)
	}
	return int64(len(data)), nil
}
