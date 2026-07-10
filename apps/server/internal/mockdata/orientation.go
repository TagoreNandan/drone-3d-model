package mockdata

import (
	"math"
	"math/rand"
	"time"
)

// Orientation holds the telemetry state values.
type Orientation struct {
	Roll       float64 `json:"roll"`
	Pitch      float64 `json:"pitch"`
	Yaw        float64 `json:"yaw"`
	AltitudeM  float64 `json:"altitude_m"`
	SpeedMps   float64 `json:"speed_mps"`
	BatteryPct float64 `json:"battery_pct"`

	// Body-frame angular velocity rate metrics, corresponding to real drone telemetry (e.g. MAVLink ATTITUDE.rollspeed/pitchspeed/yawspeed).
	// Measured in radians per second (rad/s) to resemble actual IMU/flight controller outputs.
	RollSpeed  float64 `json:"rollspeed"`
	PitchSpeed float64 `json:"pitchspeed"`
	YawSpeed   float64 `json:"yawspeed"`
}

// Config defines the limits and walk steps for telemetry.
type Config struct {
	MinRoll  float64
	MaxRoll  float64
	MinPitch float64
	MaxPitch float64
	MinAlt   float64
	MaxAlt   float64
	MinSpeed float64
	MaxSpeed float64

	RollStep  float64
	PitchStep float64
	YawStep   float64
	AltStep   float64
	SpeedStep float64

	BatteryDrainPerSec float64
	BatteryNoisePerSec float64
}

// DefaultConfig provides defaults matching the issue specs.
func DefaultConfig() Config {
	return Config{
		MinRoll:            -25.0,
		MaxRoll:            25.0,
		MinPitch:           -25.0,
		MaxPitch:           25.0,
		MinAlt:             0.0,
		MaxAlt:             120.0,
		MinSpeed:           0.0,
		MaxSpeed:           20.0,
		RollStep:           5.0,
		PitchStep:          5.0,
		YawStep:            10.0,
		AltStep:            2.0,
		SpeedStep:          1.5,
		BatteryDrainPerSec: 0.05,
		BatteryNoisePerSec: 0.005,
	}
}

// Generator runs a simple bounded random walk.
type Generator struct {
	rng    *rand.Rand
	config Config

	roll       float64
	pitch      float64
	yaw        float64
	altitudeM  float64
	speedMps   float64
	batteryPct float64

	// Internal angular speeds in degrees/second for dynamics integration
	rollSpeed  float64
	pitchSpeed float64
	yawSpeed   float64
}

// NewGenerator initializes the state deterministically.
func NewGenerator(cfg Config) *Generator {
	return &Generator{
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
		config:     cfg,
		roll:       0.0,
		pitch:      0.0,
		yaw:        0.0,
		altitudeM:  (cfg.MinAlt + cfg.MaxAlt) / 2.0,
		speedMps:   (cfg.MinSpeed + cfg.MaxSpeed) / 2.0,
		batteryPct: 100.0,
		rollSpeed:  0.0,
		pitchSpeed: 0.0,
		yawSpeed:   0.0,
	}
}

