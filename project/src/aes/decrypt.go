package aes

import (
  "math"
  "fmt"
)

func Decrypt(input []byte, key []byte) []byte {
  Log(fmt.Sprintf("\nCIPHERTEXT:\t%x", input))
  Log(fmt.Sprintf("KEY:\t\t%x", key))
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
  for r := 1; r < 4; r++ {
    for c := 0; c < 4; c++ {
      temp[c] = state[r][(4 + c - r) % Nb]
    }
    for c := 0; c < 4; c++ {
      state[r][c] = temp[c]
    }
  }
}

func subBytesInv(state [][]byte, Nb int) {
  for r := 0; r < 4; r++ {
    for c := 0; c < Nb; c++ {
      state[r][c] = isbox[state[r][c]]
    }
  }
}

func mixColumnsInv(state [][]byte, Nb int) {
  a := make([]byte, 4)
  for i := 0; i < 4; i++ {
    a[0] = state[0][i]
    a[1] = state[1][i]
    a[2] = state[2][i]
    a[3] = state[3][i]
    state[0][i] = gfMul(0x0E, a[0]) ^ gfMul(0x0B, a[1]) ^ gfMul(0x0D, a[2]) ^ gfMul(0x09, a[3])
    state[1][i] = gfMul(0x0E, a[1]) ^ gfMul(0x0B, a[2]) ^ gfMul(0x0D, a[3]) ^ gfMul(0x09, a[0])
    state[2][i] = gfMul(0x0E, a[2]) ^ gfMul(0x0B, a[3]) ^ gfMul(0x0D, a[0]) ^ gfMul(0x09, a[1])
    state[3][i] = gfMul(0x0E, a[3]) ^ gfMul(0x0B, a[0]) ^ gfMul(0x0D, a[1]) ^ gfMul(0x09, a[2])
  }
}
