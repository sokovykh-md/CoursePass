package coursepass

import (
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"courses/pkg/db"
)

const dateTimeLayout = "2006-01-02 15:04:05"

func newDBStudent(login, passwordHash, firstName, lastName, email string) *db.Student {
	return &db.Student{
		Login:        login,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		StatusID:     db.StatusEnabled,
	}
}

func newStudent(student *db.Student) *Student {
	if student == nil {
		return nil
	}
	s := Student(*student)
	return &s
}

func newCourses(courses []db.Course) []Course {
	result := make([]Course, len(courses))
	for i := range courses {
		result[i] = Course(courses[i])
	}
	return result
}

func newStudentAuth(student db.Student) studentAuth {
	return studentAuth{
		StudentID:    student.ID,
		Login:        student.Login,
		PasswordHash: student.PasswordHash,
	}
}

func newAuthToken(token string, expiresIn int) AuthToken {
	return AuthToken{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   bearerTokenType,
	}
}

func newTokenHeader() tokenHeader {
	return tokenHeader{
		Alg: jwtAlgHS256,
		Typ: jwtTyp,
	}
}

func newTokenClaims(studentID int, login string, iat, exp int64) tokenClaims {
	return tokenClaims{
		Sub:   strconv.Itoa(studentID),
		Login: login,
		Exp:   exp,
		Iat:   iat,
	}
}

func newCourse(course *db.Course) *Course {
	if course == nil {
		return nil
	}
	c := Course(*course)
	return &c
}

func formatTimePtr(v *time.Time) *string {
	if v == nil {
		return nil
	}
	s := v.Format(dateTimeLayout)
	return &s
}

func newExamStart(exam db.Exam, questionIDs []int) ExamStart {
	return ExamStart{
		ExamID:      exam.ID,
		QuestionIDs: questionIDs,
		StartedAt:   exam.CreatedAt.Format(dateTimeLayout),
		FinishedAt:  formatTimePtr(exam.FinishedAt),
	}
}

func newQuestion(question db.Question, mediaWebPath string) Question {
	return Question{
		QuestionID:   question.ID,
		QuestionText: question.QuestionText,
		QuestionType: question.QuestionType,
		PhotoURL:     newQuestionPhotoURL(question.PhotoFile, mediaWebPath),
		Options:      newQuestionOptions(question.Options),
	}
}

func newQuestions(questions []db.Question, mediaWebPath string) []Question {
	result := make([]Question, len(questions))
	for i := range questions {
		result[i] = newQuestion(questions[i], mediaWebPath)
	}

	return result
}

func newQuestionOption(option db.QuestionOption) QuestionOption {
	return QuestionOption{
		OptionID:   option.OptionID,
		OptionText: option.OptionText,
		IsCorrect:  option.IsCorrect,
	}
}

func newQuestionOptions(options db.QuestionOptions) []QuestionOption {
	return Map(options, newQuestionOption)
}

func newQuestionPhotoURL(photoFile *db.VfsFile, mediaWebPath string) *string {
	if photoFile == nil || photoFile.Path == "" {
		return nil
	}

	basePath := strings.TrimSpace(mediaWebPath)
	if basePath == "" {
		url := photoFile.Path
		return &url
	}

	url := path.Join(basePath, strings.TrimPrefix(photoFile.Path, "/"))
	return &url
}

func newExamSummary(exam db.Exam) ExamSummary {
	finalScore := 0
	if exam.FinalScore != nil {
		finalScore = int(*exam.FinalScore)
	}

	finishedAt := ""
	if exam.FinishedAt != nil {
		finishedAt = exam.FinishedAt.Format(dateTimeLayout)
	}

	return ExamSummary{
		ExamID:     exam.ID,
		CourseID:   exam.CourseID,
		Status:     exam.Status,
		FinalScore: finalScore,
		FinishedAt: finishedAt,
	}
}

func newExamSummaries(exams []db.Exam) []ExamSummary {
	return Map(exams, newExamSummary)
}

func newExamResult(examID int, status string, finalScore, correctAnswers, totalQuestions int) ExamResult {
	return ExamResult{
		ExamID:         examID,
		Status:         status,
		FinalScore:     finalScore,
		CorrectAnswers: correctAnswers,
		TotalQuestions: totalQuestions,
	}
}

func newExamState(exam db.Exam) ExamState {
	return ExamState{
		ExamID:      exam.ID,
		CourseID:    exam.CourseID,
		Status:      exam.Status,
		QuestionIDs: slices.Clone(exam.QuestionIDs),
		Answers:     newExamStateAnswers(exam.Answers),
	}
}

func newExamStateAnswers(answers db.ExamAnswers) []ExamAnswer {
	result := make([]ExamAnswer, len(answers))
	for i := range answers {
		result[i] = ExamAnswer{
			QuestionID: answers[i].QuestionID,
			OptionIDs:  slices.Clone(answers[i].OptionIDs),
		}
	}

	return result
}

func newDBExamStateAnswers(answers []ExamAnswer) db.ExamAnswers {
	result := make(db.ExamAnswers, len(answers))
	for i := range answers {
		result[i] = db.ExamAnswer{
			QuestionID: answers[i].QuestionID,
			OptionIDs:  slices.Clone(answers[i].OptionIDs),
		}
	}

	return result
}

func newDBExamAnswersUpdate(examID int, answers []ExamAnswer) *db.Exam {
	return &db.Exam{
		ID:      examID,
		Answers: newDBExamStateAnswers(answers),
	}
}

func newDBExamSubmitUpdate(examID int, status string, correctAnswers, totalQuestions int, finalScore float64, finishedAt time.Time) *db.Exam {
	return &db.Exam{
		ID:             examID,
		Status:         status,
		CorrectAnswers: &correctAnswers,
		TotalQuestions: &totalQuestions,
		FinalScore:     &finalScore,
		FinishedAt:     &finishedAt,
	}
}
