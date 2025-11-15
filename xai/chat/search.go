// Package chat provides search and reasoning functionality for xAI SDK.
package chat

import (
	"fmt"
)

// SearchParameters represents search configuration for chat requests.
type SearchParameters struct {
	// Number of search results to include
	count int32
	
	// Search domains to include
	domains []string
	
	// Search recency filter
	recency SearchRecency
}

// SearchRecency represents search recency filters.
type SearchRecency string

const (
	// SearchRecencyDefault uses default recency.
	SearchRecencyDefault SearchRecency = "default"
	
	// SearchRecencyDay includes results from the past day.
	SearchRecencyDay SearchRecency = "day"
	
	// SearchRecencyWeek includes results from the past week.
	SearchRecencyWeek SearchRecency = "week"
	
	// SearchRecencyMonth includes results from the past month.
	SearchRecencyMonth SearchRecency = "month"
	
	// SearchRecencyYear includes results from the past year.
	SearchRecencyYear SearchRecency = "year"
)

// NewSearchParameters creates new search parameters with defaults.
func NewSearchParameters() *SearchParameters {
	return &SearchParameters{
		count:    5,  // Default to 5 results
		recency:  SearchRecencyDefault,
		domains:  []string{},
	}
}

// WithCount sets the number of search results to include.
func (sp *SearchParameters) WithCount(count int32) *SearchParameters {
	sp.count = count
	return sp
}

// WithDomains sets the search domains to include.
func (sp *SearchParameters) WithDomains(domains ...string) *SearchParameters {
	sp.domains = make([]string, len(domains))
	copy(sp.domains, domains)
	return sp
}

// WithRecency sets the search recency filter.
func (sp *SearchParameters) WithRecency(recency SearchRecency) *SearchParameters {
	sp.recency = recency
	return sp
}

// Count returns the number of search results to include.
func (sp *SearchParameters) Count() int32 {
	return sp.count
}

// Domains returns the search domains to include.
func (sp *SearchParameters) Domains() []string {
	if sp.domains == nil {
		return []string{}
	}
	result := make([]string, len(sp.domains))
	copy(result, sp.domains)
	return result
}

// Recency returns the search recency filter.
func (sp *SearchParameters) Recency() SearchRecency {
	return sp.recency
}

// Validate validates the search parameters.
func (sp *SearchParameters) Validate() error {
	if sp.count < 0 || sp.count > 50 {
		return fmt.Errorf("search count must be between 0 and 50, got %d", sp.count)
	}
	
	for _, domain := range sp.domains {
		if domain == "" {
			return fmt.Errorf("search domain cannot be empty")
		}
	}
	
	// Validate recency
	validRecency := map[SearchRecency]bool{
		SearchRecencyDefault: true,
		SearchRecencyDay:     true,
		SearchRecencyWeek:    true,
		SearchRecencyMonth:   true,
		SearchRecencyYear:    true,
	}
	
	if !validRecency[sp.recency] {
		return fmt.Errorf("invalid search recency: %s", sp.recency)
	}
	
	return nil
}

// ToJSON converts search parameters to JSON representation.
func (sp *SearchParameters) ToJSON() map[string]interface{} {
	result := map[string]interface{}{
		"count": sp.count,
	}
	
	if len(sp.domains) > 0 {
		result["domains"] = sp.domains
	}
	
	if sp.recency != SearchRecencyDefault {
		result["recency"] = sp.recency
	}
	
	return result
}

// ReasoningEffort represents the reasoning effort level for chat requests.
type ReasoningEffort string

const (
	// ReasoningEffortLow uses minimal reasoning.
	ReasoningEffortLow ReasoningEffort = "low"
	
	// ReasoningEffortHigh uses maximal reasoning.
	ReasoningEffortHigh ReasoningEffort = "high"
)

// NewReasoningEffortLow creates a low reasoning effort option.
func NewReasoningEffortLow() *ReasoningEffortOption {
	return &ReasoningEffortOption{
		effort: ReasoningEffortLow,
	}
}

// NewReasoningEffortHigh creates a high reasoning effort option.
func NewReasoningEffortHigh() *ReasoningEffortOption {
	return &ReasoningEffortOption{
		effort: ReasoningEffortHigh,
	}
}

// ReasoningEffortOption represents a reasoning effort configuration.
type ReasoningEffortOption struct {
	effort ReasoningEffort
}

// Effort returns the reasoning effort level.
func (reo *ReasoningEffortOption) Effort() ReasoningEffort {
	return reo.effort
}

// Validate validates the reasoning effort option.
func (reo *ReasoningEffortOption) Validate() error {
	validEffort := map[ReasoningEffort]bool{
		ReasoningEffortLow:  true,
		ReasoningEffortHigh: true,
	}
	
	if !validEffort[reo.effort] {
		return fmt.Errorf("invalid reasoning effort: %s", reo.effort)
	}
	
	return nil
}

// ToJSON converts reasoning effort option to JSON representation.
func (reo *ReasoningEffortOption) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"reasoning_effort": reo.effort,
	}
}