package store

type UserStore struct {
	*Store
}

func NewUserStore(store *Store) *UserStore {
	return &UserStore{store}
}
