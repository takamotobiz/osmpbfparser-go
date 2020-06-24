package osmpbfparser

// New ...
func New(
	args Args,
) PBFParser {
	return &pbfParser{Args: args}
}
