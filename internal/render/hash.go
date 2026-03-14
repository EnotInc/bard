package render

import "hash/fnv"

func GetHash(s *[]rune) uint32 {
	h := fnv.New32a()
	h.Write([]byte(string(*s)))
	return h.Sum32()
}
