package icon

import (
	"context"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/sivchari/gotwtr"
)

type FetchClient interface {
	Fetch(ctx context.Context, twitterID string) (image.Image, error)
}

type fetchClient struct {
	twitterClient *gotwtr.Client
	transport     *http.Client
}

func NewFetchClient(bearerToken string, transport *http.Client) FetchClient {
	client := gotwtr.New(bearerToken, gotwtr.WithHTTPClient(transport))
	return &fetchClient{
		twitterClient: client,
		transport:     transport,
	}
}

func (f *fetchClient) Fetch(ctx context.Context, twitterID string) (image.Image, error) {
	user, err := f.twitterClient.RetrieveSingleUserWithUserName(ctx, twitterID, &gotwtr.RetrieveUserOption{
		Expansions:  nil,
		TweetFields: nil,
		UserFields: []gotwtr.UserField{
			gotwtr.UserFieldProfileImageURL,
		},
	})
	if err != nil {
		return nil, err
	}
	originalIconURL := strings.ReplaceAll(user.User.ProfileImageURL, "_normal", "")
	resp, err := f.transport.Get(originalIconURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}
