package swp08

// MockMatrix is a simple in-memory mock implementing MatrixQuery for testing
type MockMatrix struct {
	sources         [][]int // sources[destination] = source number
	numSources      int
	numDestinations int
}

// NewMockMatrix creates a mock matrix with the given size and a simple 1:1 mapping
func NewMockMatrix(numSources, numDestinations int) *MockMatrix {
	m := &MockMatrix{
		sources:         make([][]int, numDestinations),
		numSources:      numSources,
		numDestinations: numDestinations,
	}
	for dest := 0; dest < numDestinations; dest++ {
		m.sources[dest] = []int{dest % numSources} // simple mapping source = dest % numSources
	}
	return m
}

// GetSourceForDestination returns the first source connected to the given destination
func (m *MockMatrix) GetSourceForDestination(matrix, level, destination int) int {
	if destination < 0 || destination >= m.numDestinations {
		return 0 // or some default invalid source
	}
	// Return the first source for simplicity
	return m.sources[destination][0]
}

// GetMatrixSize returns the configured number of sources and destinations
func (m *MockMatrix) GetMatrixSize(matrix, level int) (int, int) {
	return m.numSources, m.numDestinations
}

// UseExtendedTallyDump optionally forces extended mode for testing
func (m *MockMatrix) UseExtendedTallyDump(matrix, level int) bool {
	// For testing, let's just return false by default
	return true
}
func (m *MockMatrix) SetSourceForDestination(matrix, level, dest, source int) bool {
	if dest < 0 || dest >= m.numDestinations || source < 0 || source >= m.numSources {
		return false
	}
	m.sources[dest][0] = source
	return true
}
