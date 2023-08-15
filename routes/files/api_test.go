package files_test

import (
	"context"
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/testutils"

	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/routes/files"
	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

func TestFilesAPIsTestSuite(t *testing.T) {
	suite.Run(t, new(FilesAPITestSuite))
}

type FilesAPITestSuite struct {
	suite.Suite

	db      *database.Database
	storage files.Storage
	handler *files.Handler

	server *grpc.Server
	client files.FilesServiceClient

	tempDir string
}

func (suite *FilesAPITestSuite) SetupSuite() {
	// Setup the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db

	// Create the directories for this suite
	tempDir, err := os.MkdirTemp(suite.T().TempDir(), "*")
	suite.Require().NoError(err)
	suite.tempDir = tempDir

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create the handler
	suite.storage = files.NewIPFSStorage("https://ipfs.desmos.network")
	suite.handler = files.NewHandler(suite.tempDir, suite.storage, suite.db)

	// Create the server
	suite.server = testutils.CreateServer(suite.db)

	// Register the service
	service := files.NewServer(suite.handler)
	files.RegisterFilesServiceServer(suite.server, service)

	// Setup the client
	conn, err := testutils.StartServerAndConnect(suite.server)
	suite.Require().NoError(err)
	suite.client = files.NewFilesServiceClient(conn)
}

func (suite *FilesAPITestSuite) TearDownSuite() {
	suite.server.Stop()
}

func (suite *FilesAPITestSuite) SetupTest() {
	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)

	// Clean all the temporary files folders
	err = os.RemoveAll(path.Join(suite.tempDir, "uploads"))
	suite.Require().NoError(err)
	utils.CreateDirIfNotExists(path.Join(suite.tempDir, "uploads"))
	err = os.RemoveAll(path.Join(suite.tempDir, "storage"))
	suite.Require().NoError(err)
	utils.CreateDirIfNotExists(path.Join(suite.tempDir, "storage"))
}

func (suite *FilesAPITestSuite) uploadFile(ctx context.Context, filePath string) (*files.UploadFileResponse, error) {
	stream, err := suite.client.UploadFile(ctx)

	file, err := os.Open(filePath)
	suite.Require().NoError(err)
	defer file.Close()

	fileInfo, err := file.Stat()
	suite.Require().NoError(err)

	buffer := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		chunk := &files.FileChunk{
			FileName: fileInfo.Name(),
			Data:     buffer[:bytesRead],
		}

		err = stream.Send(chunk)
		if err != nil {
			return nil, err
		}
	}

	return stream.CloseAndRecv()
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *FilesAPITestSuite) TestUploadMedia() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		shouldErr    bool
		check        func(res *files.UploadFileResponse)
	}{
		{
			name: "invalid session returns error",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Date(2020, 1, 1, 12, 00, 00, 000, time.UTC),
					time.Date(2020, 1, 2, 12, 00, 00, 000, time.UTC),
				))
				suite.Require().NoError(err)

				res, err := testutils.GetTestImage()
				suite.Require().NoError(err)

				outFile, err := os.Create(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer outFile.Close()

				_, err = outFile.Write(res)
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: true,
		},
		{
			name: "valid request works properly",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(3*time.Hour),
				))
				suite.Require().NoError(err)

				res, err := testutils.GetTestImage()
				suite.Require().NoError(err)

				outFile, err := os.Create(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer outFile.Close()

				_, err = outFile.Write(res)
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: false,
			check: func(res *files.UploadFileResponse) {
				suite.Require().NotEmpty(res.FileName)

				// Make sure the image hash has been saved properly
				var hash string
				err := suite.db.SQL.QueryRow(`SELECT hash FROM files_hashes WHERE file_name = $1`, res.FileName).Scan(&hash)
				suite.Require().NoError(err)
				suite.Require().Equal("L-J[3W*E#u;2%Lb:sE$OWBe@R%NH", hash)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			ctx := context.Background()
			if tc.setupContext != nil {
				ctx = tc.setupContext(ctx)
			}

			// Perform the request
			res, err := suite.uploadFile(ctx, path.Join(suite.tempDir, "temp_file.jpeg"))

			// Check the response
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(res)
			}
		})
	}
}

func (suite *FilesAPITestSuite) TestDownloadMedia() {
	testCases := []struct {
		name         string
		setup        func() (fileName string)
		setupContext func() context.Context
		shouldErr    bool
		check        func(data []byte)
	}{
		{
			name: "valid request works properly",
			setup: func() string {
				// Store the session for the upload
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(3*time.Hour),
				))
				suite.Require().NoError(err)

				// Get the test image
				imageBz, err := testutils.GetTestImage()
				suite.Require().NoError(err)

				outFile, err := os.Create(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer outFile.Close()

				_, err = outFile.Write(imageBz)
				suite.Require().NoError(err)

				// Upload the test image
				ctx := testutils.SetupContextWithAuthorization(context.Background(), "token")
				res, err := suite.uploadFile(ctx, path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)

				return res.FileName
			},
			shouldErr: false,
			check: func(data []byte) {
				expectedBz, err := testutils.GetTestImage()
				suite.Require().NoError(err)
				suite.Require().Equal(expectedBz, data)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			var fileName string
			if tc.setup != nil {
				fileName = tc.setup()
			}

			// Perform the request
			request := &files.GetFileRequest{FileName: fileName}
			stream, err := suite.client.GetFile(context.Background(), request)
			suite.Require().NoError(err)

			var data []byte

			for {
				chunk, err := stream.Recv()
				if err == io.EOF {
					break
				}
				suite.Require().NoError(err)

				data = append(data, chunk.Data...)
			}

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(data)
			}
		})
	}
}
