package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jnelle/MyFitnesspal-API-Golang/app/models"
	"github.com/jnelle/MyFitnesspal-API-Golang/internal/constants"
	"github.com/jnelle/MyFitnesspal-API-Golang/internal/pkg/utils"
	"golang.org/x/net/http2"
)

// MFPClient offers methods to get food products from MyFitnessPal
type MFPClient struct {
	AuthSignKey     []byte                      // signing key which has to be decoded first
	ClientKeys      *models.ClientKeyResponse   // client keys for creating and signing jwts
	AccessToken     *models.AccessTokenResponse // access tokens for communicating with identity api
	IDTokenResponse *models.TokenResponse       // tokens for communicating with general api
	Username        string                      // email/username from mfp acc
	Password        string                      // password from mfp acc
	JWT             string                      // own generated & signed jwt for login
	MFPCallbackCode string                      // callback code for login process
	MFPUserID       string                      // user id for login process
	UserID          string                      // domain user id for api requests
}

func doRequest(req *http.Request) (io.ReadCloser, http.Header, error) {
	client := &http.Client{
		Transport: &http2.Transport{},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.Header{}, errors.New("request failed")
	}
	fmt.Println(resp.StatusCode)

	return resp.Body, resp.Header, nil
}

func (m *MFPClient) getClientKeys() error {
	encodedAuth := utils.EncodeBase64(fmt.Sprintf("%s:%s", constants.ClientID, constants.ClientSecret))
	url := fmt.Sprintf("https://%s/clientKeys", constants.IdentityBaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", constants.UserAgent)

	response, header, err := doRequest(req)
	if err != nil {
		errorMessage := fmt.Sprintf("Couldn't fetch clientkeys\tstatus code: %v\t error: %v", header, err)
		return errors.New(errorMessage)
	}

	bodyBytes, err := io.ReadAll(response)
	if err != nil {
		log.Fatal(err)
	}

	var clientKeys *models.ClientKeyResponse
	err = sonic.Unmarshal(bodyBytes, &clientKeys)
	if err != nil {
		return err
	}

	m.ClientKeys = clientKeys
	m.AuthSignKey, _ = utils.DecodeBase64URL(m.ClientKeys.Embedded.ClientKeys[1].Key.K)

	return nil

}

func (m *MFPClient) getAccessToken() error {
	url := fmt.Sprintf("https://%s/oauth/token", constants.IdentityBaseURL)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", constants.ClientID, constants.ClientSecret)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", "143")
	req.Header.Add("User-Agent", constants.UserAgent)

	resp, header, err := doRequest(req)
	if err != nil {
		errorMessage := fmt.Sprintf("Couldn't get access tokens\tstatus code: %v\t error: %v", header, err)
		return errors.New(errorMessage)
	}

	response, err := io.ReadAll(resp)
	if err != nil {
		log.Fatal(err)
	}

	var accessTokenResponse *models.AccessTokenResponse
	err = sonic.Unmarshal(response, &accessTokenResponse)
	if err != nil {
		return err
	}

	m.AccessToken = accessTokenResponse

	return nil
}

func (m *MFPClient) login() error {
	var err error
	claims := jwt.MapClaims{
		"password": m.Password,
		"username": m.Username,
	}

	claimHeader := map[string]interface{}{
		"kid": m.ClientKeys.Embedded.ClientKeys[1].Key.Kid,
		"alg": m.ClientKeys.Embedded.ClientKeys[1].Key.Alg,
	}

	m.JWT, err = utils.GenJWT(claims, claimHeader, m.AuthSignKey)
	if err != nil {
		return err
	}

	randNum := utils.GenRandomNum(100000000, 586550506)
	bodyString := fmt.Sprintf("client_id=%s&credentials=%s&nonce=%d&redirect_uri=mfp%%3A%%2F%%2Fidentity%%2Fcallback&response_type=code&scope=openid", constants.ClientID, m.JWT, randNum)
	url := fmt.Sprintf("https://%s/oauth/authorize", constants.IdentityBaseURL)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(bodyString))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", constants.UserAgent)
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprint(len(bodyString)))

	// here appears the following error:
	// net/url.Error {Op: "Post", URL: "mfp://identity/callback?code=XXX", Err: error(*errors.errorString) *{s: "http2: unsupported scheme"}}
	// looks like go is validating urls in headers
	_, header, err := doRequest(req)
	if err != nil {
		errorMessage := fmt.Sprintf("Login failed\tstatus code: %v\t error: %v", header, err)
		return errors.New(errorMessage)
	}

	result := strings.SplitAfter(string(header.Get("Location")), constants.CallbackURL)
	m.MFPCallbackCode = result[1]
	err = m.loginCallBack()
	if err != nil {
		return errors.New("sending callbackcode failed")
	}
	return nil
}

