package mkpub

import (
	"os"
)

func WriteCollectionToHTML(collection *Collection, outputDir string) error {
	_ = collection
	_ = outputDir

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	return nil
}
