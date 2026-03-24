package vt

import (
	"net/http"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	zm "github.com/vmkteam/zenrpc-middleware"
	"github.com/vmkteam/zenrpc/v2"
)

//go:generate go tool zenrpc

const (
	AuthKey = "Authorization2"
)

const (
	NSAuth = "auth"
	NSUser = "user"

	NSCourse    = "course"
	NSExam      = "exam"
	NSQuestion  = "question"
	NSStudent   = "student"
	NSVfsFile   = "vfsFile"
	NSVfsFolder = "vfsFolder"
)

var (
	ErrUnauthorized   = httpAsRPCError(http.StatusUnauthorized)
	ErrForbidden      = httpAsRPCError(http.StatusForbidden)
	ErrNotFound       = httpAsRPCError(http.StatusNotFound)
	ErrInternal       = httpAsRPCError(http.StatusInternalServerError)
	ErrNotImplemented = httpAsRPCError(http.StatusNotImplemented)
)

var allowDebugFn = func() zm.AllowDebugFunc {
	return func(req *http.Request) bool {
		return req != nil && req.FormValue("__level") == "5"
	}
}

func httpAsRPCError(code int) *zenrpc.Error {
	return zenrpc.NewStringError(code, http.StatusText(code))
}

// New returns new zenrpc Server.
func New(dbo db.DB, logger embedlog.Logger, isDevel bool) *zenrpc.Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
	})

	commonRepo := db.NewCommonRepo(dbo)

	// middleware
	rpc.Use(
		zm.WithHeaders(),
		zm.WithDevel(isDevel),
		zm.WithNoCancelContext(),
		zm.WithMetrics("vt"),
		zm.WithSLog(logger.Print, zm.DefaultServerName, nil),
		zm.WithErrorSLog(logger.Error, zm.DefaultServerName, nil),
		zm.WithSQLLogger(dbo.DB, isDevel, allowDebugFn(), allowDebugFn()),
		zm.WithTiming(isDevel, allowDebugFn()),
		zm.WithSentry(zm.DefaultServerName),
		authMiddleware(&commonRepo, logger),
	)

	// services
	rpc.RegisterAll(map[string]zenrpc.Invoker{
		NSAuth: NewAuthService(dbo, logger),
		NSUser: NewUserService(dbo, logger),

		NSCourse:    NewCourseService(dbo, logger),
		NSExam:      NewExamService(dbo, logger),
		NSQuestion:  NewQuestionService(dbo, logger),
		NSStudent:   NewStudentService(dbo, logger),
		NSVfsFile:   NewVfsFileService(dbo, logger),
		NSVfsFolder: NewVfsFolderService(dbo, logger),
	})

	return rpc
}
