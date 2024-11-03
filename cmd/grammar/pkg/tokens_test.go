package pkg

import "testing"

func TestFixUnderscores(t *testing.T) {
	const (
		Input    string = "_hello_world"
		Expected string = "HelloWorld"
	)

	str, err := FixUnderscores([]byte(Input), true)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	} else if str != Expected {
		t.Errorf("want %q, got %q", Expected, str)
	}
}
