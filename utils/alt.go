package utils

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var raccAlts map[string]string

func LoadRaccAlts(fn string) error {
	content, err := os.ReadFile(fn)
	if err != nil {
		return errors.Wrap(err, "reading alt.json")
	}
	if err := json.Unmarshal(content, &raccAlts); err != nil {
		return errors.Wrap(err, "unmarshaling alt.json")
	}
	return nil
}

func GetAlt(index string) string {
	alt, ok := raccAlts[index]
	if !ok {
		return "a raccoon"
	}
	return alt
}

func GetAlti(index int) string {
	return GetAlt(strconv.Itoa(index))
}
