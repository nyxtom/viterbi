package viterbi

func MarkovValue() map[string]float64 {
	return make(map[string]float64)
}
func MarkovValueMap() map[string]map[string]float64 {
	return make(map[string]map[string]float64)
}
func MarkovValues() []map[string]map[string]float64 {
	return make([]map[string]map[string]float64, 0)
}
func MarkovPath() map[string]string {
	return make(map[string]string)
}
func MarkovPathMap() map[string]map[string]string {
	return make(map[string]map[string]string)
}
func MarkovPaths() []map[string]map[string]string {
	return make([]map[string]map[string]string, 0)
}
