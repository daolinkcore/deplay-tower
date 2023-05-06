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

    TE11_ID := ID(11)
    TE11 := &TestEntity{
        Entity:      &Entity{Id: TE11_ID},
        onEnterFunc: newLogEnterFunc(TE11_ID),
        onLeaveFunc: newLogLeaveFunc(TE11_ID),
    }
    TE11Coord := Coord{X: 0, Y: 0}
    log.Println("-> Entity Add: ", TE11_ID)
    assert.Equal(t, true, tw.Add(TE11.Setup(), TE11Coord))
    assert.Equal(t, 0, TE11.enters)
    assert.Equal(t, 0, TE11.leaves)

    TE22_ID := ID(22)
    TE22 := &TestEntity{
        Entity:      &Entity{Id: TE22_ID},
        onEnterFunc: newLogEnterFunc(TE22_ID),
        onLeaveFunc: newLogLeaveFunc(TE22_ID),
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

    TE33_ID := ID(33)
    TE33 := &TestEntity{
        Entity:      &Entity{Id: TE33_ID},
        onEnterFunc: newLogEnterFunc(TE33_ID),
        onLeaveFunc: newLogLeaveFunc(TE33_ID),
    }
    E33Coord := Coord{X: 30, Y: 30}
    log.Println("-> Entity Add: ", TE33_ID)
    assert.Equal(t, true, tw.Add(TE33.Setup(), E33Coord))
    assert.Equal(t, 0, TE33.enters)
    assert.Equal(t, 0, TE33.leaves)
    //
    log.Println("-> Entity Update: ", TE22_ID)
    tw.Update(TE22.Setup(), TE22Coord, E33Coord)
    assert.Equal(t, 1, TE33.enters)
    assert.Equal(t, 1, TE22.leaves)
    assert.Equal(t, 1, TE11.leaves)

    log.Println("-> Entity Update: ", TE22_ID)
    tw.Update(TE22.Setup(), E33Coord, TE11Coord)
    assert.Equal(t, 1, TE33.leaves)
    assert.Equal(t, 2, TE11.enters)

    log.Println("-> Entity Remove: ", TE22_ID)
    tw.Remove(TE22.Setup())
    assert.Equal(t, 2, TE11.leaves)
}

var _ Callback = (*TestEntity)(nil)

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

func newLogEnterFunc(id ID) func(entity *Entity) {
    return func(entity *Entity) {
        log.Printf("=> %d receive ENTER, target: %+v\n", id, entity)
    }
}

func newLogLeaveFunc(id ID) func(entity *Entity) {
    return func(entity *Entity) {
        log.Printf("=> %d receive LEAVE, target: %+v\n", id, entity)
    }
}
