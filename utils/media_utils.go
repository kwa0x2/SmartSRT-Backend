package utils

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/alfg/mp4"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
)

func IsValidMediaFile(fileType string) bool {
	switch fileType {
	case ".mp4", ".mp3", ".wav":
		return true
	default:
		return false
	}
}

func GetMP4Duration(file io.Reader) (float64, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return 0, err
	}

	mp4file, err := mp4.OpenFromBytes(buf.Bytes())
	if err != nil {
		return 0, err
	}

	duration := float64(mp4file.Moov.Mvhd.Duration) / float64(mp4file.Moov.Mvhd.Timescale)
	return math.Floor(duration), nil
}

func GetMP3Duration(file io.Reader) (float64, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return 0, err
	}

	decoder, err := mp3.NewDecoder(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}

	sampleRate := decoder.SampleRate()
	if sampleRate <= 0 {
		return 0, fmt.Errorf("invalid sample rate")
	}

	duration := float64(decoder.Length()) / float64(4) / float64(sampleRate)
	return math.Floor(duration), nil
}

func GetWAVDuration(file io.Reader) (float64, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return 0, err
	}

	decoder := wav.NewDecoder(bytes.NewReader(buf.Bytes()))
	if !decoder.IsValidFile() {
		return 0, fmt.Errorf("invalid WAV file")
	}

	decoder.ReadInfo()

	audioData, err := decoder.FullPCMBuffer()
	if err != nil {
		return 0, err
	}

	duration := float64(audioData.NumFrames()) / float64(decoder.SampleRate)
	return math.Floor(duration), nil
}

func GetMediaDuration(file io.Reader, fileType string) (float64, error) {
	switch fileType {
	case ".mp3":
		return GetMP3Duration(file)
	case ".mp4":
		return GetMP4Duration(file)
	case ".wav":
		return GetWAVDuration(file)
	default:
		return 0, fmt.Errorf("unsupported file type: %s", fileType)
	}
}
