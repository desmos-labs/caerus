package files

import (
	"fmt"
	"strings"

	"github.com/desmos-labs/caerus/utils"
)

type Storage interface {
	// UploadFile uploads a file to the storage and returns the uploaded file name
	// that can be used to retrieve the file later on using the GetFile method.
	UploadFile(filePath string) (string, error)

	// GetFile returns the file associated with the given name
	GetFile(fileName string) ([]byte, error)
}

func StorageFromEnvVariables() Storage {
	storageType := utils.GetEnvOr(EnvFileStorageType, StorageTypeIPFS)

	switch {
	case strings.EqualFold(storageType, StorageTypeIPFS):
		return NewIPFStorageFromEnvVariables()

	default:
		panic(fmt.Errorf("invalid storage type %s", storageType))
	}
}
