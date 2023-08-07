package files_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/caerus/server/testutils"

	"github.com/desmos-labs/caerus/server/database"
	"github.com/desmos-labs/caerus/server/routes/files"
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
	r       *gin.Engine

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
	suite.handler = files.NewHandler("http://localhost", suite.tempDir, suite.storage, suite.db)

	// Create the router
	router, err := testutils.CreateRouter()
	suite.Require().NoError(err)
	suite.r = router

	// Register the transactions APIs
	files.Register(suite.r, suite.handler)
}

func (suite *FilesAPITestSuite) SetupTest() {
	err := os.RemoveAll(path.Join(suite.tempDir, "uploads"))
	suite.Require().NoError(err)
	utils.CreateDirIfNotExists(path.Join(suite.tempDir, "uploads"))
	err = os.RemoveAll(path.Join(suite.tempDir, "storage"))
	suite.Require().NoError(err)
	utils.CreateDirIfNotExists(path.Join(suite.tempDir, "storage"))
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *FilesAPITestSuite) TestUploadMedia() {
	testCases := []struct {
		name         string
		setup        func()
		buildRequest func() (*http.Request, error)
		shouldErr    bool
		check        func(w *httptest.ResponseRecorder)
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
			buildRequest: func() (*http.Request, error) {
				file, err := os.Open(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer file.Close()

				var body = &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile(files.FileFormKey, path.Base(file.Name()))
				_, err = io.Copy(part, file)
				suite.Require().NoError(err)
				writer.Close()

				req, err := http.NewRequest("POST", "/media", body)
				suite.Require().NoError(err)
				req.Header.Add("Content-Type", writer.FormDataContentType())
				req.Header.Add("Authorization", "Bearer token")

				return req, nil
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
			buildRequest: func() (*http.Request, error) {
				file, err := os.Open(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer file.Close()

				var body = &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile(files.FileFormKey, path.Base(file.Name()))
				_, err = io.Copy(part, file)
				suite.Require().NoError(err)
				writer.Close()

				req, err := http.NewRequest("POST", "/media", body)
				suite.Require().NoError(err)
				req.Header.Add("Content-Type", writer.FormDataContentType())
				req.Header.Add("Authorization", "Bearer token")

				return req, nil
			},
			shouldErr: false,
			check: func(w *httptest.ResponseRecorder) {
				// Make sure the response is well formatted
				resBz, err := io.ReadAll(w.Body)
				suite.Require().NoError(err)

				var res files.UploadFileResponse
				err = json.Unmarshal(resBz, &res)
				suite.Require().NoError(err)
				suite.Require().NotEmpty(res.URL)

				// Make sure the image hash has been saved properly
				var hash string
				err = suite.db.SQL.QueryRow(`SELECT hash FROM images_hashes WHERE image_url = $1`, res.URL).Scan(&hash)
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

			req, err := tc.buildRequest()
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			suite.r.ServeHTTP(w, req)

			if tc.shouldErr {
				suite.Require().NotEqual(http.StatusOK, w.Code)
			} else {
				suite.Require().Equal(http.StatusOK, w.Code)
			}

			if tc.check != nil {
				tc.check(w)
			}
		})
	}
}

func (suite *FilesAPITestSuite) TestDownloadMedia() {
	testCases := []struct {
		name         string
		setup        func() (url string)
		buildRequest func(url string) (*http.Request, error)
		shouldErr    bool
		check        func(w *httptest.ResponseRecorder)
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
				res, err := testutils.GetTestImage()
				suite.Require().NoError(err)

				outFile, err := os.Create(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer outFile.Close()

				_, err = outFile.Write(res)
				suite.Require().NoError(err)

				// Upload the test image
				file, err := os.Open(path.Join(suite.tempDir, "temp_file.jpeg"))
				suite.Require().NoError(err)
				defer file.Close()

				var body = &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile(files.FileFormKey, path.Base(file.Name()))
				_, err = io.Copy(part, file)
				suite.Require().NoError(err)
				writer.Close()

				req, err := http.NewRequest("POST", "/media", body)
				suite.Require().NoError(err)
				req.Header.Add("Content-Type", writer.FormDataContentType())
				req.Header.Add("Authorization", "Bearer token")

				w := httptest.NewRecorder()
				suite.r.ServeHTTP(w, req)
				suite.Require().Equal(http.StatusOK, w.Code)

				var uploadResponse files.UploadFileResponse
				resBz, err := io.ReadAll(w.Body)
				suite.Require().NoError(err)
				err = json.Unmarshal(resBz, &uploadResponse)
				suite.Require().NoError(err)

				return uploadResponse.URL
			},
			buildRequest: func(url string) (*http.Request, error) {
				return http.NewRequest("GET", url, nil)
			},
			shouldErr: false,
			check: func(w *httptest.ResponseRecorder) {
				expectedBz, err := testutils.GetTestImage()
				suite.Require().NoError(err)

				resBz, err := io.ReadAll(w.Body)
				suite.Require().NoError(err)
				suite.Require().Equal(expectedBz, resBz)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			var url string
			if tc.setup != nil {
				url = tc.setup()
			}

			req, err := tc.buildRequest(url)
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			suite.r.ServeHTTP(w, req)

			if tc.shouldErr {
				suite.Require().NotEqual(http.StatusOK, w.Code)
			} else {
				suite.Require().Equal(http.StatusOK, w.Code)
			}

			if tc.check != nil {
				tc.check(w)
			}
		})
	}
}
