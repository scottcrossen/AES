package aes

import (
  "fmt"
)

func Log(a ...interface{}) (n int, err error) {
  if true {
    return fmt.Println(a...)
  } else {
    return 0, nil
  }
}
