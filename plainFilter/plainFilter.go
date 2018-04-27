package plainFilter

import (
    "math"
)

type PlainFilter struct {
    Probability            float64
    num_hash               uint64
    Capacity               uint64
    num_bits_inslice       uint64
    bitmap                 []uint8
}

// Alloc a new plain filter based on the max number of keys and false positive rate.
// The number of hash values(number of slices also) and number of bits needed are
// calculated based on the paper below.
//
// http://haslab.uminho.pt/cbm/files/dbloom.pdf 
func NewPlainFilter(capacity uint64, probability float64) *PlainFilter {
    bit_size := math.Ceil(-1 * float64(capacity) * math.Log(probability) / (math.Log(2) * math.Log(2)))
    num_hash := uint64(math.Ceil(math.Log2(1 / probability)))
    num_uint8 := uint64(math.Ceil(bit_size / 8))
    num_bits_inslice := (num_uint8 * 8) / num_hash
    bitmap := make([]uint8, num_uint8)

    return &PlainFilter {
        bitmap:                   bitmap,
        num_hash:                 num_hash,
        Capacity:                 capacity,
        num_bits_inslice:         num_bits_inslice,
        Probability:              probability,
    }
}


// Add a key to the filter. The key should be converted into two independent hash
// values on the client side already.
// Return true if it's a new key
func (pf *PlainFilter) Add(hashes []uint64) bool {
    is_new_key := false
    for i := uint64(0); i < pf.num_hash; i++ {
        new_hash := getHash(i, hashes)
        uint8_index, shift_index := pf.getIndexShift(i, new_hash)
        is_set := (pf.bitmap[uint8_index] & (0x1 << shift_index)) != 0
        
        if !is_set {
            is_new_key = true
            pf.bitmap[uint8_index] |= 0x1 << shift_index
        }
    }
    return is_new_key
}


// Check if a key is in the filter. The key should be converted into two independent hash
// values on the client side already.
func (pf *PlainFilter) Has(hashes []uint64) bool {
    for i := uint64(0); i < pf.num_hash; i++ {
        new_hash := getHash(i, hashes)
        uint8_index, shift_index := pf.getIndexShift(i, new_hash)

        if (pf.bitmap[uint8_index] & (0x1 << shift_index)) == 0  {
            return false
        }
    }
    return true
}


// Calculate where the bit should be set/checked. 
func (pf *PlainFilter) getIndexShift(index uint64, hash uint64) (uint64, uint8) {
    bit_index := index * pf.num_bits_inslice + (hash % pf.num_bits_inslice)
    uint8_index := bit_index / 8
    shift_index := bit_index % 8
    return uint8_index, uint8(shift_index)

}


// Generate a new hash value based on two fixed hash values and an index
func getHash(index uint64, hashes []uint64) uint64 {
    new_hash := hashes[0] + index * hashes[1]
    return new_hash
}


