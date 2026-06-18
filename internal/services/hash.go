package services

import "hash/fnv"

// NOTE: this is save, for now.
var hash = fnv.New32()

func GetHash(s string) uint32 {
	hash.Reset()
	hash.Write([]byte(s))
	return hash.Sum32()
}
