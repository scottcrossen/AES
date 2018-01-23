package aes

import (
  "math"
  "fmt"
)

const Nb = 4

func Encrypt(input []byte, key []byte) []byte {
  Log(fmt.Sprintf("PLAINTEXT:\t\t%x", input))
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
  addRoundKey(state, keySchedule, 0, Nb)
  Log(fmt.Sprintf("round[ %d].input\t%x", 0, input))
  Log(fmt.Sprintf("round[ %d].k_sch\t%x", 0, keySchedule[0]))
  for round := 1; round < Nr; round++ {
    var logRound string
    if round < 10 {
      logRound = fmt.Sprintf("round[ %d]", round)
    } else {
      logRound = fmt.Sprintf("round[%d]", round)
    }
    Log(fmt.Sprintf("%s.start\t%x", logRound, makeOutput(state, Nb)))
    subBytes(state, Nb)
    Log(fmt.Sprintf("%s.s_box\t%x", logRound, makeOutput(state, Nb)))
    shiftRows(state, Nb)
    Log(fmt.Sprintf("%s.s_row\t%x", logRound, makeOutput(state, Nb)))
    mixColumns(state, Nb)
    Log(fmt.Sprintf("%s.m_col\t%x", logRound, makeOutput(state, Nb)))
    addRoundKey(state, keySchedule, round, Nb)
    Log(fmt.Sprintf("%s.k_sch\t%x", logRound, makeOutputInv(keySchedule[(round * 4):], Nb)))
  }
  Log(fmt.Sprintf("%s.start\t%x", "round[10]", makeOutput(state, Nb)))
  subBytes(state, Nb)
  Log(fmt.Sprintf("%s.s_box\t%x", "round[10]", makeOutput(state, Nb)))
  shiftRows(state, Nb)
  Log(fmt.Sprintf("%s.s_row\t%x", "round[10]", makeOutput(state, Nb)))
  addRoundKey(state, keySchedule, Nr, Nb)
  Log(fmt.Sprintf("%s.k_sch\t%x", "round[10]", makeOutputInv(keySchedule[(10 * 4):], Nb)))
  Log(fmt.Sprintf("%s.output\t%x", "round[10]", makeOutput(state, Nb)))
  return makeOutput(state, Nb)
}
func makeOutput(state [][]byte, Nb int) []byte {
  output := make([]byte, 4 * Nb)
  for i := 0; i < 4 * Nb; i++ {
    output[i] = state[i % 4][int(math.Floor(float64(i) / 4))]
  }
  return output
}
func makeOutputInv(keySchedule [][]byte, Nb int) []byte {
  /*for i := 0; i < len(keySchedule); i++ {
    fmt.Println(fmt.Sprintf("%x", keySchedule[i]))
  }*/
  output := make([]byte, 4 * Nb)
  for i := 0; i < 4 * Nb; i++ {
    output[i] = keySchedule[int(math.Floor(float64(i) / 4))][i % 4]
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

func mixColumns(s [][]byte, Nb int) {
  a := make([]byte, 4)
  for i := 0; i < 4; i++ {
    a[0] = s[0][i]
    a[1] = s[1][i]
    a[2] = s[2][i]
    a[3] = s[3][i]

    s[0][i] = gfMul(0x02, a[0]) ^ gfMul(0x03, a[1]) ^ a[2] ^ a[3]
    s[1][i] = gfMul(0x02, a[1]) ^ gfMul(0x03, a[2]) ^ a[3] ^ a[0]
    s[2][i] = gfMul(0x02, a[2]) ^ gfMul(0x03, a[3]) ^ a[0] ^ a[1]
    s[3][i] = gfMul(0x02, a[3]) ^ gfMul(0x03, a[0]) ^ a[1] ^ a[2]
  }
}

func addRoundKey(state [][]byte, keySchedule [][]byte, round int, Nb int) {
  for row := 0; row < 4; row++ {
    for col := 0; col < Nb; col++ {
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
