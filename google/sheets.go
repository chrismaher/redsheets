package google

import "fmt"

func (c *Client) GetRange(spreadsheetId string, readRange string) ([][]interface{}, error) {
	fmt.Printf("Getting range %s\n", readRange)
	resp, err := c.Service.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}
