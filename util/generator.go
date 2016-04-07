package util

import(
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomPassword() string {
  b := make([]byte, 10)
  for i := range b {
      b[i] = letterBytes[rand.Intn(len(letterBytes))]
  }
  return string(b)
}