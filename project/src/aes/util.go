package aes

import (
  "fmt"
  "math"
)

func gfMul(a, b byte) byte {
  var prod byte = 0
  var h byte
  for i := 0; i < 8; i++ {
    if (b & 0x01) != 0 {
      prod ^= a
    }
    h = a >> 7
    a <<= 1
    if h != 0 {
      a ^= 0x1B
    }
    b >>= 1
  }
  return prod
}

func addRoundKey(state [][]byte, keySchedule [][]byte, round int, Nb int) {
  for row := 0; row < 4; row++ {
    for col := 0; col < Nb; col++ {
      state[row][col] ^= keySchedule[round * 4 + col][row]
    }
  }
}

func getRoundString(round int) string {
  if round < 10 {
    return fmt.Sprintf("round[ %d]", round)
  } else {
    return fmt.Sprintf("round[%d]", round)
  }
}

func makeOutput(state [][]byte, Nb int) []byte {
  output := make([]byte, 4 * Nb)
  for i := 0; i < 4 * Nb; i++ {
    output[i] = state[i % 4][int(math.Floor(float64(i) / 4))]
  }
  return output
}

func makeOutputInv(keySchedule [][]byte, Nb int) []byte {
  output := make([]byte, 4 * Nb)
  for i := 0; i < 4 * Nb; i++ {
    output[i] = keySchedule[int(math.Floor(float64(i) / 4))][i % 4]
  }
  return output
}

func subWord(word []byte) []byte {
  output := make([]byte, len(word))
  for i := 0; i < len(word); i++ {
    output[i] = sbox[word[i]]
  }
  return output
}

func rotWord(word []byte) []byte {
  output := make([]byte, len(word))
  for i := 1; i < len(word); i++ {
    output[i - 1] = word[i]
  }
  output[len(word) - 1] = word[0]
  return output
}

func keyExpansion(key []byte) [][]byte {
  Nk := len(key) / 4
  Nr := Nk + 6
  output := make([][]byte, Nb * (Nr + 1))
  temp := make([]byte, 4)
  for i := 0; i < Nk; i++ {
    output[i] = []byte {
      key[4 * i],
      key[4 * i + 1],
      key[4 * i + 2],
      key[4 * i + 3],
    }
  }
  for i := Nk; i < len(output); i++ {
    output[i] = make([]byte, 4)
    for t := 0; t < 4; t++ {
      temp[t] = output[i - 1][t]
    }
    if (i % Nk == 0) {
      temp = subWord(rotWord(temp))
      for t := 0; t < 4; t++ {
        temp[t] ^= rcon[i / Nk][t]
      }
    } else if Nk > 6 && i % Nk == 4 {
      temp = subWord(temp)
    }
    for t := 0; t < 4; t++ {
      output[i][t] = output[i - Nk][t] ^ temp[t]
    }
  }
  return output
}
