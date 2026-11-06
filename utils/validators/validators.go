package validators

import (
	"Quote/utils/arrayUtils"
	"net/http"
	"regexp"
)

func ValidateEmail(email string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return pattern.MatchString(email)
}

// Phone number must start with 00 followed by at least 11 digits
func ValidatePhone(phone string) bool {
	pattern := regexp.MustCompile(`^00\d{3}\d{8,}$`)
	return pattern.MatchString(phone)
}

func ValidateEmails(emails []string) bool {
	for _, email := range emails {
		if !ValidateEmail(email) {
			return false
		}
	}
	return true
}

func ValidatePhones(phones []string) bool {
	for _, phone := range phones {
		if !ValidatePhone(phone) {
			return false
		}
	}
	return true
}

func ValidateImage(imageBuffer []byte) bool {
	filetype := http.DetectContentType(imageBuffer)

	switch filetype {
	case "image/jpeg", "image/jpg", "image/gif", "image/png":
		return true
	default:
		return false
	}
}

func ValidatePdf(pdfBuffer []byte) bool {
	filetype := http.DetectContentType(pdfBuffer)

	switch filetype {
	case "application/pdf":
		return true
	default:
		return false
	}
}

func ValidateAudioExtension(fileExtension string) bool {
	audioExtension := []string{".aa", ".aac", ".aax", ".act", ".aiff", ".alac", ".amr", ".ape", ".au", ".awb", ".dss", ".dvf", ".flac", ".gsm", ".iklax", ".ivs", ".m4a", ".m4b", ".m4p", ".mmf", ".movpkg", ".mp3", ".mpc", ".msv", ".nmf", ".ogg, .oga, .mogg", ".opus", ".ra, .rm", ".rf64", ".sln", ".tta", ".voc", ".vox", ".wav", ".wma", ".wv", ".8svx", ".cda"}

	return arrayUtils.ContainsValue(audioExtension, fileExtension)
}
