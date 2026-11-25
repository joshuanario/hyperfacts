package serializables

import (
	"encoding/json"
	"os"
	"path"
)

type Underscore struct {
	Entry     string
	Documents []string
	Note      string
}

func GetEntryFileName(id string) (string, error) {
	underscore := Underscore{}
	b, err := os.ReadFile(path.Join(WorkingDirectoryPath, "folders", id, "_"))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &underscore)
	return underscore.Entry, err
}
