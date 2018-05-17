package filterManager

import (
    "testing"
    "hash/fnv"
    "github.com/OneOfOne/xxhash"
    . "gopkg.in/check.v1"
)


func Test(t *testing.T) { TestingT(t) }

type BloomSuite struct {
    filterManager       *FilterManager
    keys                []string
}

var _ = Suite(&BloomSuite{})

func (s *BloomSuite) SetUpTest(c *C) {
    s.filterManager = NewFilterManager()
}


func (s *BloomSuite) TestAdd(c *C) {
    s.filterManager.Create("testFilter")
    c.Assert(Add("testFilter", "num1", s.filterManager), Equals, false)
    c.Assert(Add("testFilter", "num2", s.filterManager), Equals, false)

    c.Assert(Add("testFilter", "num1", s.filterManager), Equals, true)
    c.Assert(Add("testFilter", "num2", s.filterManager), Equals, true)

    c.Assert(Has("testFilter", "num1", s.filterManager), Equals, true)
    c.Assert(Has("testFilter", "num2", s.filterManager), Equals, true)
    c.Assert(Has("testFilter", "num3", s.filterManager), Equals, false)
    c.Assert(Has("testFilter1", "num1", s.filterManager), Equals, false)
}


func Add(filterName string, key string, fm *FilterManager) bool {
    h1 := fnv.New64a()
    h1.Write([]byte(key))
    hash1 := h1.Sum64()

    h2 := xxhash.New64()
    h2.Write([]byte(key))
    hash2 := h2.Sum64()

    return fm.Add(filterName, []uint64{hash1, hash2})
}


func Has(filterName string, key string, fm *FilterManager) bool {
    h1 := fnv.New64a()
    h1.Write([]byte(key))
    hash1 := h1.Sum64()

    h2 := xxhash.New64()
    h2.Write([]byte(key))
    hash2 := h2.Sum64()

    return fm.Has(filterName, []uint64{hash1, hash2})
}

