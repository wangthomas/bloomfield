package filterManager

import (
    "testing"
    "hash/fnv"
    "github.com/OneOfOne/xxhash"
    . "gopkg.in/check.v1"
    pb "github.com/wangthomas/bloomfield/interfaces/gRPC/bloomfieldpb"
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

    res := s.filterManager.Add("testFilter", getHashes([]string{"num1", "num2"}))

    // Keys were not in filter before. Should return false
    c.Assert(res[0], Equals, false)
    c.Assert(res[1], Equals, false)

    res = s.filterManager.Add("testFilter", getHashes([]string{"num1", "num2"}))

    // Keys were added in filter above. Should return true
    c.Assert(res[0], Equals, true)
    c.Assert(res[1], Equals, true)

    res = s.filterManager.Has("testFilter", getHashes([]string{"num1", "num2", "num3"}))

    // num3 is not in filter. Should return false
    c.Assert(res[0], Equals, true)
    c.Assert(res[1], Equals, true)
    c.Assert(res[2], Equals, false)

}


func getHashes(keys []string) []*pb.Hashes {
    var hashes []*pb.Hashes
    for _, key := range keys {
        h1 := fnv.New64a()
        h1.Write([]byte(key))
        hash1 := h1.Sum64()

        h2 := xxhash.New64()
        h2.Write([]byte(key))
        hash2 := h2.Sum64()

        hashes = append(hashes, &pb.Hashes{Hash1:hash1, Hash2:hash2})
    }

    return hashes
}

