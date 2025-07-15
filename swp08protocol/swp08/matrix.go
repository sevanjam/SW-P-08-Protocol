package swp08

// MatrixQuery defines an interface for SW-P-08 to query the source-to-destination mapping

type MatrixQuery interface {
	GetSourceForDestination(matrix, level, destination int) int
	GetMatrixSize(matrix, level int) (int, int)

	// Optional
	SetSourceForDestination(matrix, level, destination, source int) bool
}

var matrix MatrixQuery

// Host application should call this to register its matrix implementation
func RegisterMatrixQuery(m MatrixQuery) {
	matrix = m
}
