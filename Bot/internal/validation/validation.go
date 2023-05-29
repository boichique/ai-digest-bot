package validation

import (
	"fmt"
	"regexp"
)

const (
	youtubeRe = `^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$`
)

func ValidateYoutube(source string) error {
	re := regexp.MustCompile(youtubeRe)
	if !re.MatchString(source) {
		return fmt.Errorf("Incorrect source")
	}

	return nil
}
