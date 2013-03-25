package main

import ("fmt"; "math/rand"; "time")

var affectedCubies = [][]int{
  []int {  0,  1,  2,  3,  0,  1,  2,  3 },   // U
  []int {  4,  7,  6,  5,  4,  5,  6,  7 },   // D
  []int {  0,  9,  4,  8,  0,  3,  5,  4 },   // F
  []int {  2, 10,  6, 11,  2,  1,  7,  6 },   // B
  []int {  3, 11,  7,  9,  3,  2,  6,  5 },   // L
  []int {  1,  8,  5, 10,  1,  0,  4,  7 },  // R
}

var phaseMoves = [][]int{
    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
    []int{0, 1, 2, 3, 4, 5, 7, 10, 12, 13, 14, 15, 16, 17},
    []int{0, 1, 2, 3, 4, 5, 7, 10, 13, 16},
    []int{1, 4, 7, 10, 13, 16},
}

type Cube struct {
    state []int
}

func intbool(b bool) int {
    if b {
        return 1
    }
    return 0
}

func (cube *Cube) id(phase int) []int{
    switch {
    case phase == 0 :
        return cube.state[20:32]
    case phase == 1 :
        result := make([]int, len(cube.state[31:40]))
        copy(result,cube.state[31:40])
        for i:=0; i < 12; i++ {
            result[0] |= (cube.state[i] / 8) << uint8(i)
        }
        return result
    case phase == 2:
        result := []int{0,0,0}
        for e,v := range cube.state[:8] {
            result[0] |= (v & 1) << uint8(2*e)
        }
        for e := range cube.state[8:12] {
            result[0] |= 2 << uint8(2*e)
        }
        for c := range cube.state[:8] {
            result[1] |= ((cube.state[c+12]-12) & 5) << uint(3*c)
        }
        for i := range cube.state[12:20] {
            for j := range cube.state[i+1:20] {
                result[2] ^= intbool(cube.state[i] > cube.state[j])
            }
        }
        return result
    default:
        return cube.state
    }
    return cube.state
}

func (cube *Cube) doMove(move int) Cube {
    newstate := make([]int,len(cube.state))
    copy(newstate,cube.state)
    oldstate := make([]int,len(cube.state))

    turns := move % 3 + 1;
    face := move / 3;
    for turn := 0; turn < turns; turn++ {
        copy(oldstate,newstate)
        for i,v := range affectedCubies[face][:8] {
            isCorner := intbool(i > 3)
            target := v + isCorner*12
            killer := isCorner * 12
            if (i&3) == 3 {
                killer += affectedCubies[face][i - 3]
            } else {
                killer += affectedCubies[face][i+1]
            }
            newstate[target] = oldstate[killer]
            orientationDelta := 0
            if i < 4 {
                orientationDelta = intbool(face>1 && face < 4)
            } else if face >= 2 {
                orientationDelta = 2 - (i&1)
            }
            newstate[target+20] = oldstate[killer+20] + orientationDelta;
            if turn == turns - 1 {
                newstate[target + 20] %= 2 + isCorner
            }
        }
    }
  return Cube{newstate}
}

func goalCube() Cube {
	state:=make([]int,40)
	for i := 0; i < 20; i++ {
		state[i] = i
	}
	return Cube{state}
}

func intsliceEqual(s1 []int, s2 []int) bool {
    if len(s1) != len(s2) {
        return false
    }
    for i := range s1 {
        if s1[i] != s2[i] {
            return false
        }
    }
    return true
}


func main() {
    rand.Seed(time.Now().UnixNano())
	//fmt.Println(affectedCubies)
    //fmt.Println(phaseMoves)
    goal := goalCube()
    current := goalCube()
    for i:=0; i < 30; i++ {
        current = current.doMove(rand.Intn(18))
    }
    fmt.Println(current)

    for phase := 0; phase < 4; phase++ {
        currentId := current.id(phase)
        goalId := goal.id(phase)
        states := []Cube {current}
        fmt.Println("start phase: ", phase)
        fmt.Println("currentId: ",currentId)
        fmt.Println("current cube: ", current)
        for phaseDone := intsliceEqual(currentId,goalId); !phaseDone; {
            fmt.Println(phase)
            nextStates := make([]Cube,0,len(phaseMoves[phase]))
            for _,curr := range states {
                for _,move := range phaseMoves[phase] {
                    nextState := curr.doMove(move)
                    //fmt.Println(nextState.id(phase),goalId)
                    if phaseDone = intsliceEqual(nextState.id(phase),goalId); phaseDone {
                        current = nextState
                        break
                    } else {
                        nextStates = append(nextStates, nextState)
                    }
                }
                if phaseDone {
                    break
                }
            }
            states = nextStates
        }
    }
}
