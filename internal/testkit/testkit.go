package testkit

import "testing"

const succeed = "\u2713"
const failed = "\u2717"

// Check and ErrorT functions help style the test logs a little better.
// They prepend 'succeeded' and 'failed' symbols to the beginning of Logf and Errorf.
// Example:
// want := 2
// Check(t, x == want, "want: %d, received: %d", []any{want, x}...)
// Result:  ✓       want 2, received: 2
// In case 'ok' is false: ✗       want 2, received: 3

func Check(t *testing.T, ok bool, format string, args ...any) {
	if ok {
		args = prepend(args, succeed)
		t.Logf("\t%s\t"+format, args...)
		return
	}

	ErrorT(t, format, args...)
}

func ErrorT(t *testing.T, format string, args ...any) {
	args = prepend(args, failed)
	t.Errorf("\t%s\t"+format, args...)
}

// prepend adds an element to the beginning of slice.
func prepend[S ~[]E, E any](s S, elem E) []any {
	args := []any{elem}
	for _, e := range s {
		args = append(args, e)
	}

	return args
}
