package models

import (
	"flag"

	"github.com/guotie/config"
)

func init() {
	configFn := flag.String("config", "./test.json", "config file path")
	config.ReadCfg(*configFn)

	dropTables()

	InitDB()
}

func dropTables() {
	db := GetDB(DEFAULT_DB)
	db.DropTableIfExists(Account{})

	db.DropTableIfExists(Post{})

	db.DropTableIfExists(Reply{})

	db.DropTableIfExists(Link{})

	db.DropTableIfExists(Meta{})

	db.DropTableIfExists(Options{})

	//db.DropTableIfExists(Term{})
	db.DropTableIfExists(TermRelation{})
	db.DropTableIfExists(TermTaxonomy{})

	db.DropTableIfExists(Table{})
	db.DropTableIfExists(Column{})
}
