package aes

import (
  "math"
  "fmt"
)

func Decrypt(input []byte, key []byte) []byte {
  Log(fmt.Sprintf("\nCIPHERTEXT:\t\t%x", input))
  Log(fmt.Sprintf("KEY:\t\t\t%x", key))
  Log("INVERSE CIPHER (DECRYPT):")
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
  Log(fmt.Sprintf("round[ %d].iinput\t%x", 0, input))
  Log(fmt.Sprintf("round[ %d].ik_sch\t%x", 0, makeOutputInv(keySchedule[Nr * 4:], Nb)))
  addRoundKey(state, keySchedule, Nr, Nb)
  for round := 1; round < Nr; round++ {
    roundString := getRoundString(round)
    Log(fmt.Sprintf("%s.istart\t%x", roundString, makeOutput(state, Nb)))
    shiftRowsInv(state, Nb)
    Log(fmt.Sprintf("%s.is_row\t%x", roundString, makeOutput(state, Nb)))
    subBytesInv(state, Nb)
    Log(fmt.Sprintf("%s.is_box\t%x", roundString, makeOutput(state, Nb)))
    Log(fmt.Sprintf("%s.ik_sch\t%x", roundString, makeOutputInv(keySchedule[(Nr - round) * 4:], Nb)))
    addRoundKey(state, keySchedule, Nr - round, Nb)
    Log(fmt.Sprintf("%s.ik_add\t%x", roundString, makeOutput(state, Nb)))
    mixColumnsInv(state, Nb)
  }
  roundString := getRoundString(Nr)
  Log(fmt.Sprintf("%s.istart\t%x", roundString, makeOutput(state, Nb)))
  shiftRowsInv(state, Nb)
  Log(fmt.Sprintf("%s.is_row\t%x", roundString, makeOutput(state, Nb)))
  subBytesInv(state, Nb)
  Log(fmt.Sprintf("%s.is_box\t%x", roundString, makeOutput(state, Nb)))
  Log(fmt.Sprintf("%s.ik_sch\t%x", roundString, makeOutputInv(keySchedule, Nb)))
  addRoundKey(state, keySchedule, 0, Nb)
  Log(fmt.Sprintf("%s.ioutput\t%x", roundString, makeOutput(state, Nb)))
  return makeOutput(state, Nb)
}

func shiftRowsInv(state [][]byte, Nb int) {
  temp := make([]byte, 4)
  for row := 1; row < 4; row++ {
    for col := 0; col < 4; col++ {
      temp[col] = state[row][(4 + col - row) % Nb]
    }
    for col := 0; col < 4; col++ {
      state[row][col] = temp[col]
    }
  }
}

func subBytesInv(state [][]byte, Nb int) {
  for row := 0; row < 4; row++ {
    for col := 0; col < Nb; col++ {
      state[row][col] = isbox[state[row][col]]
    }
  }
}

func mixColumnsInv(state [][]byte, Nb int) {
  temp := make([]byte, 4)
  for i := 0; i < 4; i++ {
    temp[0] = state[0][i]
    temp[1] = state[1][i]
    temp[2] = state[2][i]
    temp[3] = state[3][i]
    state[0][i] = ffMultiply(0x0E, temp[0]) ^ ffMultiply(0x0B, temp[1]) ^ ffMultiply(0x0D, temp[2]) ^ ffMultiply(0x09, temp[3])
    state[1][i] = ffMultiply(0x0E, temp[1]) ^ ffMultiply(0x0B, temp[2]) ^ ffMultiply(0x0D, temp[3]) ^ ffMultiply(0x09, temp[0])
    state[2][i] = ffMultiply(0x0E, temp[2]) ^ ffMultiply(0x0B, temp[3]) ^ ffMultiply(0x0D, temp[0]) ^ ffMultiply(0x09, temp[1])
    state[3][i] = ffMultiply(0x0E, temp[3]) ^ ffMultiply(0x0B, temp[0]) ^ ffMultiply(0x0D, temp[1]) ^ ffMultiply(0x09, temp[2])
  }
}
