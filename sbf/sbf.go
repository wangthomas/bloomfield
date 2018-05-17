package sbf


import (
    "sync"
    "sync/atomic"
    "time"
    "github.com/wangthomas/bloomfield/plainFilter"
)

// Scalable Bloom Filters
// Concept is described in the link below
// http://gsd.di.uminho.pt/members/cbm/ps/dbloom.pdf

var r float64 = 0.9
var scale_size uint64 = 4


type SBF struct {
    capacity            uint64
    hits                uint64
    misses              uint64
    checks              uint64
    keys                uint64
    probability         float64
    createDate          time.Time
    mutex               sync.RWMutex
    plainFilters        []*plainFilter.PlainFilter
}

// New returns a new SBF
func NewSBF(capacity uint64, probability float64) *SBF {
    var plainFilters []*plainFilter.PlainFilter

    p := probability * (1 - r)

    pf := plainFilter.NewPlainFilter(capacity, p)
    plainFilters = append(plainFilters, pf)

    return &SBF{
        createDate:     time.Now(),
        probability:    probability,
        capacity:       capacity,
        hits:           0,
        misses:         0,
        checks:         0,
        keys:           0,
        plainFilters:   plainFilters,
    }
}

// NewDefault returns a new Filter with default settings
// 100,000 capacity and 0.01 false positive probability
func NewSBFDefault() *SBF {
    return NewSBF(65536, 0.0001)
}


// Retues false if the key is in the SBF already.
// Otherwise add the key and return false.
func (t *SBF) Add(hashes []uint64) bool {
    // Check if the key is in SBF already
    t.mutex.RLock()
    for _, pf := range t.plainFilters {
        if pf.Has(hashes) {
            t.mutex.RUnlock()
            // Has the key already.
            return true
        }
    }
    t.mutex.RUnlock()

    // Add the key to SBF
    t.mutex.Lock()
    defer t.mutex.Unlock()
    pf := t.plainFilters[len(t.plainFilters)-1]

    if t.keys == t.capacity {
        // SBF is full. Expand it by attaching another plainFilter
        pf := plainFilter.NewPlainFilter(scale_size*pf.Capacity, r*pf.Probability)
        t.plainFilters = append(t.plainFilters, pf)
        atomic.AddUint64(&t.capacity, pf.Capacity)
    }

    // In most cases added is false. Since we checked the key is not in the filter in the
    // top half of this function. But there is a tiny chance there is a context switch happens
    // between the RWLock and we could add the same key twice. So double check added here.
    added := pf.Add(hashes)
    if !added {
        atomic.AddUint64(&t.keys, 1)
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


// Probability returns the false positive probability of the bloom filter
func (t *SBF) Probability() float64 {
    return t.probability
}

func (t *SBF) CreateDate() time.Time {
    return t.createDate
}

// Capacity returns the capacity of the bloom filter
func (t *SBF) Capacity() uint64 {
    return atomic.LoadUint64(&t.capacity)
}

// Checks returns the key check count on the bloom filter
func (t *SBF) Checks() uint64 {
    return atomic.LoadUint64(&t.checks)
}


// Keys returns the number of keys set in the bloom filter
func (t *SBF) Keys() uint64 {
    return atomic.LoadUint64(&t.keys)
}

// Hits returns the key check hits on the bloom filter
func (t *SBF) Hits() uint64 {
    return atomic.LoadUint64(&t.hits)
}

// Misses returns the key check misses on the bloom filter
func (t *SBF) Misses() uint64 {
    return atomic.LoadUint64(&t.misses)
}

