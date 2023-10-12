package account

type Service struct {
	store storage
}

func New(store storage) *Service {
	return &Service{
		store: store,
	}
}
