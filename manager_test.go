package tower

import (
    "github.com/stretchr/testify/assert"
    "log"
    "testing"
)

func TestManagerEntity(t *testing.T) {
    tw := NewManager(Options{
        MapWidth: 100, MapHeight: 100,
        TowerHeight: 10, TowerWidth: 10,
        Debug: true,
    })
    tw.Init()

    TE11_ID := ID("11")
    TE11 := &TestEntity{
        Entity:      &Entity{Id: TE11_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE11_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE11_ID),
    }
    TE11Coord := Coord{X: 0, Y: 0}
    log.Println("-> Entity Add: ", TE11_ID)
    assert.Equal(t, true, tw.Add(TE11.Setup(), TE11Coord))
    assert.Equal(t, 0, TE11.enters)
    assert.Equal(t, 0, TE11.leaves)

    TE22_ID := ID("22")
    TE22 := &TestEntity{
        Entity:      &Entity{Id: TE22_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE22_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE22_ID),
    }
    TE22Coord := Coord{X: 5, Y: 5}
    log.Println("-> Entity Add: ", TE22_ID)
    assert.Equal(t, true, tw.Add(TE22.Setup(), TE22Coord))
    assert.Equal(t, 1, TE11.enters)
    assert.Equal(t, 0, TE11.leaves)
    assert.Equal(t, 1, TE22.enters)
    assert.Equal(t, 0, TE22.leaves)
    assert.Equal(t, TE22.enters, TE11.enters)
    assert.Equal(t, TE22.leaves, TE11.leaves)

    TE33_ID := ID("33")
    TE33 := &TestEntity{
        Entity:      &Entity{Id: TE33_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE33_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE33_ID),
    }
    E33Coord := Coord{X: 30, Y: 30}
    log.Println("-> Entity Add: ", TE33_ID)
    assert.Equal(t, true, tw.Add(TE33.Setup(), E33Coord))
    assert.Equal(t, 0, TE33.enters)
    assert.Equal(t, 0, TE33.leaves)
    //
    log.Println("-> Entity Update: ", TE22_ID)
    tw.Update(TE22.Setup(), E33Coord)
    assert.Equal(t, 1, TE33.enters)
    assert.Equal(t, 1, TE22.leaves)
    assert.Equal(t, 1, TE11.leaves)

    log.Println("-> Entity Update: ", TE22_ID)
    tw.Update(TE22.Setup(), TE11Coord)
    assert.Equal(t, 1, TE33.leaves)
    assert.Equal(t, 2, TE11.enters)

    log.Println("-> Entity Remove: ", TE22_ID)
    tw.Remove(TE22.Setup())
    assert.Equal(t, 2, TE11.leaves)
}

