package tower

import (
    "log"
    "math"
)

type Options struct {
    MapWidth    float32 // width of map in pixel
    MapHeight   float32 // height of map in pixel
    TowerWidth  float32 // width of tower in pixel
    TowerHeight float32 // height of tower in pixel
    Debug       bool    // Is debug
}

type Coord struct {
    X float32
    Y float32
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
    towerSizeX := int(math.Ceil(float64(m.opts.MapWidth / m.opts.TowerWidth)))
    towerSizeY := int(math.Ceil(float64(m.opts.MapHeight / m.opts.TowerHeight)))
    m.max = ICoord{
        X: towerSizeX - 1,
        Y: towerSizeY - 1,
    }
    m.towers = make([][]*Tower, towerSizeX)
    for x := 0; x < towerSizeX; x++ {
        m.towers[x] = make([]*Tower, towerSizeY)
        for y := 0; y < towerSizeY; y++ {
            m.towers[x][y] = NewTower(ICoord{X: x, Y: y}, m.opts.Debug)
        }
    }
    if m.opts.Debug {
        log.Println("INFO: Tower manager init, options=", m.opts, "max=", m.max, "towers=", towerSizeX*towerSizeY)
    }
}

func (m *Manager) Options() Options {
    return m.opts
}

func (m *Manager) MaxCoord() ICoord {
    return m.max
}

func (m *Manager) TowerCount() int {
    return m.max.X * m.max.Y
}

// Add 在指定地图像素坐标的位置添加一个Entity。
//
// 此操作会触发OnEntityEnter函数回调
//
// @return 仅当成功添加entity时返回True，否则返回False
func (m *Manager) Add(entity *Entity, position Coord) bool {
    verifyEntity(entity)
    if !m.check(position) {
        log.Println("ERROR: Tower manager add entity, coord INVALID, position=", position.String(), "entity=", entity)
        return false
    }
    coord := m.convToTowerCoord(position)
    return m.towers[coord.X][coord.Y].add(entity)
}

// Remove 从指定地图像素坐标位置移除Entity
//
// 此操作会触发OnEntityLeave函数回调
//
// @return 仅当成功删除entity时返回True，否则返回False
func (m *Manager) Remove(entity *Entity) bool {
    verifyEntity(entity)
    if entity.tower == nil {
        return false
    }
    coord := *entity.tower
    return m.towers[coord.X][coord.Y].remove(entity)
}

// Update 将Entity从移动到新的像素坐标位置。
//
// 此操作会触发OnEntityEnter和OnEntityLeave函数回调
//
// @return 仅当成功更新entity时返回True，否则返回False
func (m *Manager) Update(entity *Entity, targetCoord Coord) bool {
    verifyEntity(entity)
    if !m.check(targetCoord) {
        log.Println("ERROR: Tower manager update entity, coord INVALID, target=", targetCoord, "entity=", entity)
        return false
    }
    prevCoord := *entity.tower
    nextCoord := m.convToTowerCoord(targetCoord)
    // 检查是否合法
    if prevCoord.X > len(m.towers) || nextCoord.X > len(m.towers) {
        log.Println("ERROR: Tower manager update entity, prev.tower=", prevCoord, "next.tower", nextCoord)
        return false
    }
    // Tower坐标没有发生切换
    if prevCoord.X == nextCoord.X && prevCoord.Y == nextCoord.Y {
        return false
    }
    m.towers[prevCoord.X][prevCoord.Y].remove(entity) // Prev tower
    m.towers[nextCoord.X][nextCoord.Y].add(entity)    // Next tower
    return true
}

// UpdateCoord 将Entity从指定地图像素坐标位置移动到新的像素坐标位置。
//
// 此操作会触发OnEntityEnter和OnEntityLeave函数回调
//
// @return 仅当成功更新entity时返回True，否则返回False
//func (m *Manager) UpdateCoord(entity *Entity, from, to Coord) bool {
//    verifyEntity(entity)
//    if !m.check(from) || !m.check(to) {
//        log.Println("ERROR: Tower manager update entity, coord INVALID, from=", from, "to=", to, "entity=", entity)
//        return false
//    }
//    prevCoord := m.convToTowerCoord(from)
//    nextCoord := m.convToTowerCoord(to)
//    return m.update(entity, prevCoord, nextCoord)
//}

// AddWatcher 从指定地图像素坐标位置，以及Tower距离，将Watcher添加到范围内的Tower列表
func (m *Manager) AddWatcher(watcher *Watcher, position Coord, towerDistance int) {
    verifyWatcher(watcher)
    m.searchTowers(position, towerDistance, func(tower *Tower) {
        tower.addWatcher(watcher)
    })
}

// RemoveWatcher 从指定地图像素坐标位置，以及Tower距离，移除范围内Tower绑定的Watcher列表
func (m *Manager) RemoveWatcher(watcher *Watcher, position Coord, towerDistance int) {
    verifyWatcher(watcher)
    m.searchTowers(position, towerDistance, func(tower *Tower) {
        tower.removeWatcher(watcher)
    })
}

// ClearWatcher 清除指定Watcher与已绑定Tower的关系
func (m *Manager) ClearWatcher(watcher *Watcher) {
    verifyWatcher(watcher)
    for tower := range watcher.watching {
        tower.removeWatcher(watcher)
    }
}

func (m *Manager) searchTowers(position Coord, dist int, onTower func(tower *Tower)) {
    coord := m.convToTowerCoord(position)
    start, end := m.coordRangeOf(coord, dist, m.max)
    for x := start.X; x <= end.X; x++ {
        for y := start.Y; y <= end.Y; y++ {
            onTower(m.towers[x][y])
        }
    }
}

func (m *Manager) coordRangeOf(coord ICoord, dist int, max ICoord) (start ICoord, end ICoord) {
    if coord.X-dist < 0 {
        start.X = 0
        end.X = 2 * dist
    } else if coord.X+dist > max.X {
        end.X = max.X
        start.X = max.X - 2*dist
    } else {
        start.X = coord.X - dist
        end.X = coord.X + dist
    }

    if coord.Y-dist < 0 {
        start.Y = 0
        end.Y = 2 * dist
    } else if coord.Y+dist > max.Y {
        end.Y = max.Y
        start.Y = max.Y - 2*dist
    } else {
        start.Y = coord.Y - dist
        end.Y = coord.Y + dist
    }
    start.X = iMax(start.X, 0)
    end.X = iMin(end.X, max.X)
    start.Y = iMax(start.Y, 0)
    end.Y = iMin(end.Y, max.Y)
    return
}

func (m *Manager) convToTowerCoord(pos Coord) ICoord {
    return ICoord{
        X: int(math.Floor(float64(pos.X / m.opts.TowerWidth))),
        Y: int(math.Floor(float64(pos.Y / m.opts.TowerHeight))),
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
