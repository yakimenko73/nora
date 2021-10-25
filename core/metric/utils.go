package metric

import (
	"sort"
)

func sortedInsert(entries []*ChartEntry, entry *ChartEntry) []*ChartEntry {
	idx := sort.Search(len(entries), func(i int) bool { return entries[i].Timestamp.After(entry.Timestamp) })
	entries = append(entries, &ChartEntry{})

	copy(entries[idx+1:], entries[idx:])
	entries[idx] = entry

	return entries
}

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}
