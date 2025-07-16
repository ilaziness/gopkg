package main

import (
	"fmt"
	"time"
)

// 备忘录模式 - 在不破坏封装性的前提下，捕获一个对象的内部状态，并在该对象之外保存这个状态

// 备忘录接口
type Memento interface {
	GetState() string
	GetTimestamp() time.Time
}

// 具体备忘录
type ConcreteMemento struct {
	state     string
	timestamp time.Time
}

func NewMemento(state string) *ConcreteMemento {
	return &ConcreteMemento{
		state:     state,
		timestamp: time.Now(),
	}
}

func (cm *ConcreteMemento) GetState() string {
	return cm.state
}

func (cm *ConcreteMemento) GetTimestamp() time.Time {
	return cm.timestamp
}

// 发起人 - 文本编辑器
type TextEditor struct {
	content string
}

func NewTextEditor() *TextEditor {
	return &TextEditor{content: ""}
}

func (te *TextEditor) Write(text string) {
	te.content += text
	fmt.Printf("写入文本: '%s'\n", text)
	fmt.Printf("当前内容: '%s'\n", te.content)
}

func (te *TextEditor) GetContent() string {
	return te.content
}

func (te *TextEditor) SetContent(content string) {
	te.content = content
}

// 创建备忘录
func (te *TextEditor) CreateMemento() Memento {
	fmt.Printf("创建备忘录，保存状态: '%s'\n", te.content)
	return NewMemento(te.content)
}

// 从备忘录恢复
func (te *TextEditor) RestoreFromMemento(memento Memento) {
	te.content = memento.GetState()
	fmt.Printf("从备忘录恢复，内容: '%s' (保存时间: %s)\n",
		te.content, memento.GetTimestamp().Format("15:04:05"))
}

// 管理者 - 历史记录管理器
type HistoryManager struct {
	mementos []Memento
	current  int
}

func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		mementos: make([]Memento, 0),
		current:  -1,
	}
}

func (hm *HistoryManager) SaveState(memento Memento) {
	// 如果当前不在最新位置，删除后面的历史记录
	if hm.current < len(hm.mementos)-1 {
		hm.mementos = hm.mementos[:hm.current+1]
	}

	hm.mementos = append(hm.mementos, memento)
	hm.current++
	fmt.Printf("保存历史记录 #%d\n", hm.current+1)
}

func (hm *HistoryManager) Undo() Memento {
	if hm.current > 0 {
		hm.current--
		fmt.Printf("撤销到历史记录 #%d\n", hm.current+1)
		return hm.mementos[hm.current]
	}
	fmt.Println("无法撤销，已经是最早的状态")
	return nil
}

func (hm *HistoryManager) Redo() Memento {
	if hm.current < len(hm.mementos)-1 {
		hm.current++
		fmt.Printf("重做到历史记录 #%d\n", hm.current+1)
		return hm.mementos[hm.current]
	}
	fmt.Println("无法重做，已经是最新的状态")
	return nil
}

func (hm *HistoryManager) ShowHistory() {
	fmt.Println("=== 历史记录 ===")
	for i, memento := range hm.mementos {
		marker := "  "
		if i == hm.current {
			marker = "→ "
		}
		fmt.Printf("%s#%d: '%s' (%s)\n",
			marker, i+1, memento.GetState(),
			memento.GetTimestamp().Format("15:04:05"))
	}
	fmt.Println()
}

// 游戏存档示例
type GameState struct {
	level  int
	score  int
	health int
	x, y   int
}

func (gs *GameState) String() string {
	return fmt.Sprintf("Level:%d Score:%d Health:%d Position:(%d,%d)",
		gs.level, gs.score, gs.health, gs.x, gs.y)
}

// 游戏存档备忘录
type GameMemento struct {
	state     *GameState
	timestamp time.Time
	name      string
}

func NewGameMemento(state *GameState, name string) *GameMemento {
	// 深拷贝状态
	stateCopy := &GameState{
		level:  state.level,
		score:  state.score,
		health: state.health,
		x:      state.x,
		y:      state.y,
	}
	return &GameMemento{
		state:     stateCopy,
		timestamp: time.Now(),
		name:      name,
	}
}

func (gm *GameMemento) GetGameState() *GameState {
	return gm.state
}

func (gm *GameMemento) GetName() string {
	return gm.name
}

func (gm *GameMemento) GetTimestamp() time.Time {
	return gm.timestamp
}

// 游戏角色
type GameCharacter struct {
	state *GameState
}

func NewGameCharacter() *GameCharacter {
	return &GameCharacter{
		state: &GameState{
			level:  1,
			score:  0,
			health: 100,
			x:      0,
			y:      0,
		},
	}
}

