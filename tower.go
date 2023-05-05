package tower

import "log"

type (
    ID int64
)

type Entity struct {
    Id       ID                  // 地图唯一EntityID
    Value    any                 // 功能扩展字段
    Callback Callback            // 事件回调
    tower    *ICoord             // 当前所在Tower坐标
    watching map[*Tower]struct{} // 已Watched的Tower列表（仅Watcher使用）
}

// Callback Entity进入退出回调接口
type Callback interface {
    // OnEntityEnter 当Entity进入当前Tower坐标范围时，回调此函数
    OnEntityEnter(other *Entity)

    // OnEntityLeave 当Entity离开当前Tower坐标范围时，回调此函数
    OnEntityLeave(other *Entity)
}

type Tower struct {
    debug    bool
    coord    ICoord // Tower坐标
    entities map[ID]*Entity
    watchers map[ID]*Entity
}

func NewTower(coord ICoord, debug bool) *Tower {
    return &Tower{
        coord:    coord,
        debug:    debug,
        entities: make(map[ID]*Entity),
        watchers: make(map[ID]*Entity),
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
    t.foreachCallableEntity(func(callable *Entity) {
        if callable != entity {
            callable.Callback.OnEntityEnter(entity)
        }
    })
    for _, exists := range t.entities {
        if exists != entity {
            entity.Callback.OnEntityEnter(exists)
        }
    }
    return true
}

func (t *Tower) remove(entity *Entity) bool {
    if _, ok := t.entities[entity.Id]; !ok {
        return false
    }
    entity.tower = nil
    delete(t.entities, entity.Id)
    t.foreachCallableEntity(func(callable *Entity) {
        if callable != entity {
            callable.Callback.OnEntityLeave(entity)
        }
    })
    for _, exists := range t.entities {
        if exists != entity {
            entity.Callback.OnEntityLeave(exists)
        }
    }
    return true
}

func (t *Tower) foreachCallableEntity(callback func(node *Entity)) {
    for _, entity := range t.entities {
        callback(entity)
    }
    for _, watcher := range t.watchers {
        callback(watcher)
    }
}

func (t *Tower) addWatcher(watcher *Entity) {
    if _, ok := t.watchers[watcher.Id]; ok {
        return
    }
    t.watchers[watcher.Id] = watcher
    // Link watching
    if watcher.watching == nil {
        watcher.watching = make(map[*Tower]struct{})
    }
    watcher.watching[t] = struct{}{}
}

func (t *Tower) removeWatcher(watcher *Entity) {
    if _, ok := t.watchers[watcher.Id]; !ok {
        return
    }
    delete(t.watchers, watcher.Id)
    delete(watcher.watching, t)
}
