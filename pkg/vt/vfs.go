package vt

import (
	"context"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type VfsFileService struct {
	zenrpc.Service
	embedlog.Logger
	vfsRepo db.VfsRepo
}

func NewVfsFileService(dbo db.DB, logger embedlog.Logger) *VfsFileService {
	return &VfsFileService{
		Logger:  logger,
		vfsRepo: db.NewVfsRepo(dbo),
	}
}

func (s VfsFileService) dbSort(ops *ViewOps) db.OpFunc {
	v := s.vfsRepo.DefaultVfsFileSort()
	if ops == nil {
		return v
	}

	switch ops.SortColumn {
	case db.Columns.VfsFile.ID, db.Columns.VfsFile.FolderID, db.Columns.VfsFile.Title, db.Columns.VfsFile.Path, db.Columns.VfsFile.Params, db.Columns.VfsFile.IsFavorite, db.Columns.VfsFile.MimeType, db.Columns.VfsFile.FileSize, db.Columns.VfsFile.FileExists, db.Columns.VfsFile.CreatedAt, db.Columns.VfsFile.StatusID:
		v = db.WithSort(db.NewSortField(ops.SortColumn, ops.SortDesc))
	}

	return v
}

// Count returns count VfsFiles according to conditions in search params.
//
//zenrpc:search VfsFileSearch
//zenrpc:return int
//zenrpc:500 Internal Error
func (s VfsFileService) Count(ctx context.Context, search *VfsFileSearch) (int, error) {
	count, err := s.vfsRepo.CountVfsFiles(ctx, search.ToDB())
	if err != nil {
		return 0, InternalError(err)
	}
	return count, nil
}

// Get returns а list of VfsFiles according to conditions in search params.
//
//zenrpc:search VfsFileSearch
//zenrpc:viewOps ViewOps
//zenrpc:return []VfsFileSummary
//zenrpc:500 Internal Error
func (s VfsFileService) Get(ctx context.Context, search *VfsFileSearch, viewOps *ViewOps) ([]VfsFileSummary, error) {
	list, err := s.vfsRepo.VfsFilesByFilters(ctx, search.ToDB(), viewOps.Pager(), s.dbSort(viewOps), s.vfsRepo.FullVfsFile())
	if err != nil {
		return nil, InternalError(err)
	}
	vfsFiles := make([]VfsFileSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		if vfsFile := NewVfsFileSummary(&list[i]); vfsFile != nil {
			vfsFiles = append(vfsFiles, *vfsFile)
		}
	}
	return vfsFiles, nil
}

// GetByID returns a VfsFile by its ID.
//
//zenrpc:id int
//zenrpc:return VfsFile
//zenrpc:500 Internal Error
//zenrpc:404 Not Found
func (s VfsFileService) GetByID(ctx context.Context, id int) (*VfsFile, error) {
	db, err := s.byID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewVfsFile(db), nil
}

func (s VfsFileService) byID(ctx context.Context, id int) (*db.VfsFile, error) {
	db, err := s.vfsRepo.VfsFileByID(ctx, id, s.vfsRepo.FullVfsFile())
	if err != nil {
		return nil, InternalError(err)
	} else if db == nil {
		return nil, ErrNotFound
	}
	return db, nil
}

// Add adds a VfsFile from the query.
//
//zenrpc:vfsFile VfsFile
//zenrpc:return VfsFile
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
func (s VfsFileService) Add(ctx context.Context, vfsFile VfsFile) (*VfsFile, error) {
	if ve := s.isValid(ctx, vfsFile, false); ve.HasErrors() {
		return nil, ve.Error()
	}

	db, err := s.vfsRepo.AddVfsFile(ctx, vfsFile.ToDB())
	if err != nil {
		return nil, InternalError(err)
	}
	return NewVfsFile(db), nil
}

// Update updates the VfsFile data identified by id from the query.
//
//zenrpc:vfsFiles VfsFile
//zenrpc:return VfsFile
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s VfsFileService) Update(ctx context.Context, vfsFile VfsFile) (bool, error) {
	if _, err := s.byID(ctx, vfsFile.ID); err != nil {
		return false, err
	}

	if ve := s.isValid(ctx, vfsFile, true); ve.HasErrors() {
		return false, ve.Error()
	}

	ok, err := s.vfsRepo.UpdateVfsFile(ctx, vfsFile.ToDB())
	if err != nil {
		return false, InternalError(err)
	}
	return ok, nil
}

