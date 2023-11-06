package account

const (
	DefMaxRows = 100
)

type Service struct {
	store storage
}

func New(store storage) *Service {
	return &Service{
		store: store,
	}
}
