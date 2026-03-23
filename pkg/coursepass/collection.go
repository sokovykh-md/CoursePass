package coursepass

//go:generate colgen -imports courses/pkg/db
//colgen:CourseSummary:map(db.Course)

func Map[S, T any](in []S, convert func(S) T) []T {
	out := make([]T, len(in))
	for i := range in {
		out[i] = convert(in[i])
	}

	return out
}
