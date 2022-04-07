package main

import (
	"fmt"
	"math/rand"
)

/*
色子游戏。
获胜分数条件： 先得一百分。
角色： 玩家和对手。
角色的分数： 当前分数和上一轮分数。
角色的行动： 掷色子，停止，获胜。
角色的策略： 最多掷K次色子，停止，获胜。
*/
const (
	获胜分数 = 100 // The winning score in a game of Pig
	轮次   = 10  // The number of games per series to simulate
)

type 得分 struct {
	玩家, 对手, 一手 int
}

type 动作 func(当前 得分) (结果 得分, 结束否 bool)

func 掷色子(分 得分) (得分, bool) {
	点数 := rand.Intn(6) + 1 // A random int in [1, 6]
	if 点数 == 1 {
		return 得分{分.对手, 分.玩家, 0}, true
	}
	return 得分{分.玩家, 分.对手, 点数 + 分.一手}, false
}

func 停止(分 得分) (得分, bool) {
	return 得分{分.对手, 分.玩家 + 分.一手, 0}, true
}

// 策略
type 策略 func(得分) 动作

func 最多掷色子(目标分 int) 策略 {
	return func(分 得分) 动作 {
		if 分.一手 >= 目标分 {
			return 停止
		}
		return 掷色子
	}
}

func 游戏(策略0, 策略1 策略) int {
	策略集合 := []策略{策略0, 策略1}
	var 分 得分
	var 结束否 bool
	当前玩家 := rand.Intn(2) // Randomly decide who plays first
	for 分.玩家+分.一手 < 获胜分数 {
		动作 := 策略集合[当前玩家](分)
		分, 结束否 = 动作(分)
		if 结束否 {
			当前玩家 = (当前玩家 + 1) % 2
		}
	}
	return 当前玩家
}

func 回合(策略集合 []策略) ([]int, int) {
	赢家集合 := make([]int, len(策略集合))
	for i := 0; i < len(策略集合); i++ {
		for j := i + 1; j < len(策略集合); j++ {

			for k := 0; k < 轮次; k++ {
				赢家 := 游戏(策略集合[i], 策略集合[j])
				if 赢家 == 0 {
					赢家集合[i]++
				} else {
					赢家集合[j]++
				}
			}
		}
	}
	游戏数 := 轮次 * (len(策略集合) - 1)
	return 赢家集合, 游戏数
}

func 胜率(vals ...int) string {
	总游戏数 := 0
	for _, v := range vals {
		总游戏数 += v
	}
	var s string
	for _, v := range vals {
		if s != "" {
			s += ", "
		}
		p := 100 * float64(v) / float64(总游戏数)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", v, 总游戏数, p)
	}
	return s

}

/*
色子游戏。
获胜分数条件： 先得一百分。
角色： 玩家和对手。
角色的分数： 当前分数和上一轮分数。
角色的行动： 掷色子，停止，获胜。
角色的策略： 最多掷K次色子，停止，获胜。
*/

func main() {
	策略集合 := make([]策略, 获胜分数)
	for k := range 策略集合 {
		策略集合[k] = 最多掷色子(k + 1)
	}
	赢家集合, 游戏数 := 回合(策略集合)

	for k := range 策略集合 {
		fmt.Printf("策略%d: %s\n",
			k+1,
			胜率(赢家集合[k], 游戏数-赢家集合[k]))
	}

}
