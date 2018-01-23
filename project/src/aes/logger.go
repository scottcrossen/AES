package aes

import (
  "fmt"
)

const verbose = 1

func Log(a ...interface{}) (n int, err error) {
  if (verbose >= 0) {
    return fmt.Println(a...)
  } else {
    return 0, nil
  }
}
