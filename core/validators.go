package core

import (
	"fileProcessor/definition"
	"strings"
)

func ExtensionListValidator(extensions []string) definition.Validator {
	if extensions == nil || len(extensions) == 0 || extensions[0] == "" {
		return nil
	}
	return func(path string, isDir bool) (ok bool) {
		if isDir {
			return false
		}
		for _, ext := range extensions {
			if strings.HasSuffix(path, strings.Trim(ext, "\t ")) {
				return true
			}
		}
		return false
	}
}
