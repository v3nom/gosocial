package fbapi

import (
	"errors"
	"net/url"

	"golang.org/x/net/context"

	"io/ioutil"

	"encoding/json"

	"strings"

	"time"

	"google.golang.org/appengine/urlfetch"
)

// AppID facebook app id
var AppID string

// AppRedirectURL redirect URL
var AppRedirectURL string

// AppSecret facebook app secret
var AppSecret string

// APIServerVersion facebook API version
var APIServerVersion = "v2.10"

var apiServer = "https://www.facebook.com/"

// GetLoginURL gets fb login url
func GetLoginURL(scopes []string) string {
	url, _ := url.Parse(getFacebookAPIServerURL() + "/dialog/oauth")
	q := url.Query()
	q.Set("client_id", AppID)
	q.Set("redirect_uri", AppRedirectURL)
	q.Set("scope", strings.Join(scopes, ","))
	url.RawQuery = q.Encode()
	return url.String()
}

// LoginUserWithResponseQuery logins user with code
func LoginUserWithResponseQuery(ctx context.Context, query url.Values) (*UserResponse, error) {
	if query.Get("error") != "" {
		return nil, errors.New(query.Get("error_reason"))
	}
	code := query.Get("code")
	if code == "" {
		return nil, errors.New("No code")
	}

	// Get access token
	response, errConfirm := confirmIdentity(ctx, code)
	if errConfirm != nil {
		return nil, errConfirm
	}
	if response.Error.Type != "" {
		return nil, errors.New(response.Error.Message)
	}

	// Get user info
	userInfo, errInfo := GetUserInfo(ctx, response.AccessToken, []string{"name", "email"})
	if errInfo != nil {
		return nil, errInfo
	}
	if userInfo.Error.Type != "" {
		return nil, errors.New(userInfo.Error.Message)
	}

	userInfo.AccessToken = response.AccessToken
	userInfo.AccessTokenExpiresIn = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)

	return userInfo, nil
}

// GetUserInfo gets user info
func GetUserInfo(ctx context.Context, token string, fields []string) (*UserResponse, error) {
	graphURL, _ := url.Parse(getFacebookAPIServerURL() + "/me")
	q := graphURL.Query()
	q.Set("access_token", token)
	q.Set("fields", strings.Join(fields, ","))
	graphURL.RawQuery = q.Encode()

	resp, err := makeGetRequest(ctx, graphURL.String())
	if err != nil {
		return nil, err
	}
	var decodedResponse UserResponse
	jsonDecover := json.NewDecoder(strings.NewReader(resp))
	errDecode := jsonDecover.Decode(&decodedResponse)
	if errDecode != nil {
		return nil, errDecode
	}
	return &decodedResponse, nil
}

func confirmIdentity(ctx context.Context, code string) (*IdentityResponse, error) {
	tokenURL, _ := url.Parse(getFacebookAPIServerURL() + "/oauth/access_token")
	q := tokenURL.Query()
	q.Add("client_id", AppID)
	q.Add("redirect_uri", AppRedirectURL)
	q.Add("client_secret", AppSecret)
	q.Add("code", code)
	tokenURL.RawQuery = q.Encode()

	data, err := makeGetRequest(ctx, tokenURL.String())
	if err != nil {
		return nil, err
	}

	var decodedResponse IdentityResponse
	jsonDecover := json.NewDecoder(strings.NewReader(data))
	errDecode := jsonDecover.Decode(&decodedResponse)
	if errDecode != nil {
		return nil, errDecode
	}
	return &decodedResponse, nil
}

func makeGetRequest(ctx context.Context, url string) (string, error) {
	client := urlfetch.Client(ctx)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, errBody := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errBody
	}
	return string(data), nil
}

func getFacebookAPIServerURL() string {
	return apiServer + APIServerVersion
}
