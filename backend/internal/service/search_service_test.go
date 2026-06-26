package service

import (
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSearchService_Search(t *testing.T) {
	mockClient := mocks.NewClient(t)
	expected := []domain.ItemDetail{{ID: "1", Title: "Test"}}
	mockClient.On("Search", "test").Return(expected, nil)

	svc := NewSearchService(mockClient)
	results, err := svc.Search("test")

	assert.NoError(t, err)
	assert.Equal(t, expected, results)
	mockClient.AssertExpectations(t)
}

func TestSearchService_SearchWithFilters(t *testing.T) {
	mockClient := mocks.NewClient(t)
	expected := []domain.ItemDetail{{ID: "2", Title: "Filtered"}}
	mockClient.On("SearchWithFilters", "query", "book", "fantasy", "rating:desc", 10, 0).Return(expected, nil)

	svc := NewSearchService(mockClient)
	results, err := svc.SearchWithFilters("query", "book", "fantasy", "rating:desc", 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, expected, results)
	mockClient.AssertExpectations(t)
}

func TestSearchService_GetItems(t *testing.T) {
	mockClient := mocks.NewClient(t)
	expected := []domain.ItemDetail{{ID: "3", Title: "Item"}}
	mockClient.On("GetItems", []string{"3"}, "book").Return(expected, nil)

	svc := NewSearchService(mockClient)
	results, err := svc.GetItems([]string{"3"}, "book")

	assert.NoError(t, err)
	assert.Equal(t, expected, results)
	mockClient.AssertExpectations(t)
}

func TestSearchService_GetGenres(t *testing.T) {
	mockClient := mocks.NewClient(t)
	expected := []string{"fantasy", "adventure"}
	mockClient.On("GetGenres", "book").Return(expected, nil)

	svc := NewSearchService(mockClient)
	genres, err := svc.GetGenres("book")

	assert.NoError(t, err)
	assert.Equal(t, expected, genres)
	mockClient.AssertExpectations(t)
}

func TestSearchService_ErrorHandling(t *testing.T) {
	mockClient := mocks.NewClient(t)
	mockClient.On("Search", "error").Return(nil, assert.AnError)

	svc := NewSearchService(mockClient)
	results, err := svc.Search("error")

	assert.Error(t, err)
	assert.Nil(t, results)
	mockClient.AssertExpectations(t)
}
