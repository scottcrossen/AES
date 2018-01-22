package aes

import (
  "math"
)

const Nb = 4

func Encrypt(input []byte, key []byte) []byte {
  keySchedule := keyExpansion(key)
  Nr := len(keySchedule) / Nb - 1
  state := [][]byte{
    make([]byte, 4),
    make([]byte, 4),
    make([]byte, 4),
    make([]byte, 4),
  }
  for i := 0; i < 4 * Nb; i++ {
    state[i % 4][int(math.Floor(float64(i) / 4))] = input[i]
  }
  addRoundKey(state, keySchedule, 0, Nb)
  for round := 1; round < Nr; round++ {
    subBytes(state, Nb)
    shiftRows(state, Nb)
    mixColumns(state, Nb)
    addRoundKey(state, keySchedule, 0, Nb)
  }
  subBytes(state, Nb)
  shiftRows(state, Nb)
  addRoundKey(state, keySchedule, Nr, Nb)
  output := make([]byte, 4 * Nb)
  for i := 0; i < 4 * Nb; i++ {
    output[i] = state[i % 4][int(math.Floor(float64(i) / 4))]
  }
  return output
}

func subBytes(state [][]byte, Nb int) {
  for r := 0; r < 4; r++ {
    for c := 0; c < Nb; c++ {
      state[r][c] = sbox[state[r][c]]
    }
  }
  // return state
}

func shiftRows(state [][]byte, Nb int) {
  temp := make([]byte, 4)
  for r := 1; r < 4; r++ {
    for c := 0; c < 4; c++ {
      temp[c] = state[r][(c + r) % Nb]
    }
    for c := 0; c < 4; c++ {
      state[r][c] = temp[c]
    }
  }
  // return state
}

func mixColumns(state [][]byte, Nb int) {
  for c := 0; c< Nb; c++ {
    a := make([]byte, Nb)
    b := make([]byte, Nb)
    for r := 0; r < 4; r++ {
      a[r] = state[r][c]
      if state[r][c] & 0x80 != 0x00 {
        b[r] = state[r][c] << 1 ^ 0x1b
      } else {
        b[r] = state[r][c] << 1
      }
      state[0][c] = b[0] ^ a[1] ^ b[1] ^ a[2] ^ a[3]
      state[1][c] = a[0] ^ b[1] ^ a[2] ^ b[2] ^ a[3]
      state[2][c] = a[0] ^ a[1] ^ b[2] ^ a[3] ^ b[3]
      state[3][c] = a[0] ^ b[0] ^ a[1] ^ a[2] ^ b[3]
    }
  }
  // return state
}


func addRoundKey(state [][]byte, keySchedule [][]byte, round int, Nb int) {
  for row := 0; row < 4; row++ {
    for col := 0; col<Nb; col++ {
      state[row][col] ^= keySchedule[round * 4 + col][row]
    }
  }
  //return state
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
  word[len(word) - 1] = word[0]
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
