package coursepass

// colgen:Course:MapP(db.Course)
// colgen:Exam:MapP(db.Exam)
// colgen:Question:MapP(db.Question)
//
//go:generate colgen -imports courses/pkg/db
func MapP[T, M any](in []T, convert func(*T) *M) []M {
	out := make([]M, len(in))
	for i := range in {
		out[i] = *convert(&in[i])
	}

	return out
}
