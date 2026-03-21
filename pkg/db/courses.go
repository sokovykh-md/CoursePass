package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type CoursesRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewCoursesRepo returns new repository
func NewCoursesRepo(db orm.DB) CoursesRepo {
	return CoursesRepo{
		db: db,
		filters: map[string][]Filter{
			Tables.User.Name:    {StatusFilter},
			Tables.Course.Name:  {StatusFilter},
			Tables.Student.Name: {StatusFilter},
		},
		sort: map[string][]SortField{
			Tables.User.Name:     {{Column: Columns.User.CreatedAt, Direction: SortDesc}},
			Tables.Course.Name:   {{Column: Columns.Course.CreatedAt, Direction: SortDesc}},
			Tables.Exam.Name:     {{Column: Columns.Exam.CreatedAt, Direction: SortDesc}},
			Tables.Question.Name: {{Column: Columns.Question.CreatedAt, Direction: SortDesc}},
			Tables.Student.Name:  {{Column: Columns.Student.CreatedAt, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.User.Name:     {TableColumns},
			Tables.Course.Name:   {TableColumns},
			Tables.Exam.Name:     {TableColumns, Columns.Exam.Course, Columns.Exam.Student},
			Tables.Question.Name: {TableColumns, Columns.Question.Course, Columns.Question.PhotoFile},
			Tables.Student.Name:  {TableColumns},
		},
	}
}

// WithTransaction is a function that wraps CoursesRepo with pg.Tx transaction.
func (cr CoursesRepo) WithTransaction(tx *pg.Tx) CoursesRepo {
	cr.db = tx
	return cr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (cr CoursesRepo) WithEnabledOnly() CoursesRepo {
	f := make(map[string][]Filter, len(cr.filters))
	for i := range cr.filters {
		f[i] = make([]Filter, len(cr.filters[i]))
		copy(f[i], cr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	cr.filters = f

	return cr
}

/*** User ***/

// FullUser returns full joins with all columns
func (cr CoursesRepo) FullUser() OpFunc {
	return WithColumns(cr.join[Tables.User.Name]...)
}

// DefaultUserSort returns default sort.
func (cr CoursesRepo) DefaultUserSort() OpFunc {
	return WithSort(cr.sort[Tables.User.Name]...)
}

// UserByID is a function that returns User by ID(s) or nil.
func (cr CoursesRepo) UserByID(ctx context.Context, id int, ops ...OpFunc) (*User, error) {
	return cr.OneUser(ctx, &UserSearch{ID: &id}, ops...)
}

// OneUser is a function that returns one User by filters. It could return pg.ErrMultiRows.
func (cr CoursesRepo) OneUser(ctx context.Context, search *UserSearch, ops ...OpFunc) (*User, error) {
	obj := &User{}
	err := buildQuery(ctx, cr.db, obj, search, cr.filters[Tables.User.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// UsersByFilters returns User list.
func (cr CoursesRepo) UsersByFilters(ctx context.Context, search *UserSearch, pager Pager, ops ...OpFunc) (users []User, err error) {
	err = buildQuery(ctx, cr.db, &users, search, cr.filters[Tables.User.Name], pager, ops...).Select()
	return
}

// CountUsers returns count
func (cr CoursesRepo) CountUsers(ctx context.Context, search *UserSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, cr.db, &User{}, search, cr.filters[Tables.User.Name], PagerOne, ops...).Count()
}

// AddUser adds User to DB.
func (cr CoursesRepo) AddUser(ctx context.Context, user *User, ops ...OpFunc) (*User, error) {
	q := cr.db.ModelContext(ctx, user)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.User.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return user, err
}

// UpdateUser updates User in DB.
func (cr CoursesRepo) UpdateUser(ctx context.Context, user *User, ops ...OpFunc) (bool, error) {
	q := cr.db.ModelContext(ctx, user).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.User.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteUser set statusId to deleted in DB.
func (cr CoursesRepo) DeleteUser(ctx context.Context, id int) (deleted bool, err error) {
	user := &User{ID: id, StatusID: StatusDeleted}

	return cr.UpdateUser(ctx, user, WithColumns(Columns.User.StatusID))
}

/*** Course ***/

// FullCourse returns full joins with all columns
func (cr CoursesRepo) FullCourse() OpFunc {
	return WithColumns(cr.join[Tables.Course.Name]...)
}

// DefaultCourseSort returns default sort.
func (cr CoursesRepo) DefaultCourseSort() OpFunc {
	return WithSort(cr.sort[Tables.Course.Name]...)
}

// CourseByID is a function that returns Course by ID(s) or nil.
func (cr CoursesRepo) CourseByID(ctx context.Context, id int, ops ...OpFunc) (*Course, error) {
	return cr.OneCourse(ctx, &CourseSearch{ID: &id}, ops...)
}

// OneCourse is a function that returns one Course by filters. It could return pg.ErrMultiRows.
func (cr CoursesRepo) OneCourse(ctx context.Context, search *CourseSearch, ops ...OpFunc) (*Course, error) {
	obj := &Course{}
	err := buildQuery(ctx, cr.db, obj, search, cr.filters[Tables.Course.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// CoursesByFilters returns Course list.
func (cr CoursesRepo) CoursesByFilters(ctx context.Context, search *CourseSearch, pager Pager, ops ...OpFunc) (courses []Course, err error) {
	err = buildQuery(ctx, cr.db, &courses, search, cr.filters[Tables.Course.Name], pager, ops...).Select()
	return
}

// CountCourses returns count
func (cr CoursesRepo) CountCourses(ctx context.Context, search *CourseSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, cr.db, &Course{}, search, cr.filters[Tables.Course.Name], PagerOne, ops...).Count()
}

// AddCourse adds Course to DB.
func (cr CoursesRepo) AddCourse(ctx context.Context, course *Course, ops ...OpFunc) (*Course, error) {
	q := cr.db.ModelContext(ctx, course)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Course.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return course, err
}

// UpdateCourse updates Course in DB.
func (cr CoursesRepo) UpdateCourse(ctx context.Context, course *Course, ops ...OpFunc) (bool, error) {
	q := cr.db.ModelContext(ctx, course).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Course.ID, Columns.Course.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteCourse set statusId to deleted in DB.
func (cr CoursesRepo) DeleteCourse(ctx context.Context, id int) (deleted bool, err error) {
	course := &Course{ID: id, StatusID: StatusDeleted}

	return cr.UpdateCourse(ctx, course, WithColumns(Columns.Course.StatusID))
}

/*** Exam ***/

// FullExam returns full joins with all columns
func (cr CoursesRepo) FullExam() OpFunc {
	return WithColumns(cr.join[Tables.Exam.Name]...)
}

// DefaultExamSort returns default sort.
func (cr CoursesRepo) DefaultExamSort() OpFunc {
	return WithSort(cr.sort[Tables.Exam.Name]...)
}

// ExamByID is a function that returns Exam by ID(s) or nil.
func (cr CoursesRepo) ExamByID(ctx context.Context, id int, ops ...OpFunc) (*Exam, error) {
	return cr.OneExam(ctx, &ExamSearch{ID: &id}, ops...)
}

// OneExam is a function that returns one Exam by filters. It could return pg.ErrMultiRows.
func (cr CoursesRepo) OneExam(ctx context.Context, search *ExamSearch, ops ...OpFunc) (*Exam, error) {
	obj := &Exam{}
	err := buildQuery(ctx, cr.db, obj, search, cr.filters[Tables.Exam.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// ExamsByFilters returns Exam list.
func (cr CoursesRepo) ExamsByFilters(ctx context.Context, search *ExamSearch, pager Pager, ops ...OpFunc) (exams []Exam, err error) {
	err = buildQuery(ctx, cr.db, &exams, search, cr.filters[Tables.Exam.Name], pager, ops...).Select()
	return
}

// CountExams returns count
func (cr CoursesRepo) CountExams(ctx context.Context, search *ExamSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, cr.db, &Exam{}, search, cr.filters[Tables.Exam.Name], PagerOne, ops...).Count()
}

// AddExam adds Exam to DB.
func (cr CoursesRepo) AddExam(ctx context.Context, exam *Exam, ops ...OpFunc) (*Exam, error) {
	q := cr.db.ModelContext(ctx, exam)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Exam.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return exam, err
}

// UpdateExam updates Exam in DB.
func (cr CoursesRepo) UpdateExam(ctx context.Context, exam *Exam, ops ...OpFunc) (bool, error) {
	q := cr.db.ModelContext(ctx, exam).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Exam.ID, Columns.Exam.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteExam deletes Exam from DB.
func (cr CoursesRepo) DeleteExam(ctx context.Context, id int) (deleted bool, err error) {
	exam := &Exam{ID: id}

	res, err := cr.db.ModelContext(ctx, exam).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

/*** Question ***/

// FullQuestion returns full joins with all columns
func (cr CoursesRepo) FullQuestion() OpFunc {
	return WithColumns(cr.join[Tables.Question.Name]...)
}

// DefaultQuestionSort returns default sort.
func (cr CoursesRepo) DefaultQuestionSort() OpFunc {
	return WithSort(cr.sort[Tables.Question.Name]...)
}

// QuestionByID is a function that returns Question by ID(s) or nil.
func (cr CoursesRepo) QuestionByID(ctx context.Context, id int, ops ...OpFunc) (*Question, error) {
	return cr.OneQuestion(ctx, &QuestionSearch{ID: &id}, ops...)
}

// OneQuestion is a function that returns one Question by filters. It could return pg.ErrMultiRows.
func (cr CoursesRepo) OneQuestion(ctx context.Context, search *QuestionSearch, ops ...OpFunc) (*Question, error) {
	obj := &Question{}
	err := buildQuery(ctx, cr.db, obj, search, cr.filters[Tables.Question.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// QuestionsByFilters returns Question list.
func (cr CoursesRepo) QuestionsByFilters(ctx context.Context, search *QuestionSearch, pager Pager, ops ...OpFunc) (questions []Question, err error) {
	err = buildQuery(ctx, cr.db, &questions, search, cr.filters[Tables.Question.Name], pager, ops...).Select()
	return
}

// CountQuestions returns count
func (cr CoursesRepo) CountQuestions(ctx context.Context, search *QuestionSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, cr.db, &Question{}, search, cr.filters[Tables.Question.Name], PagerOne, ops...).Count()
}

// AddQuestion adds Question to DB.
func (cr CoursesRepo) AddQuestion(ctx context.Context, question *Question, ops ...OpFunc) (*Question, error) {
	q := cr.db.ModelContext(ctx, question)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Question.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return question, err
}

// UpdateQuestion updates Question in DB.
func (cr CoursesRepo) UpdateQuestion(ctx context.Context, question *Question, ops ...OpFunc) (bool, error) {
	q := cr.db.ModelContext(ctx, question).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Question.ID, Columns.Question.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteQuestion deletes Question from DB.
func (cr CoursesRepo) DeleteQuestion(ctx context.Context, id int) (deleted bool, err error) {
	question := &Question{ID: id}

	res, err := cr.db.ModelContext(ctx, question).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

/*** Student ***/

// FullStudent returns full joins with all columns
func (cr CoursesRepo) FullStudent() OpFunc {
	return WithColumns(cr.join[Tables.Student.Name]...)
}

// DefaultStudentSort returns default sort.
func (cr CoursesRepo) DefaultStudentSort() OpFunc {
	return WithSort(cr.sort[Tables.Student.Name]...)
}

// StudentByID is a function that returns Student by ID(s) or nil.
func (cr CoursesRepo) StudentByID(ctx context.Context, id int, ops ...OpFunc) (*Student, error) {
	return cr.OneStudent(ctx, &StudentSearch{ID: &id}, ops...)
}

// OneStudent is a function that returns one Student by filters. It could return pg.ErrMultiRows.
func (cr CoursesRepo) OneStudent(ctx context.Context, search *StudentSearch, ops ...OpFunc) (*Student, error) {
	obj := &Student{}
	err := buildQuery(ctx, cr.db, obj, search, cr.filters[Tables.Student.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// StudentsByFilters returns Student list.
func (cr CoursesRepo) StudentsByFilters(ctx context.Context, search *StudentSearch, pager Pager, ops ...OpFunc) (students []Student, err error) {
	err = buildQuery(ctx, cr.db, &students, search, cr.filters[Tables.Student.Name], pager, ops...).Select()
	return
}

// CountStudents returns count
func (cr CoursesRepo) CountStudents(ctx context.Context, search *StudentSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, cr.db, &Student{}, search, cr.filters[Tables.Student.Name], PagerOne, ops...).Count()
}

// AddStudent adds Student to DB.
func (cr CoursesRepo) AddStudent(ctx context.Context, student *Student, ops ...OpFunc) (*Student, error) {
	q := cr.db.ModelContext(ctx, student)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Student.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return student, err
}

// UpdateStudent updates Student in DB.
func (cr CoursesRepo) UpdateStudent(ctx context.Context, student *Student, ops ...OpFunc) (bool, error) {
	q := cr.db.ModelContext(ctx, student).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Student.ID, Columns.Student.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteStudent set statusId to deleted in DB.
func (cr CoursesRepo) DeleteStudent(ctx context.Context, id int) (deleted bool, err error) {
	student := &Student{ID: id, StatusID: StatusDeleted}

	return cr.UpdateStudent(ctx, student, WithColumns(Columns.Student.StatusID))
}