// Delete deletes the VfsFile by its ID.
//
//zenrpc:id int
//zenrpc:return isDeleted
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s VfsFileService) Delete(ctx context.Context, id int) (bool, error) {
	if _, err := s.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := s.vfsRepo.DeleteVfsFile(ctx, id)
	if err != nil {
		return false, InternalError(err)
	}
	return ok, err
}

// Validate verifies that VfsFile data is valid.
//
//zenrpc:vfsFile VfsFile
//zenrpc:return []FieldError
//zenrpc:500 Internal Error
func (s VfsFileService) Validate(ctx context.Context, vfsFile VfsFile) ([]FieldError, error) {
	isUpdate := vfsFile.ID != 0
	if isUpdate {
		_, err := s.byID(ctx, vfsFile.ID)
		if err != nil {
			return nil, err
		}
	}

	ve := s.isValid(ctx, vfsFile, isUpdate)
	if ve.HasInternalError() {
		return nil, ve.Error()
	}

	return ve.Fields(), nil
}

func (s VfsFileService) isValid(ctx context.Context, vfsFile VfsFile, isUpdate bool) Validator {
	var v Validator

	if v.CheckBasic(ctx, vfsFile); v.HasInternalError() {
		return v
	}

	// check fks
	if vfsFile.FolderID != 0 {
		item, err := s.vfsRepo.VfsFolderByID(ctx, vfsFile.FolderID)
		if err != nil {
			v.SetInternalError(err)
		} else if item == nil {
			v.Append("folderId", FieldErrorIncorrect)
		}
	}

	// custom validation starts here
	return v
}

type VfsFolderService struct {
	zenrpc.Service
	embedlog.Logger
	vfsRepo db.VfsRepo
}

func NewVfsFolderService(dbo db.DB, logger embedlog.Logger) *VfsFolderService {
	return &VfsFolderService{
		Logger:  logger,
		vfsRepo: db.NewVfsRepo(dbo),
	}
}

func (s VfsFolderService) dbSort(ops *ViewOps) db.OpFunc {
	v := s.vfsRepo.DefaultVfsFolderSort()
	if ops == nil {
		return v
	}

	switch ops.SortColumn {
	case db.Columns.VfsFolder.ID, db.Columns.VfsFolder.ParentFolderID, db.Columns.VfsFolder.Title, db.Columns.VfsFolder.IsFavorite, db.Columns.VfsFolder.CreatedAt, db.Columns.VfsFolder.StatusID:
		v = db.WithSort(db.NewSortField(ops.SortColumn, ops.SortDesc))
	}

	return v
}

// Count returns count VfsFolders according to conditions in search params.
//
//zenrpc:search VfsFolderSearch
//zenrpc:return int
//zenrpc:500 Internal Error
func (s VfsFolderService) Count(ctx context.Context, search *VfsFolderSearch) (int, error) {
	count, err := s.vfsRepo.CountVfsFolders(ctx, search.ToDB())
	if err != nil {
		return 0, InternalError(err)
	}
	return count, nil
}

// Get returns а list of VfsFolders according to conditions in search params.
//
//zenrpc:search VfsFolderSearch
//zenrpc:viewOps ViewOps
//zenrpc:return []VfsFolderSummary
//zenrpc:500 Internal Error
func (s VfsFolderService) Get(ctx context.Context, search *VfsFolderSearch, viewOps *ViewOps) ([]VfsFolderSummary, error) {
	list, err := s.vfsRepo.VfsFoldersByFilters(ctx, search.ToDB(), viewOps.Pager(), s.dbSort(viewOps), s.vfsRepo.FullVfsFolder())
	if err != nil {
		return nil, InternalError(err)
	}
	vfsFolders := make([]VfsFolderSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		if vfsFolder := NewVfsFolderSummary(&list[i]); vfsFolder != nil {
			vfsFolders = append(vfsFolders, *vfsFolder)
		}
	}
	return vfsFolders, nil
}

