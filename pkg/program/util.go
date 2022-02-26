package program

import "math/rand"

// RandomSource is a customizable randoom source for replacing for tests or seeding.
type RandomSource interface {
	// Get should return a number in the interval [low, high)
	Get(low int, high int) int
}

// DefaultRandomSource is a RandomSource that uses rand.Intn
type DefaultRandSource struct {
}

// Get implementation for RandomSource.
func (r *DefaultRandSource) Get(low int, high int) int {
	return rand.Intn(high-low) + low
}

// TestingRandSource is an implementation of RandomSource that uses a predefined
// list of values.
type TestingRandSource struct {
	vals []int
}

// Get implementation for RandomSource.
func (r *TestingRandSource) Get(low int, high int) int {
	result := r.vals[0]
	r.vals = r.vals[1:]
	return result
}

// AddMore appends more random values to the internal list in order.
func (r *TestingRandSource) AddMore(vals ...int) {
	r.vals = append(r.vals, vals...)
}

// NewTestRandSource creates a new random source for testing pre-populated
// with the passsed values.
func NewTestRandSource(val ...int) *TestingRandSource {
	return &TestingRandSource{
		vals: val,
	}
}
