package server_solution

import "time"

const FieldWidth int = 30
const FieldHeight int = 10
const MaxTravelers int = FieldHeight*FieldWidth - (FieldWidth+FieldHeight)*3

const EmptyId string = "x"
const SquatterId string = "*"
const DangerId string = "#"
const VerticalPath string = "|"
const HorizontalPath string = "-"

const BoardRefreshRate time.Duration = time.Millisecond * 2000
const NewTravelerCooldown time.Duration = time.Millisecond * 6000
const MovementCooldownInf int = 3000 // in ms
const MovementCooldownSup int = 4000 // in ms
const NewSquatterCooldown time.Duration = time.Millisecond * 5 * 1000
const NewDangerCooldown time.Duration = time.Millisecond * 10 * 1000
const SquatterLiveTime time.Duration = time.Millisecond * 12 * 1000
const DangerLiveTime time.Duration = time.Millisecond * 17 * 1000

const AreCommentLabels bool = true
