package filterManager

import (
    "sync"
    "github.com/wangthomas/bloomfield/sbf"
    pb "github.com/wangthomas/bloomfield/interfaces/gRPC/bloomfieldpb"
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


func (t *FilterManager) Add(name string, hashes []*pb.Hashes) []bool {
    var res []bool

    // Create the filter if it does not exist
    t.mutex.Lock()

    if _, exists := t.filters[name]; !exists {
        t.filters[name] = sbf.NewSBFDefault()
    }
    filter := t.filters[name]
    t.mutex.Unlock()

    //return []bool{filter.Add([]uint64{hashes[0].Hash1, hashes[0].Hash2})}

    for _, hash := range hashes {
        res = append(res, filter.Add([]uint64{hash.Hash1, hash.Hash2}))
    }

    return res
}


func (t *FilterManager) Has(name string, hashes []*pb.Hashes) []bool {
    var res []bool

    t.mutex.RLock()
    if _, exists := t.filters[name]; !exists {
        t.mutex.RUnlock()
        for range hashes {
            res = append(res, false)
        }
        return res
    }
    filter := t.filters[name]
    t.mutex.RUnlock()

    //return []bool{filter.Has([]uint64{hashes[0].Hash1, hashes[0].Hash2})}

    for _, hash := range hashes {
        res = append(res, filter.Has([]uint64{hash.Hash1, hash.Hash2}))
    }

    return res
}


func (t *FilterManager) Drop(name string) {
    t.mutex.Lock()
    defer t.mutex.Unlock()
    if _, exists := t.filters[name]; !exists {
        return
    }
    delete(t.filters, name)
}

