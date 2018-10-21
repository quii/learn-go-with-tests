package todo

import "testing"

type TodoList struct {
	items []string
}

func New() *TodoList {
	return &TodoList{}
}

func (t *TodoList) Pending() []string {
	return t.items
}

func (t *TodoList) Put(s string) {
	t.items = append(t.items, s)
}

func (t *TodoList) MarkAsDone(s string) {
	var newList []string
	for _, item := range t.items {
		if item != s {
			newList = append(newList, item)
		}
	}
	t.items = newList
}

func TestToDo(t *testing.T) {
	t.Run("empty todo returns empty", func(t *testing.T) {
		todoList := New()
		todos := todoList.Pending()

		if len(todos) != 0 {
			t.Error("expected todos to be empty")
		}
	})

	t.Run("returns pending item", func(t *testing.T) {
		todoList := New()
		item := "stroke the cat"
		todoList.Put(item)
		todos := todoList.Pending()

		assertTodoLength(t, todos, 1)
		assertFirstTodoEquaL(t, todos, item)
	})

	t.Run("mark as done", func(t *testing.T) {
		todoList := New()
		item := "stroke the cat"
		todoList.Put(item)
		todoList.MarkAsDone(item)

		assertTodoLength(t, todoList.Pending(), 0)
	})
}

func assertTodoLength(t *testing.T, list []string, want int) {
	t.Helper()
	got := len(list)
	if got != want {
		t.Errorf("expected list of size %d, got %d", want, got)
	}
}

func assertFirstTodoEquaL(t *testing.T, todos []string, item string) {
	t.Helper()
	if todos[0] != item {
		t.Errorf("want %s, got %s", item, todos[0])
	}
}
