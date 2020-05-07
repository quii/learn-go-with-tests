package http_handlers_revisited

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type MockUserService struct {
	AddFunc func(user User) (string, error)
	UsersAdded []User
}

func (m *MockUserService) Add(user User) (insertedID string, err error) {
	m.UsersAdded = append(m.UsersAdded, user)
	return m.AddFunc(user)
}

func TestRegisterUser(t *testing.T) {
	t.Run("can add valid users", func(t *testing.T) {
		user := User{Name: "CJ"}
		expectedInsertedID := "whatever"

		service := &MockUserService{
			AddFunc: func(user User) (string, error) {
				return expectedInsertedID, nil
			},
		}
		server := NewUserServer(service)

		req := httptest.NewRequest(http.MethodGet, "/", userToJSON(user))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusCreated)

		if res.Body.String() != expectedInsertedID {
			t.Errorf("expected body of %q but got %q", res.Body.String(), expectedInsertedID)
		}

		if len(service.UsersAdded)!= 1 {
			t.Fatalf("expected 1 user added but got %d", len(service.UsersAdded))
		}

		if !reflect.DeepEqual(service.UsersAdded[0], user) {
			t.Errorf("the user added %+v was not what was expected %+v", service.UsersAdded[0], user)
		}
	})

	t.Run("returns 400 bad request if body is not valid user JSON", func(t *testing.T) {
		server := NewUserServer(nil)

		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("trouble will find me"))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusBadRequest)
	})

	t.Run("returns a 500 internal server error if the service fails", func(t *testing.T) {
		user := User{Name: "CJ"}

		service := &MockUserService{
			AddFunc: func(user User) (string, error) {
				return "", errors.New("couldn't add new user")
			},
		}
		server := NewUserServer(service)

		req := httptest.NewRequest(http.MethodGet, "/", userToJSON(user))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusInternalServerError)
	})
}

func assertStatus(t *testing.T, got, want int){
	t.Helper()
	if got != want {
		t.Errorf("wanted http status %d but got %d", got, want)
	}
}

func userToJSON(user User) io.Reader {
	stuff, _ := json.Marshal(user)
	return bytes.NewReader(stuff)
}
