package plainFilter

import (
	"testing"
	"hash/fnv"
	"github.com/OneOfOne/xxhash"
	. "gopkg.in/check.v1"
)


func Test(t *testing.T) { TestingT(t) }

type BloomSuite struct {
	filter *PlainFilter
}

var _ = Suite(&BloomSuite{})

func (s *BloomSuite) SetUpTest(c *C) {
	cap := uint64(10000)
	s.filter = NewPlainFilter(cap, 0.0001)
	Add("test", s.filter)
	Add("test_company_0x12446755073609551590", s.filter)
}


func (s *BloomSuite) TestNew(c *C) {

	c.Assert(s.filter.probability, Equals, 0.0001)
	c.Assert(s.filter.capacity, Equals, uint64(10000))
	c.Assert(s.filter.num_hash, Equals, uint64(14))
	c.Assert(s.filter.num_bits_inslice, Equals, uint64(13693))
	c.Assert(len(s.filter.bitmap), Equals, 23963)
}


func (s *BloomSuite) TestHas(c *C) {	
	c.Assert(Has("test", s.filter), Equals, true)
	c.Assert(Has("test_company_0x12446755073609551590", s.filter), Equals, true)
}


func (s *BloomSuite) TestNotHas(c *C) {
	c.Assert(Has("new_test", s.filter), Equals, false)
	c.Assert(Has("test_company_0x12446755073609551591", s.filter), Equals, false)
}


func Add(s string, filter *PlainFilter) {
	h1 := fnv.New64a()
	h1.Write([]byte(s))
	hash1 := h1.Sum64()

	h2 := xxhash.New64()
	h2.Write([]byte(s))
	hash2 := h2.Sum64()

	filter.Add([]uint64{hash1, hash2})
}

func Has(s string, filter *PlainFilter) bool {
	h1 := fnv.New64a()
	h1.Write([]byte(s))
	hash1 := h1.Sum64()

	h2 := xxhash.New64()
	h2.Write([]byte(s))
	hash2 := h2.Sum64()

	return filter.Has([]uint64{hash1, hash2})
}



