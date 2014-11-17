package viterbi

import "sort"
import "math"

// Viterbi represents the hidden markov model viterbi path finding algorithm using the trigram model
func Viterbi(state *MarkovState, transitionProbFn func(w, u, v string) float64) []string {
	if state == nil || state.InitialStates == nil {
		return []string{}
	}

	pi := MarkovValues()
	path := MarkovPaths()

	// Initialize the base state where pi(0, *, *) = 1 and pi(0, u, v) = 0 for all (u, v)
	initial := MarkovValueMap()
	initial["*"] = MarkovValue()
	initial["*"]["*"] = 1.0
	pi = append(pi, initial)

	for _, stateU := range state.PossibleStates {
		if stateU == "*" {
			continue
		}

		pi[0][stateU] = MarkovValue()
		for _, stateV := range state.PossibleStates {
			if stateV == "*" {
				continue
			}

			pi[0][stateU][stateV] = 1.0
		}
	}

	for k := 1; k <= len(state.InputTokens); k++ {
		tempPath := MarkovPathMap()
		tempPi := MarkovValueMap()
		prevPi := pi[k-1]

		statesW := getPossibleStates(k-2, state.InitialStates, "*")
		statesU := getPossibleStates(k-1, state.InitialStates, "*")
		statesV := getPossibleStates(k, state.InitialStates, "*")

		// Get the emission probability for the token at k - 1
		emissionsV := state.InitialStateProbabilities[k-1]

		// Iterate over the states U (states at k - 1)
		for _, u := range statesU {
			tempPi[u] = MarkovValue()
			tempPath[u] = MarkovPath()
			var argMax *ProbabilityClass
			for vIndex, v := range statesV {
				// Iterate over the given emission for token k with all states at k-2
				emissionV := emissionsV[vIndex]
				for _, w := range statesW {
					// Using the transition probability between w->u->v
					// calculate the maximum probability path
					transitionQ := transitionProbFn(w, u, v)
					logEmissionV := math.Log(emissionV)
					if logEmissionV == float64(0) {
						logEmissionV = float64(-1)
					}

					// given the prior probability, transition probability and logarithmic emission of v,
					// take the arg max by comparing it with our stored value
					prevPiProb := prevPi[w][u]
					prob := prevPiProb * transitionQ * logEmissionV
					if argMax == nil || prob < argMax.P {
						argMax = new(ProbabilityClass)
						argMax.P = prob
						argMax.W = w
						argMax.U = u
						argMax.V = v
						argMax.K = k
					}
				}

				if argMax == nil {
					continue
				}

				tempPi[u][v] = argMax.P
				tempPath[u][v] = argMax.W
			}
		}

		pi, path = append(pi, tempPi), append(path, tempPath)
	}

	// Iterate of the final stop-tri-gram to locate the most probable ending path
	lastPi := pi[len(pi)-1]
	var argMax *ProbabilityClass
	for u, uMap := range lastPi {
		for v, prob := range uMap {
			distProb := transitionProbFn(u, v, "STOP")
			if distProb == float64(0) {
				distProb = math.Log(0.9999999999)
			}

			prob = prob * distProb
			if argMax == nil || argMax.P < prob {
				argMax = new(ProbabilityClass)
				argMax.P = prob
				argMax.U = u
				argMax.V = v
			}
		}
	}

	// Using the most probable ending, store the best sequence in reverse and invert the backpointers
	var y []string
	if argMax == nil {
		return y
	}

	y = []string{argMax.V, argMax.U}

	// reverse the path appropriately
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	// in each case, the following best tag is the one listed under the backpointer
	// for the current best known tag
	currentBestY := argMax.V  // y (last tag in sequence)
	currentBestY1 := argMax.U // y - 1 (2nd to last tag in sequence)
	for _, p := range path {
		tempCurrentBest := p[currentBestY1][currentBestY]
		currentBestY = currentBestY1
		currentBestY1 = tempCurrentBest
		if tempCurrentBest != "*" {
			y = append(y, tempCurrentBest)
		}
	}

	sort.Reverse(sort.StringSlice(y))
	return y
}

// getPossibleStates returns the possible states at the given k with the otherwise default state
func getPossibleStates(k int, states [][]string, defaultState string) []string {
	if k < 1 {
		return []string{defaultState}
	}

	return states[k-1]
}
