// a world item is an item that is "floating" in the world
package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/sound"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
)

func InitWorldItem() {
	sound.LoadSound("bloop", "bloop.wav")
}

const pickupRadius = 3
const attractRadius = 4
const attractForceMag = 3
const worldItemSize = 0.5

type WorldItem struct {
	physics.Body
	ItemTypeId ItemTypeID
}

var worldItems = []*WorldItem{}

func CreateWorldItem(itemTypeId ItemTypeID, pos mgl32.Vec2) {
	item := WorldItem{
		Body: physics.Body{
			Size: cxmath.Vec2{X: worldItemSize, Y: worldItemSize},
			Pos:  cxmath.Vec2{X: pos.X(), Y: pos.Y()},
		},
		ItemTypeId: itemTypeId,
	}
	physics.RegisterBody(&item.Body)
	worldItems = append(worldItems, &item)
}

func TickWorldItems(
	planet *world.Planet, dt float32, playerPos cxmath.Vec2,
) (forPlayer []*WorldItem) {
	newWorldItems := []*WorldItem{}
	forPlayer = []*WorldItem{}
	for _, item := range worldItems {
		pickedUp := item.Tick(planet, dt, playerPos)
		if pickedUp {
			forPlayer = append(forPlayer, item)
		} else {
			newWorldItems = append(newWorldItems, item)
		}
	}
	worldItems = newWorldItems
	return forPlayer
}

func DrawWorldItems(cam *camera.Camera) {
	for _, item := range worldItems {
		item.Draw(cam)
	}
}

func (item WorldItem) Draw(cam *camera.Camera) {
	if !cam.IsInBoundsF(item.Pos.X, item.Pos.Y) {
		return
	}
	spriteId := itemTypes[item.ItemTypeId].SpriteID
	z := -spriteloader.SpriteRenderDistance
	view := mgl32.Translate3D(-cam.X, -cam.Y, 0)
	world := mgl32.Translate3D(
		item.Pos.X,
		item.Pos.Y,
		z,
	).Mul4(cxmath.Scale(worldItemSize))
	modelView := view.Mul4(world)
	spriteloader.DrawSpriteQuadMatrix(
		modelView, spriteId,
		spriteloader.NewDrawOptions(),
	)
}

func (item *WorldItem) Tick(
	planet *world.Planet, dt float32,
	playerPos cxmath.Vec2,
) bool {
	item.Vel.Y -= physics.Gravity * dt / 2

	itemToPlayerDisplacement := playerPos.Sub(item.Pos)
	itemToPlayerDistSqr := itemToPlayerDisplacement.LengthSqr()
	if itemToPlayerDistSqr < attractRadius*attractRadius {
		attractDirection := itemToPlayerDisplacement.
			Mult(1 / itemToPlayerDisplacement.LengthSqr())
		attractForce := attractDirection.Mult(attractForceMag * dt)
		item.Vel = item.Vel.Add(attractForce)
	}

	//item.Move(planet, dt)
	didPickup := itemToPlayerDistSqr < pickupRadius*pickupRadius
	if didPickup {
		sound.PlaySound("bloop")
	}
	return didPickup
}
