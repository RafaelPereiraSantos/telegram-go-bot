package ziputils

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

const compressedExtension = ".zip"

// DownloadFilesAndCompress, it receives a list of files to download to compress and return the name of the
// compressed file.
func DownloadFilesAndCompress(urls []string) (string, error) {
	filesTocompressName := make([]string, 0, len(urls)+1)

	for _, url := range urls {
		fileName, err := downloadFile(url)

		if err != nil {
			return "", err
		}

		filesTocompressName = append(filesTocompressName, fileName)
	}

	zipFileName, err := compressMultipleFiles(filesTocompressName)

	if err != nil {
		return "", err
	}

	return zipFileName, nil
}

// downloadFile, downloads a file given the URL.
func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	fileName := uuid.NewString() + extractContentTypeExtension(resp.Header)

	file, err := os.Create(fileName)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func extractContentTypeExtension(header http.Header) string {
	cType := header.Get("content-type")

	if cType == "" {
		return ""
	}

	typeParts := strings.Split(cType, "/")

	return typeParts[1]
}

// compressMultipleFiles, it receives a list of files to compress alongside and return the name of the newly zip file
// containing all the files compressed.
func compressMultipleFiles(filesToCompressNames []string) (string, error) {
	zipFileName := uuid.NewString() + compressedExtension

	archive, err := os.Create(zipFileName)

	if err != nil {
		return "", err
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	defer zipWriter.Close()

	for _, fileName := range filesToCompressNames {

		fileToCompress, err := os.Open(fileName)

		if err != nil {
			return "", err
		}

		writter, err := zipWriter.Create(fileName)

		if err != nil {
			return "", err
		}

		if _, err := io.Copy(writter, fileToCompress); err != nil {
			return "", err
		}

		fileToCompress.Close()

		removeFile(fileName)
	}

	return zipFileName, nil
}

func removeFile(fileName string) error {
	return os.Remove(fileName)
}