// GetByID returns a VfsFolder by its ID.
//
//zenrpc:id int
//zenrpc:return VfsFolder
//zenrpc:500 Internal Error
//zenrpc:404 Not Found
func (s VfsFolderService) GetByID(ctx context.Context, id int) (*VfsFolder, error) {
	db, err := s.byID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewVfsFolder(db), nil
}

func (s VfsFolderService) byID(ctx context.Context, id int) (*db.VfsFolder, error) {
	db, err := s.vfsRepo.VfsFolderByID(ctx, id, s.vfsRepo.FullVfsFolder())
	if err != nil {
		return nil, InternalError(err)
	} else if db == nil {
		return nil, ErrNotFound
	}
	return db, nil
}

// Add adds a VfsFolder from the query.
//
//zenrpc:vfsFolder VfsFolder
//zenrpc:return VfsFolder
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
func (s VfsFolderService) Add(ctx context.Context, vfsFolder VfsFolder) (*VfsFolder, error) {
	if ve := s.isValid(ctx, vfsFolder, false); ve.HasErrors() {
		return nil, ve.Error()
	}

	db, err := s.vfsRepo.AddVfsFolder(ctx, vfsFolder.ToDB())
	if err != nil {
		return nil, InternalError(err)
	}
	return NewVfsFolder(db), nil
}

// Update updates the VfsFolder data identified by id from the query.
//
//zenrpc:vfsFolders VfsFolder
//zenrpc:return VfsFolder
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s VfsFolderService) Update(ctx context.Context, vfsFolder VfsFolder) (bool, error) {
	if _, err := s.byID(ctx, vfsFolder.ID); err != nil {
		return false, err
	}

	if ve := s.isValid(ctx, vfsFolder, true); ve.HasErrors() {
		return false, ve.Error()
	}

	ok, err := s.vfsRepo.UpdateVfsFolder(ctx, vfsFolder.ToDB())
	if err != nil {
		return false, InternalError(err)
	}
	return ok, nil
}

// Delete deletes the VfsFolder by its ID.
//
//zenrpc:id int
//zenrpc:return isDeleted
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s VfsFolderService) Delete(ctx context.Context, id int) (bool, error) {
	if _, err := s.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := s.vfsRepo.DeleteVfsFolder(ctx, id)
	if err != nil {
		return false, InternalError(err)
	}
	return ok, err
}

// Validate verifies that VfsFolder data is valid.
//
//zenrpc:vfsFolder VfsFolder
//zenrpc:return []FieldError
//zenrpc:500 Internal Error
func (s VfsFolderService) Validate(ctx context.Context, vfsFolder VfsFolder) ([]FieldError, error) {
	isUpdate := vfsFolder.ID != 0
	if isUpdate {
		_, err := s.byID(ctx, vfsFolder.ID)
		if err != nil {
			return nil, err
		}
	}

	ve := s.isValid(ctx, vfsFolder, isUpdate)
	if ve.HasInternalError() {
		return nil, ve.Error()
	}

	return ve.Fields(), nil
}

func (s VfsFolderService) isValid(ctx context.Context, vfsFolder VfsFolder, isUpdate bool) Validator {
	var v Validator

	if v.CheckBasic(ctx, vfsFolder); v.HasInternalError() {
		return v
	}

	// check fks
	if vfsFolder.ParentFolderID != nil {
		item, err := s.vfsRepo.VfsFolderByID(ctx, *vfsFolder.ParentFolderID)
		if err != nil {
			v.SetInternalError(err)
		} else if item == nil {
			v.Append("parentFolderId", FieldErrorIncorrect)
		}
	}

	// custom validation starts here
	return v
}
