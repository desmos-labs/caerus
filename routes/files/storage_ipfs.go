package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/desmos-labs/caerus/utils"
)

var (
	_ Storage = &IPFSStorage{}
)

type IPFSStorage struct {
	ipfsEndpoint string
}

func NewIPFSStorage(ipfsEndpoint string) *IPFSStorage {
	return &IPFSStorage{
		ipfsEndpoint: ipfsEndpoint,
	}
}

func NewIPFStorageFromEnvVariables() *IPFSStorage {
	return NewIPFSStorage(utils.GetEnvOr(EnvFileStorageIPFSEndpoint, "https://ipfs.desmos.network"))
}

type ipfsResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func (s *IPFSStorage) UploadFile(filePath string) (string, error) {
	// Upload the file to the IPFS multipart/form-data format
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new buffer to store the form data
	body := &bytes.Buffer{}

	// Create a new multipart writer
	writer := multipart.NewWriter(body)

	// Create a new form file field
	fileField, err := writer.CreateFormFile("file", "filename.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create form file field: %s", err)
	}

	// Copy the file data to the form file field
	_, err = io.Copy(fileField, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file data: %s", err)
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %s", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v0/add", s.ipfsEndpoint), body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %s", err)
	}

	// Set the content type header to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %s", err)

	}
	defer res.Body.Close()

	// Handle the response
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)

	}

	// Parse the response
	var response ipfsResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %s", err)
	}

	// Return the file hash
	return response.Hash, nil
}

func (s *IPFSStorage) GetFile(fileName string) ([]byte, error) {
	// Build the URL used to get the file
	url := fmt.Sprintf("%s/ipfs/%s", s.ipfsEndpoint, fileName)

	// Make the request
	client := http.DefaultClient
	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}
	defer res.Body.Close()

	// Handle the response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	// Return the file content
	return resBody, nil
}
