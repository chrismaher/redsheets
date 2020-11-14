package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	SecretFile string
	TokenFile  string
	*oauth2.Token
	*sheets.Service
}

func (c *Client) getToken(config *oauth2.Config) error {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Open the following link in your browser and then type the authorization code: \n%v\n", url)

	var err error
	var authCode string

	if _, err = fmt.Scan(&authCode); err != nil {
		return fmt.Errorf("Unable to read authorization code: %v", err)
	}

	if c.Token, err = config.Exchange(context.TODO(), authCode); err != nil {
		return fmt.Errorf("Unable to retrieve token from web: %v", err)
	}

	return nil
}

func (c *Client) readToken() error {
	f, err := os.Open(c.TokenFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&c.Token)
}

func (c *Client) saveToken() error {
	fmt.Printf("Saving credential file to: %s\n", c.TokenFile)

	f, err := os.OpenFile(c.TokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(c.Token)

	return nil
}

func (c *Client) getClient(config *oauth2.Config) (*http.Client, error) {
	if err := c.readToken(); err != nil {
		if err := c.getToken(config); err != nil {
			return nil, err
		}
		c.saveToken()
	}

	return config.Client(context.Background(), c.Token), nil
}

func (c *Client) Authorize() error {
	buffer, err := ioutil.ReadFile(c.SecretFile)
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(buffer, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return err
	}

	client, err := c.getClient(config)
	if err != nil {
		return err
	}

	if c.Service, err = sheets.New(client); err != nil {
		return err
	}

	return nil
}
