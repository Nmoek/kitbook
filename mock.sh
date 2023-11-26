#!/usr/bin/env sh

echo "start mock tool..."

mockgen -source=D:./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
mockgen -source=D:./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
mockgen -source=D:./internal/service/article.go -package=svcmocks -destination=./internal/service/mocks/article.mock.go

mockgen -source=D:./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/sms.mock.go

mockgen -source=D:./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
mockgen -source=D:./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go
mockgen -source=D:./internal/repository/article.go -package=repomocks -destination=./internal/repository/mocks/article.mock.go
mockgen -source=D:./internal/repository/article_author.go -package=repomocks -destination=./internal/repository/mocks/article_author.mock.go
mockgen -source=D:./internal/repository/article_reader.go -package=repomocks -destination=./internal/repository/mocks/article_reader.mock.go

mockgen -source=D:./internal/repository/dao/user.go -package=daomocks -destination=./internal/repository/dao/mocks/user.mock.go
mockgen -source=D:./internal/repository/dao/article.go -package=daomocks -destination=./internal/repository/dao/mocks/article.mock.go
mockgen -source=D:./internal/repository/dao/article_author.go -package=daomocks -destination=./internal/repository/dao/mocks/article_author.mock.go
mockgen -source=D:./internal/repository/dao/article_reader.go -package=daomocks -destination=./internal/repository/dao/mocks/article_reader.mock.go


mockgen -source=D:./internal/repository/cache/user.go -package=cachemocks -destination=./internal/repository/cache/mocks/user.mock.go
mockgen -source=D:./internal/repository/cache/code.go -package=cachemocks -destination=./internal/repository/cache/mocks/code.mock.go

mockgen -package=redismocks -destination=./internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable

mockgen -source=D:./pkg/limiter/types.go -package=limitermocks -destination=./pkg/limiter/mocks/limiter.mock.go


go mod tidy


