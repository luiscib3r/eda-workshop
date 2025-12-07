package storage

import (
	"backend/gen/core"
	"backend/gen/storage"
	"backend/internal/infrastructure/service"
	storagedb "backend/internal/storage/db"
	"backend/internal/storage/events"
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FilesService struct {
	storage.UnimplementedFilesServiceServer
	db       storagedb.Querier
	producer *StorageProducer
}

var _ storage.FilesServiceServer = (*FilesService)(nil)
var _ service.Service = (*FilesService)(nil)

func NewFilesService(
	db storagedb.Querier,
	producer *StorageProducer,
) *FilesService {
	return &FilesService{
		db:       db,
		producer: producer,
	}
}

// GetFiles implements storage.FilesServiceServer.
func (f *FilesService) GetFiles(
	ctx context.Context,
	req *storage.GetFilesRequest,
) (*storage.GetFilesResponse, error) {
	limit := req.PageSize
	if limit <= 0 {
		limit = 10
	}
	offset := max(limit*(req.PageNumber-1), 0)
	result, err := f.db.GetFiles(ctx, storagedb.GetFilesParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	files := make([]*storage.File, len(result))
	for i, file := range result {
		files[i] = &storage.File{
			FileKey:   file.ID,
			FileName:  file.FileName,
			FileType:  file.FileType,
			FileSize:  file.FileSize,
			CreatedAt: file.CreatedAt.Time.UTC().Format(time.RFC3339),
		}
	}

	var totalItems int32
	if len(result) > 0 {
		totalItems = int32(result[0].Total)
	}

	pageNumber := max(req.PageNumber, 1)

	pagination := &core.Pagination{
		PageNumber:  pageNumber,
		PageSize:    limit,
		TotalItems:  totalItems,
		HasNextPage: int32(pageNumber*limit) < totalItems,
	}

	return &storage.GetFilesResponse{
		Files:      files,
		Pagination: pagination,
	}, nil
}

// DeleteFiles implements storage.FilesServiceServer.
func (f *FilesService) DeleteFiles(
	ctx context.Context,
	req *storage.DeleteFilesRequest,
) (*emptypb.Empty, error) {
	err := f.db.DeleteFilesByIDs(ctx, req.FileKeys)
	if err != nil {
		return nil, err
	}

	event := events.NewFilesDeletedEvent(
		&storage.FilesDeletedEventData{
			FileKeys: req.FileKeys,
		},
	)

	err = f.producer.Publish(ctx, event)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Register implements service.Service.
func (f *FilesService) Register(ctx context.Context, mux *runtime.ServeMux) {
	storage.RegisterFilesServiceHandlerServer(ctx, mux, f)
}
