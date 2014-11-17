package viterbi

// ProbabilityClass represents a tri-gram state
type ProbabilityClass struct {
	P float64 // Probability of this trigram
	W string  // 1st gram of the probability class
	U string  // 2nd gram of this probability class
	V string  // 3rd gram of this probability class
	K int     // index in the path
}
