export interface CssZone {
  key: string
  name: string
  focus: string
  minSpeedPercent: number
  maxSpeedPercent?: number
  slowerPaceSeconds: number
  fasterPaceSeconds: number
}

const ZONE_DEFINITIONS = [
  { key: 'z1', name: 'Recovery', focus: 'Warm-up, recovery, and technique', min: 77, max: 87 },
  {
    key: 'z2',
    name: 'Aerobic endurance',
    focus: 'Long, steady swims and base fitness',
    min: 87,
    max: 94,
  },
  { key: 'z3', name: 'Threshold', focus: 'Improve sustainable pace', min: 95, max: 104 },
  {
    key: 'z4',
    name: 'Anaerobic / VO2 Max',
    focus: 'Speed and lactate tolerance',
    min: 104,
    max: 111,
  },
  { key: 'z5', name: 'Sprint', focus: 'Maximum speed and power', min: 111 },
]

export function parseSwimTime(value: string): number | null {
  const match = value.trim().match(/^(\d+):([0-5]\d)$/)
  if (!match) return null

  return Number(match[1]) * 60 + Number(match[2])
}

export function formatSwimTime(seconds: number | null | undefined): string {
  if (seconds == null || !Number.isFinite(seconds) || seconds < 0) return ''

  const roundedSeconds = Math.round(seconds)
  const minutes = Math.floor(roundedSeconds / 60)
  const remainingSeconds = roundedSeconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

export function calculateCssPace(
  time200m: string | number | null | undefined,
  time400m: string | number | null | undefined,
): number | null {
  const seconds200 = typeof time200m === 'string' ? parseSwimTime(time200m) : time200m
  const seconds400 = typeof time400m === 'string' ? parseSwimTime(time400m) : time400m

  if (
    seconds200 == null ||
    seconds400 == null ||
    !Number.isFinite(seconds200) ||
    !Number.isFinite(seconds400) ||
    seconds200 <= 0 ||
    seconds400 <= seconds200
  ) {
    return null
  }

  return (seconds400 - seconds200) / 2
}

export function calculateCssZones(cssPaceSeconds: number | null): CssZone[] {
  if (cssPaceSeconds == null || !Number.isFinite(cssPaceSeconds) || cssPaceSeconds <= 0) {
    return []
  }

  return ZONE_DEFINITIONS.map((zone) => ({
    key: zone.key,
    name: zone.name,
    focus: zone.focus,
    minSpeedPercent: zone.min,
    maxSpeedPercent: zone.max,
    slowerPaceSeconds: cssPaceSeconds / (zone.min / 100),
    fasterPaceSeconds: cssPaceSeconds / ((zone.max ?? zone.min) / 100),
  }))
}

export function formatPaceRange(zone: CssZone): string {
  if (zone.key === 'z5') {
    return `< ${formatSwimTime(zone.fasterPaceSeconds)} / 100 m`
  }

  return `${formatSwimTime(zone.fasterPaceSeconds)}-${formatSwimTime(zone.slowerPaceSeconds)} / 100 m`
}
