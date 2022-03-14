package model

import (
	"io"
	"os"
)

type (
	ReceivedMessage struct {
		ChatId int64
		User   string
		Text   string
	}

	ReplyMessage struct {
		Text  string
		Image *ReplyLocalImage
	}

	ReplyLocalImage struct {
		FileName string
		FilePath string
	}
)

func (ri *ReplyLocalImage) NeedsUpload() bool {
	return true
}

func (ri *ReplyLocalImage) UploadData() (string, io.Reader, error) {
	f, err := os.Open(ri.FilePath)

	if err != nil {
		return "", nil, err
	}

	ri.removeUploadedFile()

	return ri.FileName, f, nil
}

func (ri *ReplyLocalImage) SendData() string {
	return ""
}

func (ri ReplyLocalImage) removeUploadedFile() error {
	return nil
}