func (gc *GameCharacter) Move(dx, dy int) {
	gc.state.x += dx
	gc.state.y += dy
	fmt.Printf("移动到位置: (%d, %d)\n", gc.state.x, gc.state.y)
}

func (gc *GameCharacter) GainScore(points int) {
	gc.state.score += points
	fmt.Printf("获得 %d 分，总分: %d\n", points, gc.state.score)
}

func (gc *GameCharacter) TakeDamage(damage int) {
	gc.state.health -= damage
	if gc.state.health < 0 {
		gc.state.health = 0
	}
	fmt.Printf("受到 %d 点伤害，剩余生命: %d\n", damage, gc.state.health)
}

func (gc *GameCharacter) LevelUp() {
	gc.state.level++
	gc.state.health = 100 // 升级回满血
	fmt.Printf("升级到 %d 级，生命值恢复满血\n", gc.state.level)
}

func (gc *GameCharacter) GetState() *GameState {
	return gc.state
}

func (gc *GameCharacter) CreateSave(name string) *GameMemento {
	fmt.Printf("创建存档: %s\n", name)
	return NewGameMemento(gc.state, name)
}

func (gc *GameCharacter) LoadSave(memento *GameMemento) {
	gc.state = &GameState{
		level:  memento.state.level,
		score:  memento.state.score,
		health: memento.state.health,
		x:      memento.state.x,
		y:      memento.state.y,
	}
	fmt.Printf("加载存档: %s\n", memento.GetName())
	fmt.Printf("当前状态: %s\n", gc.state)
}

// 游戏存档管理器
type SaveManager struct {
	saves map[string]*GameMemento
}

func NewSaveManager() *SaveManager {
	return &SaveManager{
		saves: make(map[string]*GameMemento),
	}
}

func (sm *SaveManager) SaveGame(name string, memento *GameMemento) {
	sm.saves[name] = memento
	fmt.Printf("存档 '%s' 已保存\n", name)
}

func (sm *SaveManager) LoadGame(name string) *GameMemento {
	if memento, exists := sm.saves[name]; exists {
		return memento
	}
	fmt.Printf("存档 '%s' 不存在\n", name)
	return nil
}

func (sm *SaveManager) ListSaves() {
	fmt.Println("=== 存档列表 ===")
	for name, memento := range sm.saves {
		fmt.Printf("存档: %s - %s (%s)\n",
			name, memento.state,
			memento.timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== 文本编辑器备忘录模式示例 ===")

	// 创建文本编辑器和历史管理器
	editor := NewTextEditor()
	history := NewHistoryManager()

	// 保存初始状态
	history.SaveState(editor.CreateMemento())

	// 编辑文本
	editor.Write("Hello")
	history.SaveState(editor.CreateMemento())

	editor.Write(" World")
	history.SaveState(editor.CreateMemento())

	editor.Write("!")
	history.SaveState(editor.CreateMemento())

	fmt.Println()
	history.ShowHistory()

	// 撤销操作
	fmt.Println("=== 撤销操作 ===")
	if memento := history.Undo(); memento != nil {
		editor.RestoreFromMemento(memento)
	}

	if memento := history.Undo(); memento != nil {
		editor.RestoreFromMemento(memento)
	}

	fmt.Println()
	history.ShowHistory()

	// 重做操作
	fmt.Println("=== 重做操作 ===")
	if memento := history.Redo(); memento != nil {
		editor.RestoreFromMemento(memento)
	}

	fmt.Println()
	history.ShowHistory()

	fmt.Println("\n=== 游戏存档备忘录模式示例 ===")

	// 创建游戏角色和存档管理器
	character := NewGameCharacter()
	saveManager := NewSaveManager()

	fmt.Printf("初始状态: %s\n", character.GetState())

	// 游戏进程
	character.Move(10, 5)
	character.GainScore(100)
	saveManager.SaveGame("开始游戏", character.CreateSave("开始游戏"))

	character.Move(20, 15)
	character.GainScore(200)
	character.TakeDamage(30)
	saveManager.SaveGame("第一关", character.CreateSave("第一关"))

	character.LevelUp()
	character.Move(50, 30)
	character.GainScore(500)
	saveManager.SaveGame("第二关", character.CreateSave("第二关"))

	// 模拟游戏失败
	fmt.Println("\n--- 游戏遇到困难，角色死亡 ---")
	character.TakeDamage(100)
	fmt.Printf("当前状态: %s\n", character.GetState())

	// 显示存档列表
	saveManager.ListSaves()

	// 加载存档
	fmt.Println("=== 加载存档 ===")
	if save := saveManager.LoadGame("第一关"); save != nil {
		character.LoadSave(save)
	}

	// 继续游戏
	fmt.Println("\n--- 从存档继续游戏 ---")
	character.GainScore(50)
	character.Move(5, 5)
	fmt.Printf("继续游戏后状态: %s\n", character.GetState())
}