func (m *MFPClient) loginCallBack() error {
	bodyString := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=mfp%%3A%%2F%%2Fidentity%%2Fcallback", m.MFPCallbackCode)
	url := fmt.Sprintf("https://%s/oauth/token?auto_create_account_link=false", constants.IdentityBaseURL)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(bodyString))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprint(len(bodyString)))
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.AccessToken.AccessToken))
	req.Header.Add("User-Agent", constants.UserAgent)

	resp, header, err := doRequest(req)
	if err != nil {
		errorMessage := fmt.Sprintf("Sending callbackcode failed\tstatus code: %v\t error: %v", header, err)
		return errors.New(errorMessage)

	}

	response, err := io.ReadAll(resp)
	if err != nil {
		log.Fatal(err)
	}

	var IDTokenResponse *models.TokenResponse
	err = sonic.Unmarshal(response, &IDTokenResponse)
	if err != nil {
		return err
	}

	m.IDTokenResponse = IDTokenResponse

	// err = m.getMFPUserID()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// mfp user id is needed for login process, shouldn't be that interesting for you
// func (m *MFPClient) getMFPUserID() error {
// 	m.UserID = utils.DecodeJWT(m.IDTokenResponse.IDToken)

// 	req := fasthttp.AcquireRequest()
// 	req.SetRequestURI(fmt.Sprintf("/users/%s?fetch_profile=true&fetch_emails=true", m.UserID))
// 	req.URI().SetScheme("https")
// 	req.SetHost(constants.IdentityBaseURL)
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.AccessToken.AccessToken))
// 	req.Header.SetContentType("application/json")
// 	req.Header.SetUserAgent(constants.UserAgent)
// 	req.Header.Add("Accept-Encoding", "application/json")
// 	req.Header.Add("accept", "gzip")

// 	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))
// 	fmt.Println(req)
// 	response, header, err := doRequest(req)
// 	if err != nil || header.StatusCode() != 200 {
// 		errorMessage := fmt.Sprintf("Couldn't send callbackcode\tstatus code: %d\t error: %v", header.StatusCode(), err)
// 		return errors.New(errorMessage)

// 	}

// 	// fmt.Println(response)
// 	var MFPUserIDResponse *models.UserIDResponse
// 	err = sonic.Unmarshal(response, &MFPUserIDResponse)
// 	if err != nil {
// 		return err
// 	}

// 	m.MFPUserID = MFPUserIDResponse.AccountLinks[0].DomainUserID
// 	return nil
// }

// // Any food product that you like
// func (m *MFPClient) SearchFoodWithoutPagination(foodName string) (*models.FoodSearchResponse, error) {
// 	if len(m.Username) == 0 || len(m.Password) == 0 {
// 		return nil, errors.New("mode without authentication active")
// 	}

// 	url := utils.BuildFoodSearchURL(foodName, fmt.Sprint(constants.APIBaseURL+"/v2/search/nutrition"))
// 	req := fasthttp.AcquireRequest()

