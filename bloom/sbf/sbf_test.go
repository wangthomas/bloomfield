package sbf

import (
    "testing"
    "hash/fnv"
    "math/rand"
    "github.com/OneOfOne/xxhash"
    . "gopkg.in/check.v1"
)

const (
    letterBytes        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$^&*()_+-="
    capacity         = 100
    probability        = 0.0001
    key_length        = 50
    num_keys        = 1000
) 


func Test(t *testing.T) { TestingT(t) }

type BloomSuite struct {
    filter      *SBF
    keys        []string
}

var _ = Suite(&BloomSuite{})

func (s *BloomSuite) SetUpTest(c *C) {
    cap := uint64(capacity)
    s.filter = NewSBF(cap, probability)
    // Init s.keys with random strings
    for i := 0; i < num_keys; i++ {
        b := make([]byte, key_length)
        for j := range b {
            b[j] = letterBytes[rand.Intn(len(letterBytes))]
        }
        s.keys = append(s.keys, string(b))
    }
}


func (s *BloomSuite) TestNew(c *C) {
    c.Assert(s.filter.probability, Equals, probability)
    c.Assert(s.filter.capacity, Equals, uint64(capacity))
}


func (s *BloomSuite) TestAdd(c *C) {
    for _, key := range s.keys {
        c.Assert(Add(key, s.filter), Equals, true)
    }

    for _, key := range s.keys {
        c.Assert(Add(key, s.filter), Equals, false)
    }
}


func (s *BloomSuite) TestHas(c *C) {
    for _, key := range s.keys {
        Add(key, s.filter)
    }

    for _, key := range s.keys {
        c.Assert(Has(key, s.filter), Equals, true)
    }
}


// % is not in the letterBytes. With reasonably small false positive rate we should not fail this test case.
// But take it as a guarantee
func (s *BloomSuite) TestNotHas(c *C) {

    for _, key := range s.keys {
        Add(key, s.filter)
    }

    c.Assert(Has("2CTM26WgMtQTcUyjbRuucAM6Th2j4nHVYbzsAy1uVBOlFAuEs%", s.filter), Equals, false)
    c.Assert(Has("6sFgnz0KPRy0dDer7hfLUhE6QJNVUVXZfY%fvfQ9hSu29MpDuU", s.filter), Equals, false)
    c.Assert(Has("juXWihBRrliVwXkB9Ak9%nManCN72ia50paT7fV1fkcx9EcbP5", s.filter), Equals, false)
    c.Assert(Has("7DdxUYlDyhGxm%tBX1G4ELk0RekTboc2PKo3QGLTYEXaDwYoXg", s.filter), Equals, false)
    c.Assert(Has("%eaVQoClbdsC07SjG0j991KPWaPHOSw8FgWJfp7PEjFjcZA3Bt", s.filter), Equals, false)
}


func (s *BloomSuite) BenchmarkAdd(c *C) {
    for n := 0; n < c.N; n++ {
        Add(s.keys[0], s.filter)
    }
}


func (s *BloomSuite) BenchmarkHas(c *C) {
    for n := 0; n < c.N; n++ {
        Has(s.keys[0], s.filter)
    }
}


func Add(s string, filter *SBF) bool {
    h1 := fnv.New64a()
    h1.Write([]byte(s))
    hash1 := h1.Sum64()

    h2 := xxhash.New64()
    h2.Write([]byte(s))
    hash2 := h2.Sum64()

    return filter.Add([]uint64{hash1, hash2})
}

func Has(s string, filter *SBF) bool {
    h1 := fnv.New64a()
    h1.Write([]byte(s))
    hash1 := h1.Sum64()

    h2 := xxhash.New64()
    h2.Write([]byte(s))
    hash2 := h2.Sum64()

    return filter.Has([]uint64{hash1, hash2})
}
