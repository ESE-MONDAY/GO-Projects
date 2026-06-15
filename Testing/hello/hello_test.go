package hello

import (
	"testing"
)

func TestHello(t *testing.T) {

	t.Run("Saying hello to People", func(t *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris!"

		assertCorrectMessage(t, got, want)

	})
	t.Run("If string is empty", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World!"

		assertCorrectMessage(t, got, want)

	})
	t.Run("In spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie!"

		assertCorrectMessage(t, got, want)
	})
	t.Run("In French", func(t *testing.T) {
		got := Hello("Jean", "French")
		want := "Bonjour, Jean!"

		assertCorrectMessage(t, got, want)
	})

}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
