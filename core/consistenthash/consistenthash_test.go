package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	ring := New()
	ring.Add("192.128.1.1:8080")
	ring.Add("192.128.1.2:8080")
	ring.Add("192.128.1.3:8080")
	ring.Add("192.128.1.4:8080")
	ring.Add("192.128.1.5:8080")
	ring.Add("192.128.1.6:8080")
	ring.Add("192.128.1.7:8080")
	ring.Add("192.128.1.8:8080")
	hits := make(map[string]int)
	totalHits := 8000

	for i := 0; i < totalHits; i++ {
		key := "key" + strconv.Itoa(i)
		node := ring.Get(key)
		hits[node]++
	}
	limit := 2.0
	maxRate, minRate := 0.0, 100.0
	for _, count := range hits {
		percentage := (float64(count) / float64(totalHits)) * 100
		if percentage > maxRate {
			maxRate = percentage
		}
		if percentage < minRate {
			minRate = percentage
		}
	}
	if (maxRate - minRate) > limit {
		t.Fatalf("Load balancing average test failed: Maximum percentage: %.2f%%, Minimum percentage: %.2f%%, Difference exceeds %.2f%%", maxRate, minRate, limit)
	}
}
