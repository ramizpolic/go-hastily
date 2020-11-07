package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/user"
	"strings"
)

// TokenResponse defines token data on authentication request.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenId      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// Credentials is the parent of all OAuth token related activities.
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenId      string `json:"id_token"`
	Type         string `json:"token_type"`
	Endpoint     string `json:"endpoint"`
	Path         string `json:"path"`
	Username     string `json:"username"`
}

// credentialsPath is a shared function that defines where
// the credentials will be saved and loaded from.
func credentialsPath() (string, error) {
	myself, err := user.Current()
	if err != nil {
		return "", err
	}
	return myself.HomeDir + "/.go-hastly.json", nil
}

// Validate if credentials work as a safeguard for other commands.
func (creds *Credentials) Validate() error {
	if creds.AccessToken == "" {
		return errors.New("Invalid login credentials. Please provide valid authentication data.")
	}
	return nil
}

// Save method saves credentials to the file and updates
// path of the object.
func (creds *Credentials) Save() error {

	// obtain path
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	creds.Path = path
	credsJson, _ := json.Marshal(creds)

	// save file
	err = ioutil.WriteFile(creds.Path, credsJson, 0644)
	if err != nil {
		creds.Path = ""
		return err
	}

	return nil
}

// LoadCredentials function loads credentials saved on the system.
// It throws error if there is no such file.
func LoadCredentials() (*Credentials, error) {

	// obtain path
	path, err := credentialsPath()
	if err != nil {
		return nil, err
	}

	// read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// load data into environment
	var credentials Credentials
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		return nil, err
	}

	// verify credentials
	err = credentials.Validate()
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}

// GetCredentials obtains OAuth token required to work with backend API.
func GetCredentials(endpoint string, username string, password string) (*Credentials, error) {

	// create request to obtain oauth token
	data := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response
	var token TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	// serve credentials
	var credentials = Credentials{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenId:      token.TokenId,
		Type:         token.Type,
		Username:     username,
		Endpoint:     endpoint,
		Path:         "",
	}

	// verify credentials
	err = credentials.Validate()
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}
