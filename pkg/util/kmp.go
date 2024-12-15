/*
Knuth–Morris–Pratt algorithm
See: https://go.dev/play/p/chYGT69vBc
*/

package util

// KmpSearch returns the index (0 based) of the start of the string 'find' in 'corpus', or returns -1 on failure.
func KmpSearch(corpus, find string) int {
	m, i := 0, 0
	table := make([]int, len(corpus))
	kmpTable(find, table)
	for m+i < len(corpus) {
		if find[i] == corpus[m+i] {
			if i == len(find)-1 {
				return m
			}
			i++
		} else {
			if table[i] > -1 {
				i = table[i]
				m = m + i - table[i] //
			} else {
				i = 0
				m++
			}
		}
	}
	return -1
}

// kmpTable populates the partial match table 'table' for the string 'find'.
func kmpTable(find string, table []int) {
	pos, cnd := 2, 0
	table[0], table[1] = -1, 0
	for pos < len(find) {
		if find[pos-1] == find[cnd] {
			cnd++
			table[pos] = cnd
			pos++
		} else if cnd > 0 {
			cnd = table[cnd]
		} else {
			table[pos] = 0
			pos++
		}
	}
}
