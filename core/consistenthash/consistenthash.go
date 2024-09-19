package consistenthash

import (
	"crypto/md5"
	"sort"
	"strconv"

	huacache "github.com/huahuoao/huacache/core"
)

// computeMD5 computes the MD5 hash of the given string.
func computeMD5(s string) [16]byte {
	return md5.Sum([]byte(s))
}

// hash extracts a specific 32-bit integer from the digest (Ketama feature).
func hash(digest *[16]byte, h int) int64 {
	k := ((int64((*digest)[3+h*4]) & 0xFF) << 24) |
		((int64((*digest)[2+h*4]) & 0xFF) << 16) |
		((int64((*digest)[1+h*4]) & 0xFF) << 8) |
		(int64((*digest)[h*4]) & 0xFF)
	return k
}

// Map represents the structure of a consistent hash ring.
type Map struct {
	replicas int              // Number of virtual nodes per physical node
	keys     []int64          // Sorted hash values
	hashMap  map[int64]string // Mapping from hash values to physical node names
}

// New creates a new hash ring.
func New() *Map {
	m := &Map{
		replicas: huacache.CONSISTENTHASH_VIRTUAL_NODE_NUM, // Number of virtual nodes
		hashMap:  make(map[int64]string),
	}
	return m
}

// Add adds new physical nodes to the hash ring.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			virtualNodeKey := key + strconv.Itoa(i)
			digest := computeMD5(virtualNodeKey)
			for j := 0; j < 4; j++ {
				hash := hash(&digest, j)
				m.keys = append(m.keys, hash)
				m.hashMap[hash] = key
			}
		}
	}
	sort.Slice(m.keys, func(i, j int) bool {
		return m.keys[i] < m.keys[j]
	})
}

// Get retrieves the closest physical node for the given key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	digest := computeMD5(key)
	hash := hash(&digest, 0)
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}
