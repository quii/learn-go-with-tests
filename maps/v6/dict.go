package main

const (
	ErrNotFound         = DictErr("could not find the word you were looking for")
	ErrWordExists       = DictErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictErr("cannot update word because it does not exist")
)

type DictErr string

func (e DictErr) Error() string {
	return string(e)
}

func Search(dict map[string]string, word string) (string, error) {
	def, ok := dict[word]
	if !ok {
		return "", ErrNotFound
	}

	return def, nil
}

func Add(dict map[string]string, word, def string) error {
	_, err := Search(dict, word)
	switch err {
	case ErrNotFound:
		dict[word] = def
	case nil:
		return ErrWordExists
	default:
		return err

	}

	return nil
}

func Update(dict map[string]string, word, def string) error {
	_, err := Search(dict, word)
	switch err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		dict[word] = def
	default:
		return err

	}

	dict[word] = def
	return nil
}
