package service

import (
	"os"
)

type Download interface {
	DownloadMp4(url string) (*os.File, string, error)
	DownloadMp3(url string) (*os.File, string, error)
	DbownloadMp3(url string) ([]byte, error)
}

type Service struct {
	Download
}

type StateMachine interface {
}

func NewService() *Service {
	return &Service{
		Download: NewDownloadService(),
	}
}
