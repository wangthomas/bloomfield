package sbf


import (
    "sync"
    "sync/atomic"
    "time"
    "github.com/wangthomas/bloomfield/bloom/plainFilter"
)

// Scalable Bloom Filters
// http://gsd.di.uminho.pt/members/cbm/ps/dbloom.pdf

var r float64 = 0.9
var scale_size uint64 = 4


type SBF struct {
    capacity     	uint64
    hits         	uint64
    misses       	uint64
    checks       	uint64
    num_keys     	uint64
    probability  	float64
    createDate   	time.Time
    mutex        	sync.RWMutex
    plainFilters 	[]*plainFilter.PlainFilter

}

// New returns a new SBF
func NewSBF(capacity uint64, probability float64) *SBF {
    var plainFilters []*plainFilter.PlainFilter

    p := probability * (1 - r)

    pf := plainFilter.NewPlainFilter(capacity, p)
    plainFilters = append(plainFilters, pf)

    return &SBF{
        createDate:   time.Now(),
        capacity:     capacity,
        hits:         0,
        misses:       0,
        checks:       0,
        probability:  probability,
        num_keys:         0,
        plainFilters: plainFilters,
    }
}

// NewDefault returns a new Filter with default settings
// 100,000 capacity and 0.01 false positive probability
func NewSBFDefault(name string) *SBF {
    return NewSBF(262144, 0.0001)
}


func (t *SBF) Add(hashes []uint64) bool {
    if t.Has(hashes) {
        return false
    }

    t.mutex.Lock()
    defer t.mutex.Unlock()
    pf := t.plainFilters[len(t.plainFilters)-1]

    if t.num_keys == t.capacity {
        pf := plainFilter.NewPlainFilter(scale_size*pf.Capacity, r*pf.Probability)
        t.plainFilters = append(t.plainFilters, pf)
        t.capacity = t.capacity + pf.Capacity
    }

    added := pf.Add(hashes)
    if added {
    	t.num_keys++
    }
    
    return added
}


// Has checks if a key is in a bloom
// returns true if the key was in the filter, false otherwise
func (t *SBF) Has(hashes []uint64) bool {
    t.mutex.RLock()

    has := false

    for _, pf := range t.plainFilters {
        if pf.Has(hashes) {
            has = true
            break
        }
    }

    t.mutex.RUnlock()

    atomic.AddUint64(&t.checks, 1)
    if has {
        atomic.AddUint64(&t.hits, 1)
    } else {
        atomic.AddUint64(&t.misses, 1)
    }
    return has
}

