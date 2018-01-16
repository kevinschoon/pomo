package main

type Store struct{}

func NewStore(path string) (*Store, error) {
	return &Store{}, nil
}

func (s Store) AddTask(task Task) (int, error) { return 0, nil }

func (s Store) AddPomo(pomo Pomo) error { return nil }
