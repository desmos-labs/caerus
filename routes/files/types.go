package files

// NewUploadFileResponse returns a new UploadFileResponse instance
func NewUploadFileResponse(fileName string) *UploadFileResponse {
	return &UploadFileResponse{
		FileName: fileName,
	}
}
