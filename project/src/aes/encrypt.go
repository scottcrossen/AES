package aes

import (
  "math"
  "fmt"
)

func Encrypt(input []byte, key []byte) []byte {
  Log(fmt.Sprintf("\nPLAINTEXT:\t\t%x", input))
  Log(fmt.Sprintf("KEY:\t\t\t%x", key))
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
  Log(fmt.Sprintf("round[ %d].input\t\t%x", 0, input))
  Log(fmt.Sprintf("round[ %d].k_sch\t\t%x", 0, makeOutputInv(keySchedule, Nb)))
  addRoundKey(state, keySchedule, 0, Nb)
  for round := 1; round < Nr; round++ {
    roundString := getRoundString(round)
    Log(fmt.Sprintf("%s.start\t\t%x", roundString, makeOutput(state, Nb)))
    subBytes(state, Nb)
    Log(fmt.Sprintf("%s.s_box\t\t%x", roundString, makeOutput(state, Nb)))
    shiftRows(state, Nb)
    Log(fmt.Sprintf("%s.s_row\t\t%x", roundString, makeOutput(state, Nb)))
    mixColumns(state, Nb)
    Log(fmt.Sprintf("%s.m_col\t\t%x", roundString, makeOutput(state, Nb)))
    addRoundKey(state, keySchedule, round, Nb)
    Log(fmt.Sprintf("%s.k_sch\t\t%x", roundString, makeOutputInv(keySchedule[(round * 4):], Nb)))
  }
  roundString := getRoundString(Nr)
  Log(fmt.Sprintf("%s.start\t\t%x", roundString, makeOutput(state, Nb)))
  subBytes(state, Nb)
  Log(fmt.Sprintf("%s.s_box\t\t%x", roundString, makeOutput(state, Nb)))
  shiftRows(state, Nb)
  Log(fmt.Sprintf("%s.s_row\t\t%x", roundString, makeOutput(state, Nb)))
  addRoundKey(state, keySchedule, Nr, Nb)
  Log(fmt.Sprintf("%s.k_sch\t\t%x", roundString, makeOutputInv(keySchedule[(10 * 4):], Nb)))
  Log(fmt.Sprintf("%s.output\t%x", roundString, makeOutput(state, Nb)))
  return makeOutput(state, Nb)
}

func subBytes(state [][]byte, Nb int) {
  for row := 0; row < 4; row++ {
    for col := 0; col < Nb; col++ {
      state[row][col] = sbox[state[row][col]]
    }
  }
}

func shiftRows(state [][]byte, Nb int) {
  temp := make([]byte, 4)
  for row := 1; row < 4; row++ {
    for col := 0; col < 4; col++ {
      temp[col] = state[row][(col + row) % Nb]
    }
    for col := 0; col < 4; col++ {
      state[row][col] = temp[col]
    }
  }
}

func mixColumns(state [][]byte, Nb int) {
  temp := make([]byte, 4)
  for i := 0; i < 4; i++ {
    temp[0] = state[0][i]
    temp[1] = state[1][i]
    temp[2] = state[2][i]
    temp[3] = state[3][i]
    state[0][i] = ffMultiply(0x02, temp[0]) ^ ffMultiply(0x03, temp[1]) ^ temp[2] ^ temp[3]
    state[1][i] = ffMultiply(0x02, temp[1]) ^ ffMultiply(0x03, temp[2]) ^ temp[3] ^ temp[0]
    state[2][i] = ffMultiply(0x02, temp[2]) ^ ffMultiply(0x03, temp[3]) ^ temp[0] ^ temp[1]
    state[3][i] = ffMultiply(0x02, temp[3]) ^ ffMultiply(0x03, temp[0]) ^ temp[1] ^ temp[2]
  }
}
