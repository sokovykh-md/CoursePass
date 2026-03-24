package coursepass

import (
	"context"
	"fmt"
	"time"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
)

type CourseManager struct {
	repo db.CoursesRepo
	embedlog.Logger
}

func NewCourseManager(dbo db.DB, logger embedlog.Logger) *CourseManager {
	return &CourseManager{
		repo:   db.NewCoursesRepo(dbo),
		Logger: logger,
	}
}

func (cm *CourseManager) Summary(ctx context.Context, page, pageSize int) ([]CourseSummary, error) {
	currentTime := time.Now()

	courses, err := cm.repo.CoursesByFilters(ctx, &db.CourseSearch{
		AvailableFromTo: &currentTime,
		AvailableToFrom: &currentTime,
	}, db.Pager{
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		return nil, fmt.Errorf("failed get courses: %w", err)
	}

	return newCourseSummaries(courses), nil
}

func (cm *CourseManager) ByID(ctx context.Context, courseID int) (Course, error) {
	courseData, err := cm.repo.OneCourse(ctx, &db.CourseSearch{
		ID: &courseID,
	})
	if err != nil {
		return Course{}, fmt.Errorf("failed get coursepass: %w", err)
	}
	if courseData == nil {
		return Course{}, ErrCourseNotFound
	}

	return newCourse(*courseData), nil
}

func (cm *CourseManager) Me(ctx context.Context, studentID int) (*Student, error) {
	student, err := cm.repo.OneStudent(ctx, &db.StudentSearch{
		ID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if student == nil {
		return nil, ErrStudentNotFound
	}

	result := newStudent(*student)
	return &result, nil
}
