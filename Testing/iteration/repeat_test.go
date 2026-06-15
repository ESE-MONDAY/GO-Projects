package iteration

import "testing"

func TestRepeat(t *testing.T) {
	t.Run("Repeat a character", func(t *testing.T) {
		got := Repeat("a", 5)
		want := "aaaaa"

		assertCorrectMessage(t, got, want)

	})

}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		_ = Repeat("a", 5)
	}
}
