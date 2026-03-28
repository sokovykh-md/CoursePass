package rpc

import (
	"context"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type CoursesService struct {
	zenrpc.Service
	embedlog.Logger

	courseManager *coursepass.CourseManager
}

func NewCoursesService(dbc db.DB, logger embedlog.Logger) *CoursesService {
	return &CoursesService{
		courseManager: coursepass.NewCourseManager(dbc, logger),
		Logger:        logger,
	}
}

func (cs *CoursesService) Me(ctx context.Context) (*Student, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return nil, mapRPCError(coursepass.ErrInvalidToken)
	}

	student, err := cs.courseManager.Me(ctx, studentID)
	if err != nil {
		cs.Logger.Error(ctx, "course me failed", "err", err)
		return nil, mapRPCError(err)
	}

	return newStudent(student), nil
}

func (cs *CoursesService) List(ctx context.Context, page, pageSize int) ([]*CourseSummary, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	courses, err := cs.courseManager.Summary(ctx, page, pageSize)
	if err != nil {
		cs.Logger.Error(ctx, "course list failed", "err", err)
		return nil, mapRPCError(err)
	}

	return newCourseSummaries(courses), nil
}

func (cs *CoursesService) ByID(ctx context.Context, courseID int) (*Course, error) {
	if courseID < 1 {
		return nil, invalidParamsError("courseId", "must be greater than 0")
	}

	courseObj, err := cs.courseManager.ByID(ctx, courseID)
	if err != nil {
		cs.Logger.Error(ctx, "course by id failed", "err", err)
		return nil, mapRPCError(err)
	}

	return newCourse(courseObj), nil
}
