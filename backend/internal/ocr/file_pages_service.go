package ocr

import (
	"backend/gen/core"
	"backend/gen/ocr"
	"backend/internal/infrastructure/service"
	"backend/internal/infrastructure/storage"
	ocrdb "backend/internal/ocr/db"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilesService struct {
	ocr.UnimplementedFilePagesServiceServer
	s3 *s3.PresignClient
	db *ocrdb.Queries
}

var _ ocr.FilePagesServiceServer = (*FilesService)(nil)
var _ service.Service = (*FilesService)(nil)

func NewFilesService(
	s3 *s3.PresignClient,
	db *ocrdb.Queries,
) *FilesService {
	return &FilesService{
		s3: s3,
		db: db,
	}
}

// GetFilePages implements ocr.FilePagesServiceServer.
func (f *FilesService) GetFilePages(
	ctx context.Context,
	req *ocr.GetFilePagesRequest,
) (*ocr.GetFilePagesResponse, error) {
	limit := req.PageSize
	if limit <= 0 {
		limit = 10
	}

	offset := max(limit*(req.PageNumber-1), 0)

	fileId, err := uuid.Parse(req.FileKey)
	if err != nil {
		return nil, err
	}

	result, err := f.db.GetFilePagesByFileID(
		ctx, ocrdb.GetFilePagesByFileIDParams{
			FileID: pgtype.UUID{
				Bytes: fileId,
				Valid: true,
			},
			Limit:  limit,
			Offset: offset,
		},
	)

	pages := make([]*ocr.FilePage, len(result))
	for i, page := range result {
		pages[i] = &ocr.FilePage{
			Id:         page.ID.String(),
			PageNumber: page.PageNumber + 1,
		}

		if imageUrl, err := f.s3.PresignGetObject(ctx, &s3.GetObjectInput{
			Key:    &page.PageImageKey,
			Bucket: aws.String(storage.BUCKET_NAME),
		}); err == nil {
			pages[i].ImageUrl = imageUrl.URL
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

	return &ocr.GetFilePagesResponse{
		Pages:      pages,
		Pagination: pagination,
	}, nil
}

// GetFilePageContent implements ocr.FilePagesServiceServer.
func (f *FilesService) GetFilePageContent(
	ctx context.Context,
	req *ocr.GetFilePageContentRequest,
) (*ocr.GetFilePageContentResponse, error) {
	pageId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	content, err := f.db.GetFilePageContentByID(ctx, pgtype.UUID{
		Bytes: pageId,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	if content != nil {
		return &ocr.GetFilePageContentResponse{
			Content: *content,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "file page content not found")
}

// Register implements service.Service.
func (f *FilesService) Register(ctx context.Context, mux *runtime.ServeMux) {
	ocr.RegisterFilePagesServiceHandlerServer(ctx, mux, f)
}
