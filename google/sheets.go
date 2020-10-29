package google

import "fmt"

func (s *Service) GetRange(spreadsheetId string, readRange string) ([][]interface{}, error) {
	fmt.Printf("Getting range %s\n", readRange)
	resp, err := s.Client.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}
