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
	courses, err := cm.availableCourses(ctx, currentTime, page, pageSize)
	if err != nil {
		return nil, err
	}

	return newCourseSummaries(courses), nil
}

func (cm *CourseManager) ByID(ctx context.Context, courseID int) (Course, error) {
	courseData, err := cm.courseByID(ctx, courseID)
	if err != nil {
		return Course{}, err
	}

	return newCourse(*courseData), nil
}

func (cm *CourseManager) Me(ctx context.Context, studentID int) (*Student, error) {
	studentData, err := cm.studentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	result := newStudent(*studentData)
	return &result, nil
}

func (cm *CourseManager) availableCourses(ctx context.Context, currentTime time.Time, page, pageSize int) ([]db.Course, error) {
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

	return courses, nil
}

func (cm *CourseManager) courseByID(ctx context.Context, courseID int) (*db.Course, error) {
	courseData, err := cm.repo.OneCourse(ctx, &db.CourseSearch{
		ID: &courseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get coursepass: %w", err)
	}
	if courseData == nil {
		return nil, ErrCourseNotFound
	}

	return courseData, nil
}

func (cm *CourseManager) studentByID(ctx context.Context, studentID int) (*db.Student, error) {
	studentData, err := cm.repo.OneStudent(ctx, &db.StudentSearch{
		ID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return nil, ErrStudentNotFound
	}

	return studentData, nil
}
