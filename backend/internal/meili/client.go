package meili

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/I000000/recly/internal/domain"
)

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

func (c *Client) Search(query string) ([]domain.ItemDetail, error) {
	url := fmt.Sprintf("%s/indexes/items/search", c.baseURL)
	payload := map[string]interface{}{
		"q":      query,
		"limit":  10,
		"filter": `type = "book" OR type = "movie"`,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var result struct {
		Hits []rawHit `json:"hits"`
	}
	json.Unmarshal(respBytes, &result)

	items := make([]domain.ItemDetail, len(result.Hits))
	for i, hit := range result.Hits {
		items[i] = hit.toItemDetail()
	}
	return items, nil
}

func (c *Client) SearchWithFilters(query, itemType, genre, sort string, limit, offset int) ([]domain.ItemDetail, error) {
	filters := []string{}
	if itemType != "" && itemType != "all" {
		filters = append(filters, fmt.Sprintf("type = \"%s\"", itemType))
	}
	if genre != "" {
		filters = append(filters, fmt.Sprintf("genres = \"%s\"", strings.ToLower(genre)))
	}
	filter := strings.Join(filters, " AND ")

	if sort == "" {
		switch itemType {
		case "book":
			sort = "ratings_count:desc"
		case "movie":
			sort = "vote_count:desc"
		default:
			sort = ""
		}
	}

	payload := map[string]interface{}{
		"q":      query,
		"limit":  limit,
		"offset": offset,
		"filter": filter,
	}
	if sort != "" {
		payload["sort"] = []string{sort}
	}

	bodyBytes, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/indexes/items/search", c.baseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var result struct {
		Hits []rawHit `json:"hits"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, err
	}

	items := make([]domain.ItemDetail, len(result.Hits))
	for i, hit := range result.Hits {
		items[i] = hit.toItemDetail()
	}
	return items, nil
}

func (c *Client) GetItems(ids []string, itemType string) ([]domain.ItemDetail, error) {
	if len(ids) == 0 {
		return []domain.ItemDetail{}, nil
	}

	filterParts := make([]string, len(ids))
	for i, id := range ids {
		filterParts[i] = fmt.Sprintf("id = %q", id)
	}

	var filter string
	if itemType == "all" || itemType == "" {
		filter = strings.Join(filterParts, " OR ")
	} else {
		filter = fmt.Sprintf("type = %q AND (%s)", itemType, strings.Join(filterParts, " OR "))
	}

	url := fmt.Sprintf("%s/indexes/items/search", c.baseURL)
	payload := map[string]interface{}{
		"q":      "",
		"limit":  len(ids),
		"filter": filter,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var result struct {
		Hits []rawHit `json:"hits"`
	}
	json.Unmarshal(respBytes, &result)

	items := make([]domain.ItemDetail, len(result.Hits))
	for i, hit := range result.Hits {
		items[i] = hit.toItemDetail()
	}
	return items, nil
}

func (c *Client) GetGenres(itemType string) ([]string, error) {
	filter := ""
	if itemType != "" && itemType != "all" {
		filter = fmt.Sprintf("type = \"%s\"", itemType)
	}
	payload := map[string]interface{}{
		"q":      "",
		"filter": filter,
		"facets": []string{"genres"},
		"limit":  1,
	}
	bodyBytes, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/indexes/items/search", c.baseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("meilisearch returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		FacetDistribution struct {
			Genres map[string]int `json:"genres"`
		} `json:"facetDistribution"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	genres := make([]string, 0, len(result.FacetDistribution.Genres))
	for g := range result.FacetDistribution.Genres {
		genres = append(genres, g)
	}
	sort.Strings(genres)
	return genres, nil
}
