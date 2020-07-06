package apiClient

import (
	"errors"
	"strconv"

	"github.com/xiaomLee/grpc-end/client"

	json "github.com/json-iterator/go"
)

func GetEntrustList(serverAddress string, orgId int, coinType string) ([]string, error) {
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["coinType"] = coinType

	data, err := client.CallEndApi(serverAddress, "engine", "entrustList", params)
	if err != nil {
		return nil, err
	}

	type Resp struct {
		Success bool     `json:"success"`
		PayLoad []string `json:"payload"`
		Err     struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	resp := &Resp{}
	if err = json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Err.Message)
	}

	return resp.PayLoad, nil
}
