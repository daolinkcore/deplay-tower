package tower

import "fmt"

func (o Options) String() string {
    return fmt.Sprintf(
        "Opitons{width: %.2f, height: %.2f, tower.w: %.2f, tower.h: %.2f}",
        o.MapWidth, o.MapHeight, o.TowerWidth, o.TowerHeight,
    )
}

func (o Position) String() string {
    return fmt.Sprintf(
        "Position{x: %.2f, y: %.2f}", o.X, o.Y,
    )
}

func (o coord) String() string {
    return fmt.Sprintf(
        "x: %d, y: %d", o.X, o.Y,
    )
}

func (o *Entity) String() string {
    return fmt.Sprintf(
        "Entity{id: %v, value: %+v}", o.Id, o.Value,
    )
}

func (o *Watcher) String() string {
    return fmt.Sprintf(
        "Watcher{id: %v, value: %+v}", o.Id, o.Value,
    )
}
