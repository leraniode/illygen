// Package knowledge manages the knowledge store that feeds intelligence into Illygen nodes.
package knowledge

import "time"

// Unit is the atomic piece of knowledge in Illygen.
// It is lightweight — similar in concept to a tensor in AI, but structured around
// human-readable facts rather than numerical matrices.
//
// Training produces Units with high initial Weight.
// Exploring refines existing Units over time.
type Unit struct {
	// ID uniquely identifies this knowledge unit.
	ID string

	// Domain scopes this unit to a reasoning area (e.g. "programming", "medicine").
	// Nodes consult only units within their relevant domain.
	Domain string

	// Facts holds the actual knowledge — key/value structured data.
	Facts map[string]any

	// Weight represents how trusted or relevant this unit is (0.0 to 1.0).
	// Higher weight = more influence when nodes consult this unit.
	Weight float64

	// Refined marks whether this unit has been processed and improved by learning.
	Refined bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUnit creates a new KnowledgeUnit with sensible defaults.
func NewUnit(id, domain string, facts map[string]any) *Unit {
	now := time.Now()
	return &Unit{
		ID:        id,
		Domain:    domain,
		Facts:     facts,
		Weight:    1.0,
		Refined:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Fact returns a single fact value by key. Returns nil if not found.
func (u *Unit) Fact(key string) any {
	return u.Facts[key]
}

// Refine marks this unit as refined and updates its weight.
// Called by the learning logic during training or exploring.
func (u *Unit) Refine(newWeight float64) {
	u.Weight = clamp(newWeight, 0.0, 1.0)
	u.Refined = true
	u.UpdatedAt = time.Now()
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
