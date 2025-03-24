package registry

type Action int

func (a Action) String() string {
	switch a {
	case ActionCreate:
		return "create"
	case ActionUpdate:
		return "update"
	case ActionDelete:
		return "delete"
	default:
		return "invalid"
	}
}

const (
	ActionCreate Action = iota
	ActionUpdate
	ActionDelete
)
