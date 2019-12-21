package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

type Remind struct {
	Id        int    `db:"id,primarykey,autoincrement"`
	GuildId   string `db:"guild_id,size:255"`
	CreatorId string `db:"creator_id,size:255"`
}

type ServerConfig struct {
	Id          int    `db:"id,primarykey,autoincrement"`
	GuildId     string `db:"guild_id,size:255"`
	AdminRole   string `db:"admin_role,size:255"`
	MainChannel string `db:"main_channel,size:255"`
}

type ReminderRoles struct {
	Id      int    `db:"id,primarykey,autoincrement"`
	GuildId string `db:"guild_id,size:255"`
	RoleId  string `db:"role_id,size:255"`
}

var DbMap gorp.DbMap

func InitDB() {
	db, err := sql.Open("mysql", config.MysqlString)
	if err != nil {
		log.Panic("InitDB Unable to establish connection with database! ", err)
	}
	DbMap = gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}

	DbMap.AddTableWithName(ServerConfig{}, "ServerConfigs").SetKeys(true, "id")
	DbMap.AddTableWithName(Remind{}, "Reminds").SetKeys(true, "id")

	err = DbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Panic("InitDB Unable to create tables! ", err)
	}
}
