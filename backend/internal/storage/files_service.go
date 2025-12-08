package storage

import (
	"backend/gen/core"
	"backend/gen/storage"
	"backend/internal/infrastructure/service"
	storagedb "backend/internal/storage/db"
	"backend/internal/storage/events"
	"context"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FilesService struct {
	storage.UnimplementedFilesServiceServer
	db   *storagedb.Queries
	pool *pgxpool.Pool
}

var _ storage.FilesServiceServer = (*FilesService)(nil)
var _ service.Service = (*FilesService)(nil)

func NewFilesService(
	db *storagedb.Queries,
	pool *pgxpool.Pool,
) *FilesService {
	return &FilesService{
		db:   db,
		pool: pool,
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
			FileKey:   file.ID.String(),
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
	pageSize := int32(min(len(result), int(limit)))

	pagination := &core.Pagination{
		PageNumber:  pageNumber,
		PageSize:    pageSize,
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
	tx, err := f.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := f.db.WithTx(tx)

	ids := lo.Map(req.FileKeys, func(key string, _ int) pgtype.UUID {
		id, err := uuid.Parse(key)
		if err != nil {
			return pgtype.UUID{
				Valid: false,
			}
		}

		return pgtype.UUID{
			Bytes: id,
			Valid: true,
		}
	})

	ids = lo.Filter(ids, func(id pgtype.UUID, _ int) bool {
		return id.Valid
	})

	err = qtx.DeleteFilesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	event := events.NewFilesDeletedEvent(
		&storage.FilesDeletedEventData{
			FileKeys: req.FileKeys,
		},
	)

	eventId := event.Id
	eventType := event.Type()
	payload, err := protojson.Marshal(event.Data())
	if err != nil {
		return nil, err
	}

	err = qtx.CreateOutboxEvent(ctx, storagedb.CreateOutboxEventParams{
		EventID: pgtype.UUID{
			Bytes: eventId,
			Valid: true,
		},
		EventType: eventType,
		Payload:   payload,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Register implements service.Service.
func (f *FilesService) Register(ctx context.Context, mux *runtime.ServeMux) {
	storage.RegisterFilesServiceHandlerServer(ctx, mux, f)
}
