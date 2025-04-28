package docker

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"unicode"
)

type Image struct {
	ID         string
	Repository string
	Tag        string
	Size       uint64 // in bytes
	CreatedAt  time.Time

	// from more advanced processing
	RepoDigests []string
	Used        bool
}

func (i *Image) toJSON() ([]byte, error) {
	return json.Marshal(i)
}

type TempImage struct {
	ID         string
	Repository string
	Tag        string
	Size       string
	CreatedAt  string
}

func (i *TempImage) ToImage() (image Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", i.CreatedAt)
	if err != nil {
		return Image{}, err
	}

	parsedSize, err := ParseSize(i.Size)
	if err != nil {
		return Image{}, err
	}

	image = Image{
		ID:         i.ID,
		Repository: i.Repository,
		Tag:        i.Tag,
		Size:       parsedSize,
		CreatedAt:  t,
	}

	return image, nil
}

func ParseSize(size string) (uint64, error) {
	multipliers := map[string]uint64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
		"PB": 1024 * 1024 * 1024 * 1024 * 1024,
	}

	sizeIndex := firstNonNumeric(size)
	if sizeIndex == -1 {
		return 0, nil
	}

	sizeValue := size[:sizeIndex]
	sizeUnit := size[sizeIndex:]

	sizeValueFloat, err := strconv.ParseFloat(sizeValue, 64)
	if err != nil {
		return 0, err
	}

	multiplier, ok := multipliers[sizeUnit]
	if !ok {
		return 0, errors.New("invalid size unit")
	}

	sizeBytes := uint64(sizeValueFloat * float64(multiplier))
	return sizeBytes, nil
}

func firstNonNumeric(s string) int {
	for i, r := range s {
		if unicode.IsLetter(r) {
			return i
		}
	}
	return -1
}
