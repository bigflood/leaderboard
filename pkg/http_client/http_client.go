package http_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/bigflood/leaderboard/api"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func New(endpoint string) *Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	httpTransport := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		IdleConnTimeout:       3 * time.Minute,
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       20,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   30 * time.Second,
	}

	return &Client{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

func (client *Client) doReq(ctx context.Context, method, path string, data interface{}) error {
	req, err := http.NewRequest(method, client.endpoint+path, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		type MessageData struct {
			Message string
		}
		msgData := MessageData{}
		json.NewDecoder(resp.Body).Decode(&msgData)
		msg := msgData.Message
		if msg == "" {
			msg = resp.Status
		}
		return errors.New(msg)
	}

	return json.NewDecoder(resp.Body).Decode(data)
}

func (client *Client) UserCount(ctx context.Context) (int, error) {
	type UserCountData struct {
		Count int
	}

	data := &UserCountData{}

	err := client.doReq(ctx, http.MethodGet, "/usercount", data)
	if err != nil {
		return 0, err
	}

	return data.Count, nil
}

func (client *Client) GetUser(ctx context.Context, userId string) (api.User, error) {
	data := api.User{}

	path := fmt.Sprintf("/users/%s", userId)
	err := client.doReq(ctx, http.MethodGet, path, &data)
	return data, err
}

func (client *Client) SetUser(ctx context.Context, userId string, score int) error {
	type Data struct {
	}
	data := Data{}

	path := fmt.Sprintf("/users/%s?score=%v", userId, score)
	err := client.doReq(ctx, http.MethodPut, path, &data)
	return err
}

func (client *Client) GetRanks(ctx context.Context, rank, count int) ([]api.User, error) {
	data := []api.User{}

	path := fmt.Sprintf("/ranks?rank=%v&count=%v", rank, count)
	err := client.doReq(ctx, http.MethodGet, path, &data)
	return data, err
}
