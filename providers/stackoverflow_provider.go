package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"mcp-server/config"
)

// To get the data we want (comments, answers, etc.), we need a custom filter.
// This one was created on https://api.stackexchange.com/docs/create-filter
// This contains includes: .backoff, .error_id, .error_message, .error_name, .has_more, .items, .quota_max, .quota_remaining,
// answer.body, answer.comments, answer.score,
// comment.body,
// question.answers, question.body, question.comments, question.score, question.title
const stackExchangeAPIFilter = "!nKzQURF6Y5"

type StackOverflowProvider struct {
	apiKey     string
	client     *http.Client
	apiBaseURL string
}

// Structs to unmarshal the Stack Exchange API response
type StackAPISearchResponse struct {
	Items []Question `json:"items"`
}
type Question struct {
	Title    string    `json:"title"`
	Score    int       `json:"score"`
	Body     string    `json:"body"`
	Answers  []Answer  `json:"answers"`
	Comments []Comment `json:"comments"`
}
type Answer struct {
	Score    int       `json:"score"`
	Body     string    `json:"body"`
	Comments []Comment `json:"comments"`
}
type Comment struct {
	Score int    `json:"score"`
	Body  string `json:"body"`
}

// Bootstraps the Stackoverflow Provider
func NewStackOverflowProvider(cfg *config.StackExchangeConfig) *StackOverflowProvider {
	return &StackOverflowProvider{
		apiKey:     cfg.APIKey,
		client:     &http.Client{},
		apiBaseURL: "https://api.stackexchange.com/2.3",
	}
}

// Searches StackOverflow and maps the response to spec
func (p *StackOverflowProvider) Search(query string) (*StackAPISearchResponse, error) {
	reqURL, _ := url.Parse(p.apiBaseURL + "/search/advanced")
	q := reqURL.Query()
	q.Set("order", "desc")
	q.Set("sort", "relevance")
	q.Set("site", "providers")
	q.Set("pagesize", "10")
	q.Set("q", query)
	q.Set("filter", stackExchangeAPIFilter)
	if p.apiKey != "" {
		q.Set("key", p.apiKey)
	}
	reqURL.RawQuery = q.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to call Stack Exchange API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stack Exchange API returned non-200 status: %s", resp.Status)
	}

	var searchResult StackAPISearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode Stack Exchange API response: %w", err)
	}

	return &searchResult, nil
}