package viterbi

// MarkovState is a simple hidden markov state for storing initial state and probabilities prior to the viterbi algorithm
type MarkovState struct {
	InputTokens               []string
	PossibleStates            []string
	InitialStates             [][]string
	InitialStateProbabilities [][]float64
}
