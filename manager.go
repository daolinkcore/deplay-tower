package tower

import (
    "log"
    "math"
)

type Options struct {
    MapWidth    float64 // 地图宽度
    MapHeight   float64 // 地图高度
    TowerWidth  float64 // Tower宽度
    TowerHeight float64 // Tower高度
    Debug       bool
}

type Coord struct {
    X float64
    Y float64
}

type ICoord struct {
    X int
    Y int
}

type Manager struct {
    opts   Options
    towers [][]*Tower
    max    ICoord
}

func NewManager(options Options) *Manager {
    return &Manager{
        opts: options,
    }
}

func (m *Manager) Init() {
    wSize := int(math.Ceil(m.opts.MapWidth / m.opts.TowerWidth))
    hSize := int(math.Ceil(m.opts.MapHeight / m.opts.TowerHeight))
    m.max = ICoord{
        X: wSize - 1,
        Y: hSize - 1,
    }
    m.towers = make([][]*Tower, wSize)
    for i := 0; i < wSize; i++ {
        m.towers[i] = make([]*Tower, hSize)
        for j := 0; j < hSize; j++ {
            m.towers[i][j] = NewTower(ICoord{X: i, Y: j}, m.opts.Debug)
        }
    }
    if m.opts.Debug {
        log.Println("INFO: Tower manager init, options=", m.opts, "max=", m.max, "towers=", wSize*hSize)
    }
}

// Add 在指定地图坐标位置添加Entity。
// 此操作会触发OnEntityEnter函数回调
// @return 仅当成功添加entity时返回True，否则返回False
func (m *Manager) Add(entity *Entity, position Coord) bool {
    verifyEntity(entity)
    if !m.check(position) {
        log.Println("ERROR: Tower manager add entity, coord INVALID, position=", position.String(), "entity=", entity)
        return false
    }
    tc := m.convToTowerCoord(position)
    return m.towers[tc.X][tc.Y].add(entity)
}

// Remove 从指定地图坐标位置移除Entity
// 此操作会触发OnEntityLeave函数回调
// @return 仅当成功删除entity时返回True，否则返回False
func (m *Manager) Remove(entity *Entity) bool {
    verifyEntity(entity)
    if entity.tower != nil {
        p := *entity.tower
        return m.towers[p.X][p.Y].remove(entity)
    } else {
        return false
    }
}

// Update 将Entity从指定地图坐标位置移动到新坐标位置
// 此操作会触发OnEntityEnter和OnEntityLeave函数回调
// @return 仅当成功更新entity时返回True，否则返回False
func (m *Manager) Update(entity *Entity, from, to Coord) bool {
    verifyEntity(entity)
    if !m.check(from) || !m.check(to) {
        return false
    }
    tc1 := m.convToTowerCoord(from)
    tc2 := m.convToTowerCoord(to)
    if tc1.X == tc2.X && tc1.Y == tc2.Y {
        return false
    }
    if tc1.X > len(m.towers) || tc2.X > len(m.towers) {
        log.Println("ERROR: Tower manager update entity, old.pos=", from, "old.at=", tc1, "new.pos=", to, "new.at", tc2)
        return false
    }
    oldTower := m.towers[tc1.X][tc1.Y]
    newTower := m.towers[tc2.X][tc2.Y]
    oldTower.remove(entity)
    newTower.add(entity)
    return true
}

// AddWatcher 从指定地图坐标位置，以及Tower距离，添加Watcher到范围内的Tower列表
func (m *Manager) AddWatcher(watcher *Watcher, position Coord, towerDistance int) {
    verifyWatcher(watcher)
    m.searchTowers(position, towerDistance, func(tower *Tower) {
        tower.addWatcher(watcher)
    })
}

// RemoveWatcher 从指定地图坐标位置，以及Tower距离，移除范围内Tower绑定的Watcher列表
func (m *Manager) RemoveWatcher(watcher *Watcher, position Coord, towerDistance int) {
    verifyWatcher(watcher)
    m.searchTowers(position, towerDistance, func(tower *Tower) {
        tower.removeWatcher(watcher)
    })
}

// ClearWatcher 清除指定Watcher全部绑定已绑定关系
func (m *Manager) ClearWatcher(watcher *Watcher) {
    verifyWatcher(watcher)
    for tower := range watcher.watching {
        tower.removeWatcher(watcher)
    }
}

func (m *Manager) searchTowers(position Coord, dist int, onEach func(tower *Tower)) {
    ip := m.convToTowerCoord(position)
    start, end := m.coordRangeOf(ip, dist, m.max)
    for i := start.X; i <= end.X; i++ {
        for j := start.Y; j <= end.Y; j++ {
            onEach(m.towers[i][j])
        }
    }
}

func (m *Manager) coordRangeOf(pos ICoord, dist int, max ICoord) (start ICoord, end ICoord) {
    if pos.X-dist < 0 {
        start.X = 0
        end.X = 2 * dist
    } else if pos.X+dist > max.X {
        end.X = max.X
        start.X = max.X - 2*dist
    } else {
        start.X = pos.X - dist
        end.X = pos.X + dist
    }

    if pos.Y-dist < 0 {
        start.Y = 0
        end.Y = 2 * dist
    } else if pos.Y+dist > max.Y {
        end.Y = max.Y
        start.Y = max.Y - 2*dist
    } else {
        start.Y = pos.Y - dist
        end.Y = pos.Y + dist
    }
    start.X = iMax(start.X, 0)
    end.X = iMin(end.X, max.X)
    start.Y = iMax(start.Y, 0)
    end.Y = iMin(end.Y, max.Y)
    return
}

func (m *Manager) convToTowerCoord(pos Coord) ICoord {
    return ICoord{
        X: int(math.Floor(pos.X / m.opts.TowerWidth)),
        Y: int(math.Floor(pos.Y / m.opts.TowerHeight)),
    }
}

func (m *Manager) check(coord Coord) bool {
    if coord.X < 0 || coord.Y < 0 || coord.X >= m.opts.MapWidth || coord.Y >= m.opts.MapHeight {
        return false
    }
    return true
}

func verifyEntity(entity *Entity) {
    if entity.Callback == nil {
        log.Fatalln("ERROR: Tower manager verify, nil CALLBACK function, entity=", entity)
    }
}

func verifyWatcher(watcher *Watcher) {
    if watcher.Callback == nil {
        log.Fatalln("ERROR: Tower manager verify, nil CALLBACK function, watcher=", watcher)
    }
}

func iMax(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func iMin(a, b int) int {
    if a < b {
        return a
    }
    return b
}
