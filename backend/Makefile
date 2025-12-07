.PHONY: generate
generate:
	buf generate

.PHONY: sql
sql:
	sqlc generate

.PHONY: tools
tools:
	go install \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
	go install github.com/bufbuild/buf/cmd/buf@v1.60.0 && \
	go install -tags 'pgx5' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
	buf dep update