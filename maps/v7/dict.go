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

type Dict map[string]string

func (d Dict) Search(word string) (string, error) {
	def, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}

	return def, nil
}

func (d Dict) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		d[word] = def
	case nil:
		return ErrWordExists
	default:
		return err

	}

	return nil
}

func (d Dict) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		d[word] = def
	default:
		return err

	}

	return nil
}

func (d Dict) Delete(word string) {
	delete(d, word)
}
