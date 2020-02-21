package core

import (
	"fileProcessor/definition"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var filesCreated = make(map[string]struct{})

func FileConcatenatorHandler(outputFileName string) definition.Handler {
	return func(path string) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("cant read file %v, err: %v\r\n", path, err)
			return
		}
		fname := path[strings.LastIndex(path, string(os.PathSeparator)):]
		err = appendFile(outputFileName, fname, data)
		if err != nil {
			fmt.Printf("cant write to file %v, err: %v\r\n", outputFileName, err)
			return
		}
	}
}

func appendFile(outputFileName string, filename string, data []byte) error {
	if _, ok := filesCreated[outputFileName]; !ok {
		_, err := os.Create(outputFileName)
		if err != nil {
			return fmt.Errorf("cant create file %v, err: %v", outputFileName, err)
		}
		filesCreated[outputFileName] = struct{}{}
	}
	f, err := os.OpenFile(outputFileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("cant open file %v, err: %v\r\n", outputFileName, err)
	}
	defer f.Close()

	_, err = f.Write([]byte("\r\n" + filename + ":\r\n{\r\n"))
	if err != nil {
		return fmt.Errorf("cant write file %v, err: %v\r\n", outputFileName, err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("cant write file %v, err: %v\r\n", outputFileName, err)
	}

	_, err = f.Write([]byte("\r\n}\r\n"))
	if err != nil {
		return fmt.Errorf("cant write file %v, err: %v\r\n", outputFileName, err)
	}
	return nil
}
