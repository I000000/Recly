package meili

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/I000000/recly/internal/domain"
)

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
