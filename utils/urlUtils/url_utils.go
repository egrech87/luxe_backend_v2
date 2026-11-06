package urlUtils

import "net/http"

// Taken from https://github.com/cyfrost/url-expander/blob/master/local_server.go
func Unshorten(url string) (unshortenedUrl string, err error) {
	unshortenedUrl = url

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			unshortenedUrl = req.URL.String()
			return nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return unshortenedUrl, nil
}
