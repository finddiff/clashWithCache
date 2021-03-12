package tunnel

import (
	"github.com/finddiff/clashWithCache/Persistence"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
	R "github.com/finddiff/clashWithCache/rules"
	"github.com/xujiajun/nutsdb"
)

const MapDomainRule = "map-domain-rule"
const MapIPRule = "map-domain-rule"

var (
	TakeSpaceValue = []byte{}
)

func AddDomainRule(key string, value string) {
	err := Persistence.DB.Update(func(tx *nutsdb.Tx) error {
		//add new to maps
		log.Debugln("AddDomainRule add new to maps key:%v", key)
		err := tx.Put(MapDomainRule, []byte(key), []byte(value), 0)
		if err != nil {
			log.Errorln("tx.Put(MapDomainRule, []byte(key), []byte(value), 0) %v", err)
		}
		return nil
	})
	if err != nil {
		log.Errorln("AddDomainRule db.Update(func(tx *nutsdb.Tx) error  %v", err)
	}
}

func DeleteDomainRule(key string) {
	err := Persistence.DB.Update(func(tx *nutsdb.Tx) error {
		//add new to maps
		log.Debugln("AddDomainRule add new to maps key:%v", key)
		err := tx.Delete(MapDomainRule, []byte(key))
		if err != nil {
			log.Errorln("tx.Delete(MapDomainRule, []byte(key)) %v", err)
		}
		return nil
	})
	if err != nil {
		log.Errorln("AddDomainRule db.Update(func(tx *nutsdb.Tx) error  %v", err)
	}
}

func LoadDomainRule() []C.Rule {
	rules = []C.Rule{}
	err := Persistence.DB.View(func(tx *nutsdb.Tx) error {
		entries, _ := tx.GetAll(MapDomainRule)
		for _, entry := range entries {
			key := string(entry.Key)
			value := string(entry.Value)
			log.Infoln("LoadDomainRule add R.NewDomainKeyword(key:%v, value:%v)", key, value)
			rules = append(rules, R.NewDomainKeyword(key, value))
		}
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}
	return rules
}

func AddIPRule(key string, value string) {
	err := Persistence.DB.Update(func(tx *nutsdb.Tx) error {
		//add new to maps
		log.Debugln("AddDomainRule add new to maps key:%v", key)
		err := tx.Put(MapIPRule, []byte(key), []byte(value), 0)
		if err != nil {
			log.Errorln("tx.Put(MapDomainRule, []byte(key), []byte(value), 0) %v", err)
		}
		return nil
	})
	if err != nil {
		log.Errorln("AddDomainRule db.Update(func(tx *nutsdb.Tx) error  %v", err)
	}
}

func DeleteIPRule(key string) {
	err := Persistence.DB.Update(func(tx *nutsdb.Tx) error {
		//add new to maps
		log.Debugln("AddDomainRule add new to maps key:%v", key)
		err := tx.Delete(MapIPRule, []byte(key))
		if err != nil {
			log.Errorln("tx.Delete(MapDomainRule, []byte(key)) %v", err)
		}
		return nil
	})
	if err != nil {
		log.Errorln("AddDomainRule db.Update(func(tx *nutsdb.Tx) error  %v", err)
	}
}

func LoadIPRule() []C.Rule {
	rules = []C.Rule{}
	err := Persistence.DB.View(func(tx *nutsdb.Tx) error {
		entries, _ := tx.GetAll(MapIPRule)
		for _, entry := range entries {
			key := string(entry.Key)
			value := string(entry.Value)
			log.Infoln("LoadIPRule add R.NewIPCIDR(key:%v, value:%v, R.WithIPCIDRNoResolve(true))", key, value)
			newRule, err := R.NewIPCIDR(key, value, R.WithIPCIDRNoResolve(true))
			if err != nil {
				continue
			}
			rules = append(rules, newRule)
		}
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}
	return rules
}

func LoadRule(rules []C.Rule) []C.Rule {
	newRules := []C.Rule{}
	newRules = append(newRules, LoadDomainRule()...)
	newRules = append(newRules, LoadIPRule()...)
	newRules = append(newRules, rules...)
	return newRules
}
