package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BreedService struct {
	client    *http.Client
	cache     []Breed
	cacheTime time.Time
	mutex     sync.RWMutex
	cacheTTL  time.Duration
}

func NewBreedService() *BreedService {
	return &BreedService{
		client:   &http.Client{Timeout: 30 * time.Second},
		cacheTTL: time.Hour, // Cache for 1 hour
	}
}

func (s *BreedService) GetBreeds() ([]Breed, error) {
	s.mutex.RLock()
	if len(s.cache) > 0 && time.Since(s.cacheTime) < s.cacheTTL {
		breeds := make([]Breed, len(s.cache))
		copy(breeds, s.cache)
		s.mutex.RUnlock()
		return breeds, nil
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.cache) > 0 && time.Since(s.cacheTime) < s.cacheTTL {
		breeds := make([]Breed, len(s.cache))
		copy(breeds, s.cache)
		return breeds, nil
	}

	resp, err := s.client.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch breeds from TheCatAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TheCatAPI returned status %d", resp.StatusCode)
	}

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, fmt.Errorf("failed to decode breeds response: %w", err)
	}

	s.cache = breeds
	s.cacheTime = time.Now()

	result := make([]Breed, len(breeds))
	copy(result, breeds)
	return result, nil
}

func (s *BreedService) ValidateBreed(breedName string) (bool, error) {
	breeds, err := s.GetBreeds()
	if err != nil {
		return false, err
	}

	for _, breed := range breeds {
		if breed.Name == breedName {
			return true, nil
		}
	}

	return false, nil
}

func (s *BreedService) GetBreedNames() ([]string, error) {
	breeds, err := s.GetBreeds()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(breeds))
	for i, breed := range breeds {
		names[i] = breed.Name
	}

	return names, nil
}
