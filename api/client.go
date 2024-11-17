package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"speedrunner/models"
)

const (
	baseURL    = "https://www.speedrun.com/api/v1"
	userAgent  = "speedrun-browser/1.0"
	rateLimit  = 100
	ratePeriod = time.Minute
)

type Client struct {
	httpClient *http.Client
	lastReqs   []time.Time
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		lastReqs: make([]time.Time, 0, rateLimit),
	}
}

// SearchGames searches for games using the bulk API
func (c *Client) SearchGames(query string) ([]models.Game, error) {
	params := url.Values{}
	params.Add("name", query)
	params.Add("_bulk", "yes") // Use bulk mode for better performance

	var resp models.APIResponse
	resp.Data = make([]interface{}, 0)

	err := c.get("/games?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	games := make([]models.Game, 0, len(resp.Data))
	for _, item := range resp.Data {
		if gameData, err := json.Marshal(item); err == nil {
			var game models.Game
			if err := json.Unmarshal(gameData, &game); err == nil {
				games = append(games, game)
			}
		}
	}

	return games, nil
}

// SearchUsers searches for users
func (c *Client) SearchUsers(query string) ([]models.User, error) {
	params := url.Values{}
	params.Add("lookup", query) // Use lookup for comprehensive search

	var resp models.APIResponse
	resp.Data = make([]interface{}, 0)

	err := c.get("/users?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(resp.Data))
	for _, item := range resp.Data {
		if userData, err := json.Marshal(item); err == nil {
			var user models.User
			if err := json.Unmarshal(userData, &user); err == nil {
				users = append(users, user)
			}
		}
	}

	return users, nil
}

func (c *Client) get(endpoint string, v interface{}) error {
	if err := c.checkRateLimit(); err != nil {
		return err
	}

	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 420 {
		return fmt.Errorf("rate limit exceeded")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) checkRateLimit() error {
	now := time.Now()
	cutoff := now.Add(-ratePeriod)

	validReqs := c.lastReqs[:0]
	for _, t := range c.lastReqs {
		if t.After(cutoff) {
			validReqs = append(validReqs, t)
		}
	}
	c.lastReqs = validReqs

	if len(c.lastReqs) >= rateLimit {
		return fmt.Errorf("rate limit exceeded: %d requests per %v", rateLimit, ratePeriod)
	}

	c.lastReqs = append(c.lastReqs, now)
	return nil
}
