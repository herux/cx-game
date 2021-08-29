package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/events"
	"github.com/skycoin/cx-game/world"
)

const (
	drillSpeed     float32 = 0.5
	drillJumpSpeed float32 = 15
)

func AiHandlerDrill(agent *agents.Agent, ctx AiContext) {
	dist := ctx.PlayerPos.X() - agent.PhysicsState.Pos.X
	directionX := math32.Sign(dist)
	if math32.Abs(dist) > ctx.WorldWidth/2 {
		directionX *= -1
	}
	agent.PhysicsState.Direction = directionX * -1
	if math32.Abs(dist) > 1 {
		agent.AnimationPlayback.PlayRepeating("Walk")
		agent.PhysicsState.Vel.X = directionX * walkSpeed
	} else {
		agent.AnimationPlayback.PlayRepeating("Attack")
		agent.PhysicsState.Vel.X = 0
	}

	isCollisionHorizontal := agent.PhysicsState.Collisions.Horizontal()
	if isCollisionHorizontal {
		events.OnSpiderCollisionHorizontal.Trigger(events.SpiderEventData{
			Agent: agent,
		})
		ctx.World.Planet.DamageTile(int(agent.PhysicsState.Pos.X), int(agent.PhysicsState.Pos.Y+0.5), world.TopLayer)
	}

	// doJump :=
	// 	agent.PhysicsState.Collisions.Horizontal() &&
	// 		agent.PhysicsState.IsOnGround() && !agent.PhysicsState.Collisions.VerticalAbove()
	// if doJump {
	// 	events.OnSpiderBeforeJump.Trigger(events.SpiderEventData{
	// 		Agent: agent,
	// 	})

	// 	agent.PhysicsState.Vel.Y = drillJumpSpeed
	// 	// trigger an event when spiderdrill jump
	// 	events.OnSpiderJump.Trigger(events.SpiderEventData{
	// 		Agent: agent,
	// 	})
	// } else {
	// 	agent.PhysicsState.Vel.Y -= constants.Gravity * constants.PHYSICS_TICK
	// }
}