func TestManagerWatcher(t *testing.T) {
    tw := NewManager(Options{
        MapWidth: 100, MapHeight: 100,
        TowerHeight: 10, TowerWidth: 10,
        Debug: true,
    })
    tw.Init()
    TWMAP_ID := ID("99")
    TWMAP := &TestWatcher{
        Watcher:     &Watcher{Id: TWMAP_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TWMAP_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TWMAP_ID),
    }
    TWMAPCoord := Coord{X: 50, Y: 50}
    log.Println("-> Watcher Add: ", TWMAP_ID)
    tw.AddWatcher(TWMAP.Setup(), TWMAPCoord, 10)

    TWCenter_ID := ID("88")
    TWCenter := &TestWatcher{
        Watcher:     &Watcher{Id: TWCenter_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TWCenter_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TWCenter_ID),
    }
    TWCenterCoord := Coord{X: 50, Y: 50}
    log.Println("-> Watcher Add: ", TWCenter_ID)
    tw.AddWatcher(TWCenter.Setup(), TWCenterCoord, 1)

    TE11_ID := ID("11")
    TE11 := &TestEntity{
        Entity:      &Entity{Id: TE11_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE11_ID),
        onLeaveFunc: newLogLeaveFunc("Entity", TE11_ID),
    }
    TE11Coord := Coord{X: 0, Y: 0}
    log.Println("-> Entity Add: ", TE11_ID)
    assert.Equal(t, true, tw.Add(TE11.Setup(), TE11Coord))
    assert.Equal(t, 1, TWMAP.enters)
    assert.Equal(t, 0, TWCenter.enters)

    TE22_ID := ID("22")
    TE22 := &TestEntity{
        Entity:      &Entity{Id: TE22_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE22_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE22_ID),
    }
    E22Coord := Coord{X: 30, Y: 30}
    log.Println("-> Entity Add: ", TE22_ID)
    assert.Equal(t, true, tw.Add(TE22.Setup(), E22Coord))
    assert.Equal(t, 2, TWMAP.enters)
    assert.Equal(t, 0, TWCenter.enters)

    TE33_ID := ID("33")
    TE33 := &TestEntity{
        Entity:      &Entity{Id: TE33_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE33_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE33_ID),
    }
    E33Coord := Coord{X: 40, Y: 40}
    log.Println("-> Entity Add: ", TE33_ID)
    assert.Equal(t, true, tw.Add(TE33.Setup(), E33Coord))
    assert.Equal(t, 3, TWMAP.enters)
    assert.Equal(t, 1, TWCenter.enters)

    log.Println("-> Entity Remove: ", TE33_ID)
    assert.Equal(t, true, tw.Remove(TE33.Setup()))
    assert.Equal(t, 1, TWMAP.leaves)
    assert.Equal(t, 1, TWCenter.leaves)

    log.Println("-> Watcher Remove: ", TWMAP)
    tw.ClearWatcher(TWMAP.Watcher)

    TE44_ID := ID("44")
    TE44 := &TestEntity{
        Entity:      &Entity{Id: TE44_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE44_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE44_ID),
    }
    E44Coord := Coord{X: 40, Y: 40}
    log.Println("-> Entity Add: ", TE44_ID)
    assert.Equal(t, true, tw.Add(TE44.Setup(), E44Coord))
    assert.Equal(t, 3, TWMAP.enters)
    assert.Equal(t, 2, TWCenter.enters)
    log.Println("-> Entity Remove: ", TE44_ID)
    assert.Equal(t, true, tw.Remove(TE44.Setup()))
    assert.Equal(t, 1, TWMAP.leaves)
    assert.Equal(t, 2, TWCenter.leaves)

    TE55_ID := ID("55")
    TE55 := &TestEntity{
        Entity:      &Entity{Id: TE55_ID},
        onEnterFunc: newLogEnterFunc("Watcher", TE55_ID),
        onLeaveFunc: newLogLeaveFunc("Watcher", TE55_ID),
    }
    E55Coord := Coord{X: 40, Y: 40}
    log.Println("-> Entity Add: ", TE55_ID)
    assert.Equal(t, true, tw.Add(TE55.Setup(), E55Coord))
    assert.Equal(t, 3, TWMAP.enters)
    assert.Equal(t, 3, TWCenter.enters)
    log.Println("-> Entity Remove: ", TE55_ID)
    assert.Equal(t, true, tw.Remove(TE55.Setup()))
    assert.Equal(t, 1, TWMAP.leaves)
    assert.Equal(t, 3, TWCenter.leaves)
}

//// Entity

var _ EntityCallback = (*TestEntity)(nil)

type TestEntity struct {
    *Entity
    enters      int
    leaves      int
    onEnterFunc func(entity *Entity)
    onLeaveFunc func(entity *Entity)
}

func (t *TestEntity) Setup() *Entity {
    t.Entity.Callback = t
    return t.Entity
}

func (t *TestEntity) OnEntityEnter(other *Entity) {
    t.enters++
    if t.onEnterFunc != nil {
        t.onEnterFunc(other)
    }
}

func (t *TestEntity) OnEntityLeave(other *Entity) {
    t.leaves++
    if t.onLeaveFunc != nil {
        t.onLeaveFunc(other)
    }
}

//// Watcher

var _ WatcherCallback = (*TestWatcher)(nil)

type TestWatcher struct {
    *Watcher
    enters      int
    leaves      int
    onEnterFunc func(entity *Entity)
    onLeaveFunc func(entity *Entity)
}

func (t *TestWatcher) Setup() *Watcher {
    t.Watcher.Callback = t
    return t.Watcher
}

func (t *TestWatcher) OnWatchingEnter(other *Entity) {
    t.enters++
    if t.onEnterFunc != nil {
        t.onEnterFunc(other)
    }
}

func (t *TestWatcher) OnWatchingLeave(other *Entity) {
    t.leaves++
    if t.onLeaveFunc != nil {
        t.onLeaveFunc(other)
    }
}

func newLogEnterFunc(kind string, id ID) func(entity *Entity) {
    return func(entity *Entity) {
        log.Printf("=> %s[%s] receive ENTER, target: %+v\n", kind, id, entity)
    }
}

func newLogLeaveFunc(kind string, id ID) func(entity *Entity) {
    return func(entity *Entity) {
        log.Printf("=> %s[%s] receive LEAVE, target: %+v\n", kind, id, entity)
    }
}
