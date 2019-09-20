package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func GenToken() string {
	seed := time.Now().String() + fmt.Sprint(rand.Intn(10000))
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])
}