// 	req.SetRequestURI(url)
// 	req.Header.SetUserAgent(constants.UserAgentAPI)
// 	req.Header.Add("device_id", uuid.NewString())
// 	req.Header.Add("mfp-flow-id", uuid.NewString())
// 	req.Header.Add("api-version", constants.APIVersion)
// 	req.Header.Add("Screen-Density", constants.ScreenDensity)
// 	req.Header.Add("Screen-Height", constants.ScreenHeight)
// 	req.Header.Add("Screen-Width", constants.ScreenWidth)
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.IDTokenResponse.AccessToken))
// 	req.Header.Add("mfp-user-id", m.MFPUserID)
// 	req.Header.Add("mfp-client-id", constants.MFPClientID)
// 	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

// 	response, header, err := doRequest(req)
// 	if err != nil || header.StatusCode() != 200 {
// 		errorMessage := fmt.Sprintf("Couldn't fetch food%s\tstatus code: %d\t error: %v", foodName, header.StatusCode(), err)
// 		return nil, errors.New(errorMessage)
// 	}

// 	var foodResponse *models.FoodSearchResponse
// 	err = sonic.Unmarshal(response, &foodResponse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return foodResponse, nil
// }

// // InitialLoad does the initial login part starts at getting clientkeys till creating
// // JWTs for login
func (m *MFPClient) InitialLoad() error {
	// if len(m.Username) != 0 && len(m.Password) != 0 {
	err := m.getClientKeys()
	if err != nil {
		return err
	}

	err = m.getAccessToken()
	if err != nil {
		return err
	}

	err = m.login()
	if err != nil {
		return err
	}

	log.Println("Using mode without authentication")
	return nil
}

// // Refreshes the user token
// func (m *MFPClient) RefreshToken() error {
// 	if len(m.Username) == 0 || len(m.Password) == 0 {
// 		return errors.New("mode without authentication active")
// 	}

// 	bodyString := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s&user_id=%s", m.IDTokenResponse.RefreshToken, constants.ClientID, constants.ClientSecret, m.MFPUserID)
// 	req := fasthttp.AcquireRequest()
// 	req.SetRequestURI("/oauth/token")
// 	req.URI().SetScheme("https")
// 	req.SetHost(constants.IdentityBaseURL)
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.AccessToken.AccessToken))
// 	req.Header.SetContentType("application/x-www-form-urlencoded")
// 	req.Header.SetContentLength(len(bodyString))
// 	req.SetBodyString(bodyString)
// 	req.Header.SetUserAgent(constants.UserAgent)
// 	req.Header.SetMethodBytes([]byte(fasthttp.MethodPost))

// 	response, header, err := doRequest(req)
// 	if err != nil || header.StatusCode() != 200 {
// 		errorMessage := fmt.Sprintf("Couldn't get refresh token\tstatus code: %d\t error: %v", header.StatusCode(), header.Header())
// 		return errors.New(errorMessage)
// 	}

// 	var accessTokenResponse *models.AccessTokenResponse
// 	err = sonic.Unmarshal(response, &accessTokenResponse)
// 	if err != nil {
// 		return err
// 	}

// 	m.AccessToken = accessTokenResponse

// 	return nil
// }

// // It could be possible that this endpoint won't work in the future
// func (m *MFPClient) SearchFoodWithoutAuth(foodName string, offset int, page int) (*models.FoodSearchResponseWithoutAuth, error) {
// 	path := fmt.Sprintf("/public/nutrition?q=%s&page=%d&per_page=%d", foodName, page, offset)

// 	req := fasthttp.AcquireRequest()
// 	req.SetRequestURI(path)
// 	req.URI().SetScheme("https")
// 	req.SetHost(constants.APIBaseURL)
// 	req.Header.SetUserAgent(constants.UserAgent)
// 	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

// 	response, header, err := doRequest(req)
// 	if err != nil || header.StatusCode() != 200 {
// 		errorMessage := fmt.Sprintf("Request failed\tstatus code: %d\t error: %v", header.StatusCode(), err)
// 		return nil, errors.New(errorMessage)
// 	}

// 	var foodResponse *models.FoodSearchResponseWithoutAuth
// 	err = sonic.Unmarshal(response, &foodResponse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return foodResponse, nil
// }
