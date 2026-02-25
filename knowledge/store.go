package knowledge

import (
	"fmt"
	"sync"
)

// Store is the persistent knowledge base that feeds all nodes in Illygen.
// It is namespaced by domain so nodes only consult relevant knowledge.
// The store is safe for concurrent access.
type Store struct {
	mu    sync.RWMutex
	units map[string]*Unit // keyed by unit ID
}

// NewStore creates an empty in-memory KnowledgeStore.
func NewStore() *Store {
	return &Store{
		units: make(map[string]*Unit),
	}
}

// Add inserts a KnowledgeUnit into the store.
// Returns an error if a unit with the same ID already exists.
func (s *Store) Add(unit *Unit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.units[unit.ID]; exists {
		return fmt.Errorf("illygen/knowledge: unit %q already exists", unit.ID)
	}
	s.units[unit.ID] = unit
	return nil
}

// Set inserts or replaces a KnowledgeUnit by ID.
func (s *Store) Set(unit *Unit) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.units[unit.ID] = unit
}

// Get retrieves a KnowledgeUnit by ID.
func (s *Store) Get(id string) (*Unit, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.units[id]
	return u, ok
}

// Domain returns all KnowledgeUnits belonging to a given domain,
// ordered by weight descending (most trusted first).
func (s *Store) Domain(domain string) []*Unit {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Unit
	for _, u := range s.units {
		if u.Domain == domain {
			result = append(result, u)
		}
	}
	sortByWeight(result)
	return result
}

// Remove deletes a KnowledgeUnit from the store by ID.
func (s *Store) Remove(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.units, id)
}

// Size returns the total number of knowledge units in the store.
func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.units)
}

// sortByWeight sorts units by weight descending (simple insertion sort — stores are small).
func sortByWeight(units []*Unit) {
	for i := 1; i < len(units); i++ {
		for j := i; j > 0 && units[j].Weight > units[j-1].Weight; j-- {
			units[j], units[j-1] = units[j-1], units[j]
		}
	}
}
