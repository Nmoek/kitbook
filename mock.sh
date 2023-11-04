#!/usr/bin/env sh

echo "start mock tool..."

mockgen -source=D:./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
mockgen -source=D:./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
mockgen -source=D:./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
mockgen -source=D:./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go

go mod tidy


