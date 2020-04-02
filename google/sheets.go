package google

import (
	"log"
)

func (s *Service) GetRange(spreadsheetId string, readRange string) ([][]interface{}, error) {
	log.Printf("Getting range %s", readRange)
	resp, err := s.Client.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}
