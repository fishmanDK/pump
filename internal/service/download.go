package service

import (
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

// type YoutubeClient struct{
// 	YoutubeClient youtube.Client
// }

type DownloadService struct {
	YoutubeClient youtube.Client
}

func NewDownloadService() *DownloadService {
	return &DownloadService{
		YoutubeClient: youtube.Client{},
	}
}
func getVideoID(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("error: bad url")
	}

	// Создаем регулярное выражение для извлечения идентификатора видео из URL
	re := regexp.MustCompile(`(?:youtu\.be\/|youtube\.com\/(?:shorts\/|watch\?(?:.*&)?v=|(?:embed|v)\/))([\w-]{11})`)

	// Извлекаем идентификатор видео из URL
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", errors.New("error: invalid youtube url")
	}
	videoID := matches[1]

	return videoID, nil
}

func (d *DownloadService) DownloadMp4(url string) (*os.File, string, error) {
	videoID, err := getVideoID(url)
	if err != nil {
		return nil, "", err
	}
	client := d.YoutubeClient

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}
	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	outputDir := "videos"
	videoTitle := video.Title[:]
	filePath := path.Join(outputDir, fmt.Sprintf("%s.mp4", videoTitle))
	file, err := os.Create(filePath)
	defer file.Close()
	io.Copy(file, stream)
	//if err != nil {
	//	panic(err)
	//}

	return file, videoTitle, nil

}

func (d *DownloadService) DownloadMp3(url string) (*os.File, string, error) {
	videoID, err := getVideoID(url)
	if err != nil {
		return nil, "", err
	}
	client := d.YoutubeClient

	video, err := client.GetVideo(videoID)

	if err != nil {
		panic(err)
	}
	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	outputDir := "videos"
	//videoTitle := video.Title[:]
	videoTitle := strings.ReplaceAll(video.Title, " ", "_")
	filePath := path.Join(outputDir, fmt.Sprintf("%s.mp3", videoTitle))
	file, err := os.Create(filePath)
	defer file.Close()
	io.Copy(file, stream)
	//if err != nil {
	//	panic(err)
	//}

	return file, filePath, nil

}

func (d *DownloadService) DbownloadMp3(url string) ([]byte, error) {
	videoID, err := getVideoID(url)
	if err != nil {
		return nil, err
	}
	client := d.YoutubeClient

	video, err := client.GetVideo(videoID)

	if err != nil {
		panic(err)
	}
	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	fileBytes, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil

}
