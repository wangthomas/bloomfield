package filterManager

import (
    "sync"
    "github.com/wangthomas/bloomfield/sbf"
)

type FilterManager struct {
    filters map[string]*sbf.SBF
    mutex   sync.RWMutex
}

func NewFilterManager() *FilterManager {
    filterMap := make(map[string]*sbf.SBF)
    return &FilterManager{
        filters: filterMap,
    }
}


func (t *FilterManager) Create(name string) {

    t.mutex.Lock()
    defer t.mutex.Unlock()
    if _, exists := t.filters[name]; !exists {
        t.filters[name] = sbf.NewSBFDefault()
    }
}


func (t *FilterManager) Add(name string, hashes []uint64) bool {

    // Create the filter if it does not exist
    t.mutex.Lock()

    if _, exists := t.filters[name]; !exists {
        t.filters[name] = sbf.NewSBFDefault()
    }
    filter := t.filters[name]
    t.mutex.Unlock()

    return filter.Add(hashes)
}


func (t *FilterManager) Has(name string, hashes []uint64) bool {

    t.mutex.RLock()
    if _, exists := t.filters[name]; !exists {
        t.mutex.RUnlock()
        return false
    }
    filter := t.filters[name]
    t.mutex.RUnlock()
    return filter.Has(hashes)
}


func (t *FilterManager) Drop(name string) {
    t.mutex.Lock()
    defer t.mutex.Unlock()
    if _, exists := t.filters[name]; !exists {
        return
    }
    delete(t.filters, name)
}

