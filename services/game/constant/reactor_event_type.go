package constant

import "fmt"

type ReactorEventType int

const (
	ReactorEventTypeHit            ReactorEventType = 0
	ReactorEventTypeDirectionalHit ReactorEventType = 2
	ReactorEventTypeTouch          ReactorEventType = 9
	ReactorEventTypeItem           ReactorEventType = 100
)

func (t ReactorEventType) String() string {
	switch t {
	case ReactorEventTypeHit:
		return "Hit"
	case ReactorEventTypeDirectionalHit:
		return "DirectionalHit"
	case ReactorEventTypeTouch:
		return "Touch"
	case ReactorEventTypeItem:
		return "Item"
	default:
		return fmt.Sprintf("ReactorEventType(%d)", int(t))
	}
}
