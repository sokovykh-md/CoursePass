package rpc

import (
	"courses/pkg/coursepass"
)

func newToken(token coursepass.AuthToken) *Token {
	return &Token{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newStudent(student *coursepass.Student) *Student {
	return &Student{
		StudentID: student.StudentID,
		Login:     student.Login,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
	}
}

func newCourse(c coursepass.Course) *Course {
	return &Course{
		CourseID:      c.CourseID,
		Title:         c.Title,
		Description:   c.Description,
		TimeLimit:     c.TimeLimit,
		AvailableType: c.AvailableType,
		AvailableFrom: c.AvailableFrom,
		AvailableTo:   c.AvailableTo,
	}
}

func newCourseSummary(course coursepass.CourseSummary) *CourseSummary {
	return &CourseSummary{
		CourseID:      course.CourseID,
		Title:         course.Title,
		TimeLimit:     course.TimeLimit,
		AvailableType: course.AvailableType,
		AvailableFrom: course.AvailableFrom,
		AvailableTo:   course.AvailableTo,
	}
}

func newExamStart(start coursepass.ExamStart) *ExamStart {
	return &ExamStart{
		ExamID:      start.ExamID,
		QuestionIDs: start.QuestionIDs,
		StartedAt:   start.StartedAt,
		FinishedAt:  start.FinishedAt,
	}
}

func newQuestion(question coursepass.Question) *Question {
	return &Question{
		QuestionID:   question.QuestionID,
		QuestionText: question.QuestionText,
		QuestionType: question.QuestionType,
		PhotoURL:     question.PhotoURL,
		Options:      NewQuestionOptions(question.Options),
	}
}

func NewQuestionOption(option coursepass.QuestionOption) *QuestionOption {
	return &QuestionOption{
		OptionID:   option.OptionID,
		OptionText: option.OptionText,
	}
}

func newExamResult(result coursepass.ExamResult) *ExamResult {
	return &ExamResult{
		ExamID:         result.ExamID,
		Status:         result.Status,
		FinalScore:     result.FinalScore,
		CorrectAnswers: result.CorrectAnswers,
		TotalQuestions: result.TotalQuestions,
	}
}

func newExamSummary(summary coursepass.ExamSummary) *ExamSummary {
	return &ExamSummary{
		ExamID:     summary.ExamID,
		CourseID:   summary.CourseID,
		Status:     summary.Status,
		FinalScore: summary.FinalScore,
		FinishedAt: summary.FinishedAt,
	}
}
