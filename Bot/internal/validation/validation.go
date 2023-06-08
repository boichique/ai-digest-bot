package validation

import (
	"fmt"
	"regexp"
)

const (
	YoutubeChannelRe = `^@[a-zA-Z0-9_-]{1,}`
	YoutubeVideoRe   = `^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/watch\?v=.{1,}`
)

func ValidateLink(pattern string, source string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(source) {
		return fmt.Errorf("incorrect source")
	}

	return nil
}
