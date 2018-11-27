GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ino.out
GOOS=linux
GOARCH=amd64

DB_USER=root
DB_PASSWORD=password
DB_NAME=ino

# テスト用のマイグレーション
migrate/init:
	# TODO: マイグレーション(up)
	mysql -u root -h localhost --protocol tcp -p$(DB_PASSWORD) < "migrations/0_init_up.sql"
	mysql -u root -h localhost --protocol tcp -p$(DB_PASSWORD) < "migrations/1_add_time_up.sql"
	mysql -u root -h localhost --protocol tcp -p$(DB_PASSWORD) < "migrations/2_tokens_init_up.sql"
	# 確認
	mysql -u root -h localhost --protocol tcp -e "show tables from $(DB_NAME)" -p$(DB_PASSWORD)

migrate/down:
	# TODO: マイグレーション(down)
	mysql -u root -h localhost --protocol tcp -p$(DB_PASSWORD) < "migrations/2_tokens_init_down.sql"
	mysql -u root -h localhost --protocol tcp -p$(DB_PASSWORD) < "migrations/1_add_time_down.sql"
	mysql -u root -h localhost --protocol tcp -p$(DB_PASSWORD) < "migrations/0_init_down.sql"

test:
	go test -v ./...

deps:
	which dep || go get -v -u github.com/golang/dep/cmd/dep
	dep ensure

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ino.out
