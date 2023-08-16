package files

import (
	"io"
	"os"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/utils"
)

var (
	_ FilesServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func NewServerFromEnvVariables(db Database) *Server {
	return &Server{
		handler: NewHandlerFromEnvVariables(db),
	}
}

// UploadFile implements FilesServiceServer
func (s *Server) UploadFile(stream FilesService_UploadFileServer) error {
	_, err := authentication.GetAuthenticatedUserData(stream.Context())
	if err != nil {
		return err
	}

	var fileName string
	var fileData []byte
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// If we read the EOF, it means the upload has completed,
			// and we are ready to process the file properly

			// First of all, store the file
			filePath, err := s.handler.SaveFile(fileName, fileData)
			if err != nil {
				return utils.UnwrapError(stream.Context(), err)
			}

			// Now, handle the request
			tempFilePath, res, err := s.handler.UploadFile(filePath)
			if err != nil {
				return utils.UnwrapError(stream.Context(), err)
			}
			os.Remove(tempFilePath)

			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		fileName = chunk.FileName
		fileData = append(fileData, chunk.Data...)
	}
}

// GetFile implements FilesServiceServer
func (s *Server) GetFile(request *GetFileRequest, stream FilesService_GetFileServer) error {
	// Download the file contents into a temporary file
	tempFilePath, err := s.handler.GetFile(request.FileName)
	if err != nil {
		return utils.UnwrapError(stream.Context(), err)
	}
	defer os.Remove(tempFilePath)

	// Open the file
	file, err := os.Open(tempFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Send the file using a stream and a buffer as support
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := &FileChunk{
			Data: buffer[:bytesRead],
		}
		if err := stream.Send(chunk); err != nil {
			return err
		}
	}

	return nil
}
