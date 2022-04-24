package apu

type Pulse struct {
	enabled         bool
	channel         byte
	lengthEnabled   bool
	lengthValue     byte
	timerPeriod     uint16
	timerValue      uint16
	dutyMode        byte
	dutyValue       byte
	sweepReload     bool
	sweepEnabled    bool
	sweepNegate     bool
	sweepShift      byte
	sweepPeriod     byte
	sweepValue      byte
	envelopeEnabled bool
	envelopeLoop    bool
	envelopeStart   bool
	envelopePeriod  byte
	envelopeValue   byte
	envelopeVolume  byte
	constantVolume  byte
}
