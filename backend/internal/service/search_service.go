package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/I000000/recly/internal/domain"
)

type SearchService struct {
	meiliURL   string
	meiliKey   string
	httpClient *http.Client
}

func NewSearchService(meiliURL, meiliKey string) *SearchService {
	return &SearchService{
		meiliURL:   meiliURL,
		meiliKey:   meiliKey,
		httpClient: &http.Client{},
	}
}

// ---------- helpers ----------

func normalizeGenres(raw interface{}) []string {
	switch v := raw.(type) {
	case string:
		parts := strings.Split(v, ",")
		res := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				res = append(res, p)
			}
		}
		return res
	case []interface{}:
		res := make([]string, 0, len(v))
		for _, item := range v {
			switch s := item.(type) {
			case string:
				res = append(res, s)
			case map[string]interface{}:
				if name, ok := s["children"]; ok {
					if str, ok := name.(string); ok {
						res = append(res, str)
					} else {
						res = append(res, fmt.Sprint(name))
					}
				} else {
					res = append(res, fmt.Sprint(s))
				}
			default:
				res = append(res, fmt.Sprint(item))
			}
		}
		return res
	case []string:
		return v
	}
	return nil
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	if lower == "nan" || lower == "null" || lower == "<nil>" || s == "" {
		return ""
	}
	return s
}

func parseYear(dateStr string) int {
	dateStr = cleanString(dateStr)
	if dateStr == "" {
		return 0
	}
	if idx := strings.Index(dateStr, "-"); idx != -1 {
		dateStr = dateStr[:idx]
	}
	if y, err := strconv.Atoi(dateStr); err == nil {
		return y
	}
	return 0
}

func parseEuropeanFloat(s string) float64 {
	s = cleanString(s)
	if s == "" {
		return 0
	}
	normalized := strings.Replace(s, ",", ".", 1)
	if f, err := strconv.ParseFloat(normalized, 64); err == nil {
		return f
	}
	return 0
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return 0
		}
		return val
	case string:
		return parseEuropeanFloat(val)
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case string:
		cleaned := cleanString(val)
		if cleaned == "" {
			return 0
		}
		if i, err := strconv.Atoi(cleaned); err == nil {
			return i
		}
		// на случай "77,0" – берём целую часть
		if f := parseEuropeanFloat(cleaned); f != 0 {
			return int(f)
		}
	}
	return 0
}

// ---------- raw hit ----------

type rawHit struct {
	ID              interface{} `json:"id"`
	Title           string      `json:"title"`
	Type            string      `json:"type"`
	ImageURL        string      `json:"image_url"`
	Genres          interface{} `json:"genres"`
	AverageRating   interface{} `json:"average_rating"`
	VoteAverage     interface{} `json:"vote_average"`
	RatingsCount    interface{} `json:"ratings_count"`
	VoteCount       interface{} `json:"vote_count"`
	PublicationYear interface{} `json:"publication_year"`
	ReleaseDate     interface{} `json:"release_date"`
	Description     string      `json:"description"`
	Overview        string      `json:"overview"`
	Authors         string      `json:"authors"`
	Director        string      `json:"director"`
	Cast            string      `json:"cast"`
	Runtime         interface{} `json:"runtime"`
}

func (r *rawHit) toItemDetail() domain.ItemDetail {
	item := domain.ItemDetail{
		ID:       fmt.Sprint(r.ID),
		Title:    r.Title,
		Type:     r.Type,
		Image:    r.ImageURL,
		Genres:   normalizeGenres(r.Genres),
		Authors:  r.Authors,
		Director: cleanString(r.Director),
		Cast:     cleanString(r.Cast),
		Runtime:  toInt(r.Runtime),
	}

	if r.Type == "book" {
		item.Year = toInt(r.PublicationYear)
		item.Rating = toFloat(r.AverageRating)
		item.Popularity = toInt(r.RatingsCount)
		item.Description = r.Description
	} else { // movie
		releaseStr := cleanString(fmt.Sprint(r.ReleaseDate))
		if releaseStr != "" {
			item.Year = parseYear(releaseStr)
		}
		item.Rating = toFloat(r.VoteAverage)
		item.Popularity = toInt(r.VoteCount)
		item.Description = cleanString(r.Overview)
	}
	return item
}

// ---------- public methods ----------

func (s *SearchService) Search(query string) ([]domain.ItemDetail, error) {
	url := fmt.Sprintf("%s/indexes/items/search", s.meiliURL)
	payload := map[string]interface{}{
		"q":      query,
		"limit":  10,
		"filter": `type = "book" OR type = "movie"`,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.meiliKey)

	resp, err := s.httpClient.Do(req)
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

func (s *SearchService) SearchWithFilters(query, itemType, genre, sort string, limit, offset int) ([]domain.ItemDetail, error) {
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
	url := fmt.Sprintf("%s/indexes/items/search", s.meiliURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.meiliKey)

	resp, err := s.httpClient.Do(req)
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

func (s *SearchService) GetItems(ids []string, itemType string) ([]domain.ItemDetail, error) {
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

	url := fmt.Sprintf("%s/indexes/items/search", s.meiliURL)
	payload := map[string]interface{}{
		"q":      "",
		"limit":  len(ids),
		"filter": filter,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.meiliKey)

	resp, err := s.httpClient.Do(req)
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
