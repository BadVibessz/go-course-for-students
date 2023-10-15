package ads

import "errors"

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	Published bool
}

// todo: default params?
func (a Ad) New() {

}

// todo: create module validator and publish on github then import here
func (a *Ad) Validate() (bool, error) {

	titleLen := len([]rune(a.Title))
	textLen := len([]rune(a.Text))

	if titleLen == 0 {
		return false, errors.New("title cannot be empty")
	} else if titleLen >= 100 {
		return false, errors.New("title length cannot be more than 99")
	}

	if textLen == 0 {
		return false, errors.New("text cannot be empty")
	} else if textLen >= 500 {
		return false, errors.New("text length cannot be more than 499")
	}

	return true, nil
}

func ValidateID(id int64) (bool, error) {

	if id < 0 {
		return false, errors.New("id cannot be negative value")
	}
	return true, nil
}
