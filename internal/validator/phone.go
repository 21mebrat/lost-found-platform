package validator

import (
	"errors"
	"regexp"
	"strings"
)

func NormalizePhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)

	replacer := strings.NewReplacer(
		" ", "",
		"-", "",
		"(", "",
		")", "",
	)
	phone = replacer.Replace(phone)

	switch {
	case strings.HasPrefix(phone, "+251"):
	case strings.HasPrefix(phone, "251"):
		phone = "+" + phone

	case strings.HasPrefix(phone, "0"):
		phone = "+251" + phone[1:]

	default:
		return "", errors.New("invalid phone number format")
	}

	return phone, nil
}

func ValidatePhone(phone string) (string, error) {
	phone, err := NormalizePhone(phone)
	if err != nil {
		return "", err
	}

	// +251 followed by 7 or 9 and then 8 digits
	re := regexp.MustCompile(`^\+251[79]\d{8}$`)

	if !re.MatchString(phone) {
		return "", errors.New("invalid Ethiopian mobile number")
	}

	return phone, nil
}
