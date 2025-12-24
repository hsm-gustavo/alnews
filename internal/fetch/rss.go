package fetch

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hsm-gustavo/alnews/internal/model"
)

const DefaultURL = "https://archlinux.org/feeds/news"

func FetchRSS(url string) (model.RSS, error) {
	if url == "" {
		url = DefaultURL
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET",url, nil)
	if err != nil {
		return model.RSS{}, err
	}

	req.Header.Set("User-Agent", "alnews-cli/1.0 (+github.com/hsm-gustavo/alnews)")

	resp, err := client.Do(req)
	if err != nil {
		return model.RSS{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.RSS{}, fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.RSS{}, err
	}

	var feed model.RSS
	if err := xml.Unmarshal(body, &feed); err != nil {
		return model.RSS{}, fmt.Errorf("failed to decode xml: %w", err)
	}

	return feed, nil
}