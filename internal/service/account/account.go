package account

const (
	DefaultMaxRows = 100
)

type Service struct {
	store storage
	log   logger
}

func New(store storage, log logger) *Service {
	return &Service{
		store: store,
		log:   log,
	}
}
