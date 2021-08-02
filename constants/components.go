package constants

import "github.com/skycoin/cx-game/components/types"

const (
	AI_HANDLER_NULL types.AgentAiHandlerID = iota
	AI_HANDLER_WALK
	AI_HANDLER_LEAP
	AI_HANDLER_DRILL

	NUM_AI_HANDLERS
)

//agent drawhandler constants
const (
	DRAW_HANDLER_NULL types.AgentDrawHandlerID = iota
	DRAW_HANDLER_QUAD
	DRAW_HANDLER_ANIM

	NUM_AGENT_DRAW_HANDLERS // DO NOT SET THIS MANUALLY
)

//agent physics constants
const (
	PHYSICS_HANDLER_NULL types.AgentPhysicsHandlerID = iota

	NUM_AGENT_PHYSICS_HANDLERS // DO NOT SET THIS MANUALLY
)

//particle drawhandler constants
const (
	PARTICLE_DRAW_HANDLER_NULL types.ParticleDrawHandlerId = iota
	PARTICLE_DRAW_HANDLER_SOLID
	PARTICLE_DRAW_HANDLER_TRANSPARENT
	PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED

	NUM_PARTICLE_DRAW_HANDLERS // DO NOT SET THIS MANUALLY
)

//particle physicshandler constants
const (
	PARTICLE_PHYSICS_HANDLER_NULL types.ParticlePhysicsHandlerID = iota
	PARTICLE_PHYSICS_HANDLER_BOUNCE_GRAVITY
	PARTICLE_PHYSICS_HANDLER_GRAVITY
	PARTICLE_PHYSICS_HANDLER_DRIFT
	PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT
	PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT_CALLBACK

	NUM_PARTICLE_PHYSICS_HANDLERS // DO NOT SET THIS MANUALLY
)
