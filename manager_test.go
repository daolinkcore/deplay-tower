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

    E11 := &TestEntity{
        Entity: &Entity{Id: 11},
        OnEnterFunc: func(entity *Entity) {
            log.Println("=> 11 receive ENTER, target:", entity)
        },
        OnLeaveFunc: func(entity *Entity) {
            log.Println("=> 11 receive LEAVE, target:", entity)
        },
    }
    E11Coord := Coord{X: 0, Y: 0}
    log.Println("-> 11 add")
    assert.Equal(t, true, tw.Add(E11.Setup(), E11Coord))

    E22 := &TestEntity{
        Entity: &Entity{Id: 22},
        OnEnterFunc: func(entity *Entity) {
            log.Println("=> 22 receive ENTER, target:", entity)
        },
        OnLeaveFunc: func(entity *Entity) {
            log.Println("=> 22 receive LEAVE, target:", entity)
        },
    }
    E22Coord := Coord{X: 5, Y: 5}
    log.Println("-> 22 add")
    assert.Equal(t, true, tw.Add(E22.Setup(), E22Coord))

    E33 := &TestEntity{
        Entity: &Entity{Id: 33},
        OnEnterFunc: func(entity *Entity) {
            log.Println("=> 33 receive ENTER, target:", entity)
        },
        OnLeaveFunc: func(entity *Entity) {
            log.Println("=> 33 receive LEAVE, target:", entity)
        },
    }
    E33Coord := Coord{X: 30, Y: 30}
    log.Println("-> 33 add")
    assert.Equal(t, true, tw.Add(E33.Setup(), E33Coord))

    log.Println("-> 22 update")
    tw.Update(E22.Setup(), E22Coord, E33Coord)

    log.Println("-> 22 update")
    tw.Update(E22.Setup(), E33Coord, E11Coord)

    log.Println("-> 22 remove")
    tw.Remove(E22.Setup())
}

var _ Callback = (*TestEntity)(nil)

type TestEntity struct {
    *Entity
    OnEnterFunc func(entity *Entity)
    OnLeaveFunc func(entity *Entity)
}

func (t *TestEntity) Setup() *Entity {
    t.Entity.Callback = t
    return t.Entity
}

func (t *TestEntity) OnEntityEnter(other *Entity) {
    if t.OnEnterFunc != nil {
        t.OnEnterFunc(other)
    }
}

func (t *TestEntity) OnEntityLeave(other *Entity) {
    if t.OnLeaveFunc != nil {
        t.OnLeaveFunc(other)
    }
}
