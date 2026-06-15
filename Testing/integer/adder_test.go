package integer

import "testing"

func TestAdder(t *testing.T) {
	t.Run("Adding two positive numbers", func(t *testing.T) {
		got := Adder(2, 3)
		want := 5

		assertCorrectMessage(t, got, want)

	})
	t.Run("Adding a positive and a negative number", func(t *testing.T) {
		got := Adder(5, -2)
		want := 3

		assertCorrectMessage(t, got, want)

	})
	t.Run("Adding two negative numbers", func(t *testing.T) {
		got := Adder(-4, -6)
		want := -10

		assertCorrectMessage(t, got, want)

	})
}

func assertCorrectMessage(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
