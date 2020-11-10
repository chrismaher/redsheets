package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	File string
	*sheets.Service
	*oauth2.Token
}

func (c *Client) getToken(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var err error
	var authCode string

	if _, err = fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	if c.Token, err = config.Exchange(context.TODO(), authCode); err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
}

func (c *Client) readToken() error {
	f, err := os.Open(c.File)
	if err != nil {
		return err
	}
	defer f.Close()

	c.Token = new(oauth2.Token)
	if err := json.NewDecoder(f).Decode(c.Token); err != nil {
		return err
	}

	return nil
}

func (c *Client) saveToken() {
	fmt.Printf("Saving credential file to: %s\n", c.File)
	f, err := os.OpenFile(c.File, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(c.Token)
}

func (c *Client) getClient(config *oauth2.Config) *http.Client {
	if err := c.readToken(); err != nil {
		c.getToken(config)
		c.saveToken()
	}
	return config.Client(context.Background(), c.Token)
}

func (c *Client) Authorize() error {
	buffer, err := ioutil.ReadFile(c.File)
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(buffer, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return err
	}

	if c.Service, err = sheets.New(c.getClient(config)); err != nil {
		return err
	}

	return nil
}
