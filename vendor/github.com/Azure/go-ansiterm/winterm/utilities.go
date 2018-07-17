

package winterm



func addInRange(n int16, increment int16, min int16, max int16) int16 {
	return ensureInRange(n+increment, min, max)
}
