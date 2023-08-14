package files

const (
	FileFormKey = "file"
)

type UploadFileResponse struct {
	URL string `json:"url"`
}

// NewUploadFileResponse returns a new UploadFileResponse instance
func NewUploadFileResponse(url string) *UploadFileResponse {
	return &UploadFileResponse{
		URL: url,
	}
}
