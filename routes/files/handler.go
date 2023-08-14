package files

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // Required to properly decode GIF images
	_ "image/jpeg" // Required to properly decode JPEG images
	_ "image/png"  // Required to properly decode PNG images
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bbrks/go-blurhash"

	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	host         string
	uploadFolder string
	storage      Storage
	db           Database
}

// NewHandler returns a new Handler instance
func NewHandler(host string, filesBasePath string, storage Storage, db Database) *Handler {

	uploadsFolder := path.Join(filesBasePath, "uploads")
	utils.CreateDirIfNotExists(uploadsFolder)

	return &Handler{
		host:         host,
		uploadFolder: uploadsFolder,
		storage:      storage,
		db:           db,
	}
}

// NewHandlerFromEnvVariables builds a new Handler instance reading the configurations from the environment variables
func NewHandlerFromEnvVariables(db Database) *Handler {
	defaultBasePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return NewHandler(
		utils.GetEnvOr(types.EnvAPIsHost, "https://localhost:8080"),
		utils.GetEnvOr(EnvFileStorageBaseFolder, defaultBasePath),
		StorageFromEnvVariables(),
		db,
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
		return "", utils.WrapErr(http.StatusBadRequest, "Unsupported file type")
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
func (h *Handler) UploadFile(filePath string) (string, *UploadFileResponse, error) {
	fileName, err := h.storage.UploadFile(filePath)
	if err != nil {
		return filePath, nil, err
	}

	// Get the URL to download the file
	fileUrl := fmt.Sprintf("%[1]s/media/%[2]s", h.host, fileName)
	response := NewUploadFileResponse(fileUrl)

	// Read the image
	reader, err := os.Open(filePath)
	if err != nil {
		return filePath, response, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		// If the image is not supported, then just return the response
		if errors.Is(err, image.ErrFormat) {
			return filePath, response, nil
		}
		return filePath, response, err
	}

	// Create the image hash
	str, err := blurhash.Encode(4, 3, img)
	if err != nil {
		return filePath, response, err
	}

	// Save the image hash
	err = h.db.SaveMediaHash(fileUrl, str)
	if err != nil {
		return fileUrl, response, err
	}

	return filePath, response, nil
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
