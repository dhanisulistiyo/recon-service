package reconcile

type reconcileService struct {
}

func NewService() Service {
	return &reconcileService{}
}

type Service interface {
	Execute()
}
