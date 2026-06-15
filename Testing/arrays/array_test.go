package arrays

import (
	"testing"

	"slices"
)

func TestSumArray(t *testing.T) {
	t.Run("Sum of an array of integers", func(t *testing.T) {
		got := SumArray([]int{1, 2, 3, 4, 5})
		want := 15

		assertCorrectMessage(t, got, want)

	})
	t.Run("Sum with no fixed length", func(t *testing.T) {
		got := SumArray([]int{10, 20, 30, 40, 50})
		want := 150
		assertCorrectMessage(t, got, want)
	})

}

func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertCorrectMessage(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
