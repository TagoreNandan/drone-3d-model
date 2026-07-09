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
		}
	}

	// Random step scaled by delta time
	g.roll += (g.rng.Float64()*2 - 1) * g.config.RollStep * deltaSeconds
	g.pitch += (g.rng.Float64()*2 - 1) * g.config.PitchStep * deltaSeconds
	g.yaw += (g.rng.Float64()*2 - 1) * g.config.YawStep * deltaSeconds
	g.altitudeM += (g.rng.Float64()*2 - 1) * g.config.AltStep * deltaSeconds
	g.speedMps += (g.rng.Float64()*2 - 1) * g.config.SpeedStep * deltaSeconds

	// Battery slowly drains with small random variation
	batteryDelta := -(g.config.BatteryDrainPerSec * deltaSeconds) + (g.rng.Float64()*2-1)*g.config.BatteryNoisePerSec*deltaSeconds
	g.batteryPct += batteryDelta

	// Clamp states to configured ranges
	if g.roll < g.config.MinRoll {
		g.roll = g.config.MinRoll
	} else if g.roll > g.config.MaxRoll {
		g.roll = g.config.MaxRoll
	}

	if g.pitch < g.config.MinPitch {
		g.pitch = g.config.MinPitch
	} else if g.pitch > g.config.MaxPitch {
		g.pitch = g.config.MaxPitch
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
	}
}
