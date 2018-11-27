package main

import "github.com/prog470dev/inori-backend/base"

func main() {
	server := base.New()
	server.Init("db_config.yml") // プロジェクトをrootとしたパス
	server.Run()
}
