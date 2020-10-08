package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		assert.Equal(t, "Hello, Chris", Hello("Chris", ""))
	})

	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
		assert.Equal(t, "Hello, World", Hello("", ""))
	})

	t.Run("say hello in Spanish", func(t *testing.T) {
		assert.Equal(t, "Hola, Elodie", Hello("Elodie", spanish))
	})

	t.Run("say hello in French", func(t *testing.T) {
		assert.Equal(t, "Bonjour, Lauren", Hello("Lauren", french))
	})

}
