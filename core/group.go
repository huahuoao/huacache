package huacache

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/huahuoao/huacache/core/lru"
)

type GroupStatus struct {
	name     string
	size     int64
	used     int64
	keyCount int
}
type Group struct {
	name      string
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup creates a new instance of Group
func NewGroup(name string, cacheBytes int64) (*Group, error) {
	mu.Lock()
	defer mu.Unlock()
	if name == "" {
		return nil, fmt.Errorf("group name can't be empty")
	}
	_, ok := groups[name]
	if ok {
		return nil, fmt.Errorf("group %s already exists", name)
	}

	g := &Group{
		name: name,
		mainCache: cache{
			cacheBytes: cacheBytes,
			lru:        lru.New(cacheBytes, nil), // Initialize lru here
		},
	}
	groups[name] = g
	return g, nil
}

func DelGroup(name string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := groups[name]; !exists {
		return fmt.Errorf("group %s does not exist", name)
	}
	delete(groups, name)
	return nil
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return ByteView{}, fmt.Errorf("key not found in cache")

}

func (g *Group) AddOrUpdate(key string, value ByteView) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}
	return g.mainCache.add(key, value)
}

func (g *Group) Delete(key string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}
	err := g.mainCache.delete(key)
	return err
}
func (g *Group) Keys() ([]string, error) {

	return g.mainCache.lru.Keys()
}

//Group Methods

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) (*Group, error) {
	mu.RLock()
	defer mu.RUnlock()
	g, ok := groups[name]
	if !ok {
		return nil, errors.New("group not found")
	}
	return g, nil
}

// ListGroups returns a slice of all group names.
func ListGroups() ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	if groups == nil {
		return nil, fmt.Errorf("groups is not initialized")
	}

	names := make([]string, 0, len(groups))
	for name := range groups {
		names = append(names, name)
	}
	return names, nil
}

func (g *Group) GetStatus() *GroupStatus {
	mu.RLock()
	defer mu.RUnlock()
	maxbytes, nbytes, keyCount := g.mainCache.lru.GetMemoryUsedSituation()
	gs := &GroupStatus{
		name:     g.name,
		size:     maxbytes,
		used:     nbytes,
		keyCount: keyCount,
	}
	return gs
}

func (gs *GroupStatus) ToString() string {
	// 将字节转换为MB
	sizeMB := float64(gs.size) / (1024 * 1024)
	usedMB := float64(gs.used) / (1024 * 1024)

	// 计算使用率
	usageRate := (float64(gs.used) / float64(gs.size)) * 100
	return fmt.Sprintf("Group Status:\nName: %s\nSize: %.2f MB\nUsed: %.2f MB\nUsage Rate: %.2f%%\nKey Count: %d\n",
		gs.name, sizeMB, usedMB, usageRate, gs.keyCount)
}
