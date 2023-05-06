package tower

import "log"

type (
    ID int64
)

type Entity struct {
    Id       ID       // 地图唯一EntityID
    Value    any      // 功能扩展字段
    Callback Callback // 事件回调
    tower    *ICoord  // 当前所在Tower坐标
}

type Watcher struct {
    Id       ID                  // 地图唯一EntityID
    Value    any                 // 功能扩展字段
    Callback Callback            // 事件回调
    watching map[*Tower]struct{} // 已Watched的Tower列表（仅Watcher使用）
}

// Callback Entity进入退出回调接口
type Callback interface {
    // OnEntityEnter 当Entity进入当前Tower坐标范围时，回调此函数
    OnEntityEnter(other *Entity)

    // OnEntityLeave 当Entity离开当前Tower坐标范围时，回调此函数
    OnEntityLeave(other *Entity)

    //OnEntityChanged(other *Entity)
}

type Tower struct {
    debug    bool
    coord    ICoord // Tower坐标
    entities map[ID]*Entity
    watchers map[ID]*Watcher
}

func NewTower(coord ICoord, debug bool) *Tower {
    return &Tower{
        coord:    coord,
        debug:    debug,
        entities: make(map[ID]*Entity),
        watchers: make(map[ID]*Watcher),
    }
}

func (t *Tower) add(entity *Entity) bool {
    if _, ok := t.entities[entity.Id]; ok {
        return false
    }
    t.entities[entity.Id] = entity
    entity.tower = &t.coord
    if t.debug {
        log.Println("DEBUG: Tower(", t.coord, ") -> Add entity:", entity)
    }
    for _, watcher := range t.watchers {
        watcher.Callback.OnEntityEnter(entity)
    }
    for _, exists := range t.entities {
        if exists == entity {
            continue
        }
        entity.Callback.OnEntityEnter(exists)
        exists.Callback.OnEntityEnter(entity)
    }
    return true
}

func (t *Tower) remove(entity *Entity) bool {
    if _, ok := t.entities[entity.Id]; !ok {
        return false
    }
    entity.tower = nil
    delete(t.entities, entity.Id)
    for _, watcher := range t.watchers {
        watcher.Callback.OnEntityLeave(entity)
    }
    for _, remain := range t.entities {
        entity.Callback.OnEntityLeave(remain)
        remain.Callback.OnEntityLeave(entity)
    }
    return true
}

func (t *Tower) addWatcher(watcher *Watcher) {
    if _, ok := t.watchers[watcher.Id]; ok {
        return
    }
    t.watchers[watcher.Id] = watcher
    if watcher.watching == nil {
        watcher.watching = make(map[*Tower]struct{})
    }
    watcher.watching[t] = struct{}{}
    if t.debug {
        log.Println("DEBUG: Tower(", t.coord, ") -> Add watcher:", watcher)
    }
}

func (t *Tower) removeWatcher(watcher *Watcher) {
    if _, ok := t.watchers[watcher.Id]; !ok {
        return
    }
    delete(t.watchers, watcher.Id)
    delete(watcher.watching, t)
    if t.debug {
        log.Println("DEBUG: Tower(", t.coord, ") -> Remove watcher:", watcher)
    }
}
