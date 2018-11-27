package eth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/invisible73/eth/pkg/types"
)

type client struct {
	url string
}

type rpc struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     string        `json:"id"`
}

type baseResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Error   struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type listAccounts struct {
	baseResponse
	Result types.Accounts `json:"result"`
}

type sendTransaction struct {
	baseResponse
	Result string `json:"result"`
}

type tx struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"value"`
}

func (c *client) SendTransaction(transaction types.Transaction) (string, error) {
	data, err := json.Marshal(rpc{Method: "personal_sendTransaction", Params: []interface{}{tx{transaction.From, transaction.To, fmt.Sprintf("0x%x", transaction.Amount)}, "123"}})
	if err != nil {
		return "", err
	}

	body, err := c.makeRequest(data)
	if err != nil {
		return "", err
	}

	var v sendTransaction

	if err = json.Unmarshal(body, &v); err != nil {
		return "", err
	}

	if v.Error.Message != "" {
		return "", errors.New(v.Error.Message)
	}

	return v.Result, nil
}

func (c *client) ListAccounts() (types.Accounts, error) {
	data, err := json.Marshal(rpc{Method: "personal_listAccounts", Params: []interface{}{}})
	if err != nil {
		return nil, err
	}

	body, err := c.makeRequest(data)
	if err != nil {
		return nil, err
	}

	var v listAccounts

	if err = json.Unmarshal(body, &v); err != nil {
		return nil, err
	}

	if v.Error.Message != "" {
		return nil, errors.New(v.Error.Message)
	}

	return v.Result, nil
}

func (c *client) makeRequest(data []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", `application/json`)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// NewClient ...
func NewClient(url string) (Client, error) {
	if url == "" {
		return nil, errors.New("empty url")
	}

	return &client{
		url: url,
	}, nil
}
