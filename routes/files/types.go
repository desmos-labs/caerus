package files

const (
	FileFormKey = "file"
)

// NewUploadFileResponse returns a new UploadFileResponse instance
func NewUploadFileResponse(url string) *UploadFileResponse {
	return &UploadFileResponse{
		Url: url,
	}
}
