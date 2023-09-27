package files

import (
	_ "image/gif"  // Required to properly decode GIF images
	_ "image/jpeg" // Required to properly decode JPEG images
	_ "image/png"  // Required to properly decode PNG images
	"net/http"
	"os"
	"path"
	"strings"

	"google.golang.org/grpc/codes"

	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	uploadFolder string
	storage      Storage
}

// NewHandler returns a new Handler instance
func NewHandler(filesBasePath string, storage Storage) *Handler {
	uploadsFolder := path.Join(filesBasePath, "uploads")
	utils.CreateDirIfNotExists(uploadsFolder)

	return &Handler{
		uploadFolder: uploadsFolder,
		storage:      storage,
	}
}

// NewHandlerFromEnvVariables builds a new Handler instance reading the configurations from the environment variables
func NewHandlerFromEnvVariables() *Handler {
	defaultBasePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return NewHandler(
		utils.GetEnvOr(EnvFileStorageBaseFolder, defaultBasePath),
		StorageFromEnvVariables(),
	)
}

// --------------------------------------------------------------------------------------------------------------------

// createUploadsFolder creates the folder where to store the uploaded files, if not existing
func (h *Handler) createUploadsFolder() error {
	if _, err := os.Stat(h.uploadFolder); os.IsNotExist(err) {
		err = os.Mkdir(h.uploadFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveFile reads the given request and stores the associated file bytes
// into a new file within the upload folder.
func (h *Handler) SaveFile(fileName string, data []byte) (string, error) {

	// Create the upload folder if it does not exist
	err := h.createUploadsFolder()
	if err != nil {
		return "", err
	}

	mimeType := strings.Split(http.DetectContentType(data), "/")[0]
	if mimeType != "image" && mimeType != "video" {
		return "", utils.WrapErr(codes.InvalidArgument, "Unsupported file type")
	}

	// Copy the uploads file bytes to the out file
	filePath := path.Join(h.uploadFolder, fileName)
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// UploadFile writes the contents of the file located at the given path in a temporary file, and uploads them remotely.
// After the uploads is completed, the path to the temporary file is returned along with the upload response.
// Note: The caller should make sure the temporary file is deleted.
func (h *Handler) UploadFile(filePath string) (*UploadFileResponse, error) {
	fileName, err := h.storage.UploadFile(filePath)
	if err != nil {
		return nil, err
	}

	return NewUploadFileResponse(fileName), nil
}

// GetFile gets the contents of the file having the given name, and writes them on a temporary
// file located at the returned path.
// Note: The caller should make sure the temporary file is deleted.
func (h *Handler) GetFile(fileName string) (string, error) {
	bz, err := h.storage.GetFile(fileName)
	if err != nil {
		return "", nil
	}

	filePath := path.Join(h.uploadFolder, fileName)
	return filePath, os.WriteFile(filePath, bz, 0600)
}
