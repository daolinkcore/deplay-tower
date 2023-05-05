package tower

import "log"

type (
    ID   int64
    TYPE int
)

type Entity struct {
    Id       ID                  // 地图唯一EntityID
    Type     TYPE                // 类型
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
    coord    ICoord
    ids      map[ID]*Entity
    watchers map[TYPE]map[ID]*Entity
    objects  map[TYPE]map[ID]*Entity
    size     int
}

func NewTower(coord ICoord, debug bool) *Tower {
    return &Tower{
        coord:    coord,
        debug:    debug,
        ids:      make(map[ID]*Entity),
        watchers: make(map[TYPE]map[ID]*Entity),
        objects:  make(map[TYPE]map[ID]*Entity),
    }
}

func (t *Tower) add(entity *Entity) bool {
    if _, ok := t.ids[entity.Id]; ok {
        return false
    }
    objmap, ok := t.objects[entity.Type]
    if !ok {
        objmap = make(map[ID]*Entity)
        t.objects[entity.Type] = objmap
    }
    objmap[entity.Id] = entity
    t.ids[entity.Id] = entity
    entity.tower = &t.coord
    if t.debug {
        log.Println("DEBUG: Tower(", t.coord, ") -> Add entity:", entity)
    }
    t.forEachCallableEntity(func(obj *Entity) {
        if obj != entity {
            obj.Callback.OnEntityEnter(entity)
        }
    })
    return true
}

func (t *Tower) remove(entity *Entity) bool {
    if _, ok := t.ids[entity.Id]; !ok {
        return false
    }
    if objmap, ok := t.objects[entity.Type]; ok {
        delete(objmap, entity.Id)
    }
    entity.tower = nil
    delete(t.ids, entity.Id)
    t.size -= 1
    t.forEachCallableEntity(func(obj *Entity) {
        if obj != entity {
            obj.Callback.OnEntityLeave(entity)
        }
    })
    return true
}

func (t *Tower) forEachCallableEntity(callback func(node *Entity)) {
    for _, obj := range t.ids {
        callback(obj)
    }
    for _, wchmap := range t.watchers {
        for _, watcher := range wchmap {
            callback(watcher)
        }
    }
}

func (t *Tower) addWatcher(watcher *Entity) {
    wchmap, ok := t.watchers[watcher.Type]
    if !ok {
        wchmap = make(map[ID]*Entity)
        t.watchers[watcher.Type] = wchmap
    }
    if _, ok := wchmap[watcher.Id]; ok {
        return
    }
    wchmap[watcher.Id] = watcher
    // Link watching
    if watcher.watching == nil {
        watcher.watching = make(map[*Tower]struct{})
    }
    watcher.watching[t] = struct{}{}
}

func (t *Tower) removeWatcher(watcher *Entity) {
    wchmap, ok := t.watchers[watcher.Type]
    if !ok {
        return
    }
    delete(wchmap, watcher.Id)
    delete(watcher.watching, t)
}