// Next moves the random walk forward by deltaSeconds.
func (g *Generator) Next(deltaSeconds float64) Orientation {
	if deltaSeconds <= 0 {
		return Orientation{
			Roll:       g.roll,
			Pitch:      g.pitch,
			Yaw:        g.yaw,
			AltitudeM:  g.altitudeM,
			SpeedMps:   g.speedMps,
			BatteryPct: g.batteryPct,
			RollSpeed:  g.rollSpeed * (math.Pi / 180.0),
			PitchSpeed: g.pitchSpeed * (math.Pi / 180.0),
			YawSpeed:   g.yawSpeed * (math.Pi / 180.0),
		}
	}

	// Sub-step integration to maintain numerical stability for the spring-mass-damper system
	dtRemaining := deltaSeconds
	const maxStep = 0.05 // 50ms sub-steps for stability
	for dtRemaining > 0 {
		dt := dtRemaining
		if dt > maxStep {
			dt = maxStep
		}
		dtRemaining -= dt

		// Roll dynamics (stabilized towards 0 degrees using spring-mass-damper model)
		// Accelerations: random turbulence + spring centering force + drag/damping
		rollAccel := (g.rng.Float64()*2 - 1) * 150.0 - 8.0 * g.roll - 4.0 * g.rollSpeed
		g.rollSpeed += rollAccel * dt
		// Clamp speed to realistic max roll rate of 100 deg/s
		if g.rollSpeed < -100.0 {
			g.rollSpeed = -100.0
		} else if g.rollSpeed > 100.0 {
			g.rollSpeed = 100.0
		}
		g.roll += g.rollSpeed * dt

		// Pitch dynamics (stabilized towards 0 degrees using spring-mass-damper model)
		pitchAccel := (g.rng.Float64()*2 - 1) * 150.0 - 8.0 * g.pitch - 4.0 * g.pitchSpeed
		g.pitchSpeed += pitchAccel * dt
		if g.pitchSpeed < -100.0 {
			g.pitchSpeed = -100.0
		} else if g.pitchSpeed > 100.0 {
			g.pitchSpeed = 100.0
		}
		g.pitch += g.pitchSpeed * dt

		// Yaw dynamics (drifting yaw-rate, no centering spring force since heading is free)
		yawAccel := (g.rng.Float64()*2 - 1) * 40.0 - 2.0 * g.yawSpeed
		g.yawSpeed += yawAccel * dt
		if g.yawSpeed < -45.0 {
			g.yawSpeed = -45.0
		} else if g.yawSpeed > 45.0 {
			g.yawSpeed = 45.0
		}
		g.yaw += g.yawSpeed * dt
	}

	// Update altitude and speed with standard random walk
	g.altitudeM += (g.rng.Float64()*2 - 1) * g.config.AltStep * deltaSeconds
	g.speedMps += (g.rng.Float64()*2 - 1) * g.config.SpeedStep * deltaSeconds

	// Battery slowly drains with small random variation
	batteryDelta := -(g.config.BatteryDrainPerSec * deltaSeconds) + (g.rng.Float64()*2-1)*g.config.BatteryNoisePerSec*deltaSeconds
	g.batteryPct += batteryDelta

	// Clamp states to configured ranges
	if g.roll < g.config.MinRoll {
		g.roll = g.config.MinRoll
		g.rollSpeed = 0
	} else if g.roll > g.config.MaxRoll {
		g.roll = g.config.MaxRoll
		g.rollSpeed = 0
	}

	if g.pitch < g.config.MinPitch {
		g.pitch = g.config.MinPitch
		g.pitchSpeed = 0
	} else if g.pitch > g.config.MaxPitch {
		g.pitch = g.config.MaxPitch
		g.pitchSpeed = 0
	}

	if g.altitudeM < g.config.MinAlt {
		g.altitudeM = g.config.MinAlt
	} else if g.altitudeM > g.config.MaxAlt {
		g.altitudeM = g.config.MaxAlt
	}

	if g.speedMps < g.config.MinSpeed {
		g.speedMps = g.config.MinSpeed
	} else if g.speedMps > g.config.MaxSpeed {
		g.speedMps = g.config.MaxSpeed
	}

	if g.batteryPct < 0 {
		g.batteryPct = 0
	} else if g.batteryPct > 100 {
		g.batteryPct = 100
	}

	// Wrap yaw continuously between 0 and 360 degrees
	g.yaw = math.Mod(g.yaw, 360.0)
	if g.yaw < 0 {
		g.yaw += 360.0
	}

	return Orientation{
		Roll:       g.roll,
		Pitch:      g.pitch,
		Yaw:        g.yaw,
		AltitudeM:  g.altitudeM,
		SpeedMps:   g.speedMps,
		BatteryPct: g.batteryPct,
		RollSpeed:  g.rollSpeed * (math.Pi / 180.0),
		PitchSpeed: g.pitchSpeed * (math.Pi / 180.0),
		YawSpeed:   g.yawSpeed * (math.Pi / 180.0),
	}
}
