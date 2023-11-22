package main

import (
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
)

// function to hash a string 
func hash_string(str string) [32]byte {
	return sha256.Sum256([]byte(str))
} 

// function to hash a string array
func hash_string_array(str_array []string) [32]byte {
	hash := sha256.New()
	for _, str := range str_array {
		hash.Write([]byte(str))
	}
	return sha256.Sum256(hash.Sum(nil))
}

// function to calculate number of trailing zeros in a 32 byte hash
func count_trailing_zeros(hash [32]byte) int {
	zero_count := 0
	for i := 31; i >= 0; i-- {
		byte := hash[i]
		// iterate through each bit of the byte
		for j := 0; j < 8; j++ {
			if byte & 1 == 0 {
				zero_count++
			} else {
				return zero_count
			}
			byte = byte >> 1
		}
	}
	return zero_count
}

// function to convert a 32 byte hash to a string
func hash_to_string(hash [32]byte) string {
	return hex.EncodeToString(hash[:])
}

// function to generate random nonce 
func generate_nonce() [32]byte {
	nonce := make([]byte, 32)
	rand.Read(nonce)
	var static_nonce [32]byte
	copy(static_nonce[:], nonce[:])
	return static_nonce
}

// function to concatenate byte arrays 
func concatenate_hashes(arrays ...[]byte) [32]byte {
	num := len(arrays)
	result := make([]byte, 0 , num*32)
	for _, array := range arrays {
		result = append(result, array...)
	}
	
	// hash the concatenated hashes
	return sha256.Sum256(result)
}

