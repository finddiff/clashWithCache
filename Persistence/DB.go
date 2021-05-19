package Persistence

import (
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
	"github.com/xujiajun/nutsdb"
	P "path"
)

var (
	DB     *nutsdb.DB
	RuleDB *nutsdb.DB
)

func InitDB() {
	dbpath := P.Join(C.Path.HomeDir(), "DNSNUTSDB")

	opt := nutsdb.DefaultOptions
	opt.Dir = dbpath
	newdb, err := nutsdb.Open(opt)
	if err != nil {
		log.Errorln("newdb, err := nutsdb.Open(opt) err:%v", err)
	}
	DB = newdb
}

func MergeDB() {
	dbpath := P.Join(C.Path.HomeDir(), "DNSNUTSDB")

	opt := nutsdb.DefaultOptions
	opt.Dir = dbpath
	newdb, err := nutsdb.Open(opt)
	if err != nil {
		log.Errorln("newdb, err := nutsdb.Open(opt) err:%v", err)
	}
	defer newdb.Close()

	err = nil
	err = newdb.Update(func(tx *nutsdb.Tx) error {
		log.Infoln("MergeDB start working")
		newdb.Merge()
		log.Infoln("MergeDB finish working")
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}
}

func CloseDB() {
	DB.Close()
}

func InitRuleDB() {
	dbpath := P.Join(C.Path.HomeDir(), "RULENUTSDB")

	opt := nutsdb.DefaultOptions
	opt.Dir = dbpath
	newdb, err := nutsdb.Open(opt)
	if err != nil {
		log.Errorln("newdb, err := nutsdb.Open(opt) err:%v", err)
	}
	RuleDB = newdb
}

func MergeRuleDB() {
	dbpath := P.Join(C.Path.HomeDir(), "RULENUTSDB")

	opt := nutsdb.DefaultOptions
	opt.Dir = dbpath
	newdb, err := nutsdb.Open(opt)
	if err != nil {
		log.Errorln("newdb, err := nutsdb.Open(opt) err:%v", err)
	}
	defer newdb.Close()

	err = nil
	err = newdb.Update(func(tx *nutsdb.Tx) error {
		log.Infoln("MergeDB start working")
		newdb.Merge()
		log.Infoln("MergeDB finish working")
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}
}

func CloseRuleDB() {
	RuleDB.Close()
}
