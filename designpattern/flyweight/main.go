package main

import (
	"fmt"
	"time"
)

// 享元模式是一种结构型设计模式，它通过共享对象来减少内存使用和提高性能
// 以下示例，使用享元模式，使每场比赛结果对象共享球队对象，而不需要每场比赛结果都创建两个个球队对象

type TeamId uint8

const (
	Warrior TeamId = iota
	Laker
)

// Team 球队
type Team struct {
	// 球队id
	Id TeamId
	// 球队名称
	Name    string
	Players []*Player
}

type Player struct {
	Name string
	Team TeamId
}

// 比赛结果
type Match struct {
	Date         time.Time
	LocalTeam    *Team //主队
	VisitorTeam  *Team //客队
	LocalScore   uint8 //主队得分
	VisitorScore uint8 //客队得分
}

func (m *Match) ShowResult() {
	fmt.Printf("%s VS %s - %d:%d\n", m.LocalTeam.Name, m.VisitorTeam.Name,
		m.LocalScore, m.VisitorScore)
}

// 队伍池
type TeamPool struct {
	teams map[TeamId]*Team
}

func (tp *TeamPool) get(id TeamId) *Team {
	team, ok := tp.teams[id]
	if !ok {
		team = createTeam(id)
		tp.teams[id] = team
	}
	return team
}

// 创建队伍
func createTeam(id TeamId) *Team {
	switch id {
	case Warrior:
		w := &Team{
			Id:   Warrior,
			Name: "Golden State Warriors",
		}
		curry := &Player{
			Name: "Stephen Curry",
			Team: Warrior,
		}
		thompson := &Player{
			Name: "Klay Thompson",
			Team: Warrior,
		}
		w.Players = append(w.Players, curry, thompson)
		return w
	case Laker:
		l := &Team{
			Id:   Laker,
			Name: "Los Angeles Lakers",
		}
		james := &Player{
			Name: "LeBron James",
			Team: Laker,
		}
		davis := &Player{
			Name: "Anthony Davis",
			Team: Laker,
		}
		l.Players = append(l.Players, james, davis)
		return l
	default:
		fmt.Printf("Get an invalid team id %v.\n", id)
		return nil
	}
}

// 享元工厂
var factory = &TeamPool{teams: map[TeamId]*Team{}}

func main() {
	game1 := &Match{
		Date:         time.Date(2020, 1, 10, 9, 30, 0, 0, time.Local),
		LocalTeam:    factory.get(Warrior),
		VisitorTeam:  factory.get(Laker),
		LocalScore:   102,
		VisitorScore: 99,
	}
	game1.ShowResult()
	game2 := &Match{
		Date:         time.Date(2020, 1, 10, 9, 30, 0, 0, time.Local),
		LocalTeam:    factory.get(Laker),
		VisitorTeam:  factory.get(Warrior),
		LocalScore:   102,
		VisitorScore: 99,
	}
	game2.ShowResult()

	fmt.Println("Team Warrior object pattern", game1.LocalTeam == game2.VisitorTeam)
}
