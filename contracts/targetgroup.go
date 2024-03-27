package contracts

type TargetGroup interface {
	Name() string
	GetNextAvailableTarget() Target
}
