package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func (c *Config) Open(filename string) (*sql.DB, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	cstr := c.User + ":" + c.Password + "@" + c.Protocol + "(" + c.Host + ":" + c.Port + ")/" + c.Database

	// parseTime=true	:for supporting DATETIME type
	// charset=utf8mb4	:for supporting 4 byte character-set
	// interpolateParams=true	:for supporting place holder

	db, err := sql.Open("mysql", cstr+"?parseTime=true&charset=utf8mb4&interpolateParams=true")
	if err != nil {
		return nil, err
	}

	return db, err
}
