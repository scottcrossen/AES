package aes

import (
  "math"
  "fmt"
)

func Encrypt(input []byte, key []byte) []byte {
  Log(fmt.Sprintf("\nPLAINTEXT:\t\t%x", input))
  Log(fmt.Sprintf("KEY:\t\t%x", key))
  Log("CIPHER (ENCRYPT):")
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
  Log(fmt.Sprintf("round[ %d].input\t%x", 0, input))
  Log(fmt.Sprintf("round[ %d].k_sch\t%x", 0, makeOutputInv(keySchedule, Nb)))
  addRoundKey(state, keySchedule, 0, Nb)
  for round := 1; round < Nr; round++ {
    roundString := getRoundString(round)
    Log(fmt.Sprintf("%s.start\t%x", roundString, makeOutput(state, Nb)))
    subBytes(state, Nb)
    Log(fmt.Sprintf("%s.s_box\t%x", roundString, makeOutput(state, Nb)))
    shiftRows(state, Nb)
    Log(fmt.Sprintf("%s.s_row\t%x", roundString, makeOutput(state, Nb)))
    mixColumns(state, Nb)
    Log(fmt.Sprintf("%s.m_col\t%x", roundString, makeOutput(state, Nb)))
    addRoundKey(state, keySchedule, round, Nb)
    Log(fmt.Sprintf("%s.k_sch\t%x", roundString, makeOutputInv(keySchedule[(round * 4):], Nb)))
  }
  roundString := getRoundString(Nr)
  Log(fmt.Sprintf("%s.start\t%x", roundString, makeOutput(state, Nb)))
  subBytes(state, Nb)
  Log(fmt.Sprintf("%s.s_box\t%x", roundString, makeOutput(state, Nb)))
  shiftRows(state, Nb)
  Log(fmt.Sprintf("%s.s_row\t%x", roundString, makeOutput(state, Nb)))
  addRoundKey(state, keySchedule, Nr, Nb)
  Log(fmt.Sprintf("%s.k_sch\t%x", roundString, makeOutputInv(keySchedule[(10 * 4):], Nb)))
  Log(fmt.Sprintf("%s.output\t%x", roundString, makeOutput(state, Nb)))
  return makeOutput(state, Nb)
}

func subBytes(state [][]byte, Nb int) {
  for r := 0; r < 4; r++ {
    for c := 0; c < Nb; c++ {
      state[r][c] = sbox[state[r][c]]
    }
  }
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
}

func mixColumns(state [][]byte, Nb int) {
  a := make([]byte, 4)
  for i := 0; i < 4; i++ {
    a[0] = state[0][i]
    a[1] = state[1][i]
    a[2] = state[2][i]
    a[3] = state[3][i]

    state[0][i] = gfMul(0x02, a[0]) ^ gfMul(0x03, a[1]) ^ a[2] ^ a[3]
    state[1][i] = gfMul(0x02, a[1]) ^ gfMul(0x03, a[2]) ^ a[3] ^ a[0]
    state[2][i] = gfMul(0x02, a[2]) ^ gfMul(0x03, a[3]) ^ a[0] ^ a[1]
    state[3][i] = gfMul(0x02, a[3]) ^ gfMul(0x03, a[0]) ^ a[1] ^ a[2]
  }
}
