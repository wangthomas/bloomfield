package plainFilter

import (
	"testing"
	"fmt"
)

// type BloomSuite struct {
// 	filter PlainFilter
// }

// func Test(t *testing.T) {
// 	TestingT(t)
// }

func TestNewPlainFilter(t *testing.T) {
	cap := uint64(100)
	filter := NewPlainFilter(cap, 0.001)
	fmt.Println("num_hash is ", filter.num_hash)
	fmt.Println("bitmap size is ", len(filter.bitmap))
	fmt.Println("num_bits_inslice is ", filter.num_bits_inslice)
	filter.Add([]uint64{12446755073609551590, 10446755473632211101})
	fmt.Println("num_bitmap is ", filter.bitmap)
	isSet := filter.Has([]uint64{12446755073609551598, 10446755473632211100}) 
	fmt.Println("isSet is ", isSet)



	// c.Assert(filter.Capacity(), Equals, cap)
	// c.Assert(filter.Checks(), Equals, int64(0))
	// c.Assert(filter.Hits(), Equals, int64(0))
	// c.Assert(filter.Misses(), Equals, int64(0))
	// c.Assert(filter.Keys(), Equals, int64(0))
}

