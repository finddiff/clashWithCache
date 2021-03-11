package dns

import (
	"bytes"
	"encoding/gob"
	"github.com/finddiff/clashWithCache/common/cache"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
	"github.com/finddiff/clashWithCache/tunnel"
	D "github.com/miekg/dns"
	"github.com/xujiajun/nutsdb"
	P "path"
	"sync"
	"time"
)

const MapDomainIPs string = "mapDomain-IPs"
const MapDomainIPttl string = "mapDomainIP-ttl"
const MapIPDomain string = "mapIP-Domain"
const MapDomainDnsMsg string = "mapDomain-DnsMsg"
const MaxDnsMsgAge = 3 * 24 * 3600

type DnsMap struct {
	ipstr   string
	domain  string
	ttl     uint32
	mapping *cache.LruCache
}

type DnsMsgMap struct {
	key   string
	value D.Msg
}

var (
	saveMapQueue = make(chan DnsMap, 500)
	saveDnsQueue = make(chan DnsMsgMap, 500)
	mu           sync.Mutex
	startED      = false
	//db           *sql.DB
	db        *nutsdb.DB
	SplitChat = "$"
)

func DnsMapAdd(dnsMap DnsMap) {
	select {
	case saveMapQueue <- dnsMap:
	default:
		log.Debugln("DnsMapAdd is block!")
	}
	//saveMapQueue <- dnsMap
}

func handleDnsMap(dnsMap DnsMap) {
	if dnsMap.ipstr == "" {
		return
	}

	dnsMap.ttl += 30

	err := db.Update(func(tx *nutsdb.Tx) error {
		//key := []byte(dnsMap.domain)
		//val := ""

		//item, err := tx.Get(MapDomainIPs, key)
		//if err != nil {
		//	log.Debugln("tx.Get(MapDomainIPs, key:%v) error %v", key, err)
		//} else {
		//	for _, ipstr := range strings.Split(string(item.Value), SplitChat) {
		//		log.Debugln("_, ipstr := range strings.Split(string(item.Value); ipstr:%s", ipstr)
		//		if len(ipstr) == 0 {
		//			continue
		//		}
		//		if item, err1 := tx.Get(MapDomainIPttl, []byte(dnsMap.domain+ipstr)); err1 == nil || item != nil {
		//			val += SplitChat + ipstr
		//		} else {
		//			err1 := tx.Delete(MapIPDomain, []byte(dnsMap.ipstr))
		//			if err1 != nil {
		//				log.Errorln("tx.Delete(MapIPDomain, []byte(dnsMap.ipstr)) error %v", err1)
		//			}
		//		}
		//	}
		//}
		//
		//if len(val) != 0 {
		//	val += SplitChat
		//}
		//val += dnsMap.ipstr
		//
		//err = tx.Put(MapDomainIPs, key, []byte(val), 0)
		//if err != nil {
		//	log.Errorln("tx.Put(MapDomainIPs, key, []byte(val), 0) %v", err)
		//}
		//
		//log.Debugln("tx.Put(MapDomainIPs, key, []byte(val) key:%v, val:%v", string(key), val)

		//add new to maps
		log.Debugln("DnsMapAdd add new to maps ip:%s| host:%s| expire Time:%v| ttl:%d|", dnsMap.ipstr, dnsMap.domain, time.Second*time.Duration(0), dnsMap.ttl)
		//err = tx.Put(MapDomainIPttl, []byte(dnsMap.domain+dnsMap.ipstr), []byte(strconv.Itoa(int(dnsMap.ttl))), 0)
		//if err != nil {
		//	log.Errorln("tx.Put(MapDomainIPs, key, val, dnsMap.ttl) %v", err)
		//}
		err := tx.Put(MapIPDomain, []byte(dnsMap.ipstr), []byte(dnsMap.domain), 0)
		if err != nil {
			log.Errorln("tx.Put(MapDomainIPs, key, []byte(val), 0) %v", err)
		}
		return nil
	})
	if err != nil {
		log.Errorln("DnsMapAdd db.Update(func(tx *nutsdb.Tx) error  %v", err)
	}
	go tunnel.DnsPreCache(dnsMap.domain)
}

func IPDomainMapOnEvict(key interface{}, value interface{}) {
	err := db.Update(func(tx *nutsdb.Tx) error {
		//tx.Delete(MapDomainIPttl, []byte(key.(string) + value.(string)))
		tx.Delete(MapIPDomain, []byte(key.(string)))
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}
}

func loadToIPDomainMap(mapping *cache.LruCache) {
	err := db.Update(func(tx *nutsdb.Tx) error {
		//db.Merge()
		entries, _ := tx.GetAll(MapIPDomain)
		for _, entry := range entries {
			//fmt.Println(string(entry.Key), string(entry.Value))
			//domainStr := string(entry.Key)
			//ipstrList := string(entry.Value)
			ip := string(entry.Key)
			domainStr := string(entry.Value)
			log.Infoln("loadToIPDomainMap SetWithExpire ip:%s| host:%s| expire Time:%v| ttl:%d|", ip, domainStr, time.Second*time.Duration(3), 3)
			mapping.SetWithExpire(ip, domainStr, time.Now().Add(time.Second*time.Duration(3)))
			//newValue := ""
			//lastIP := ""
			//for _, ipstr := range strings.Split(ipstrList, SplitChat) {
			//	if len(ipstr) > 0 {
			//		lastIP = ipstr
			//		log.Debugln("loadToIPDomainMap tx.Get(MapDomainIPttl, []byte(domainStr+ipstr) domainStr:%s, ipstr:%s", domainStr, ipstr)
			//		ttlItem, err1 := tx.Get(MapDomainIPttl, []byte(domainStr+ipstr))
			//		if err1 != nil {
			//			log.Errorln("tx.Get(MapDomainIPttl, []byte(domainStr+ipstr)) error %v", err1)
			//			continue
			//		} else {
			//			newValue += SplitChat + ipstr
			//		}
			//
			//		//log.Debugln("string(ttlItem.Value):%v", string(ttlItem.Value))
			//		ttl, err2 := strconv.Atoi(string(ttlItem.Value))
			//		if err2 != nil {
			//			log.Errorln("strconv.Atoi(string(ttlItem.Value)) error %v", err2)
			//			continue
			//		}
			//		log.Infoln("loadToIPDomainMap SetWithExpire ip:%s| host:%s| expire Time:%v| ttl:%d|", ipstr, domainStr, time.Second*time.Duration(ttl), ttl)
			//		mapping.SetWithExpire(ipstr, domainStr, time.Now().Add(time.Second*time.Duration(ttl)))
			//	}
			//}
			//if len(newValue) == 0 {
			//	log.Infoln("loadToIPDomainMap SetWithExpire ip:%s| host:%s| expire Time:%v| ttl:%d|", lastIP, domainStr, time.Second*time.Duration(3), 3)
			//	mapping.SetWithExpire(lastIP, domainStr, time.Now().Add(time.Second*time.Duration(3)))
			//	//tx.Delete(MapDomainIPs, []byte(domainStr))
			//} else {
			//	tx.Put(MapDomainIPs, []byte(domainStr), []byte(newValue), 0)
			//}
		}
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}

	mapping.SetOnEvict(IPDomainMapOnEvict)
}

func DnsMsgAdd(dnsMsg DnsMsgMap) {
	select {
	case saveDnsQueue <- dnsMsg:
	default:
		log.Debugln("DnsMsgAdd is block!")
	}
	//saveMapQueue <- dnsMap
}

func DnsMsg2Byte(p interface{}) (rb []byte, err error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(p)
	if err != nil {
		log.Errorln("Struct2Byte gob err:%v", err)
	}
	return buf.Bytes(), err
}

func Byte2DnsMsg(buf []byte) (dnsMsg D.Msg, err error) {
	enc := gob.NewDecoder(bytes.NewReader(buf))
	err = enc.Decode(&dnsMsg)
	if err != nil {
		log.Errorln("Byte2DnsMsg gob err:%v", err)
		//return dnsMsg, err
	}
	return dnsMsg, err
}

func handleDnsMsg() {
	for dnsMsg := range saveDnsQueue {
		value, err := DnsMsg2Byte(dnsMsg.value)
		if err != nil {
			continue
		}
		err = db.Update(func(tx *nutsdb.Tx) error {
			log.Debugln("handleDnsMsg tx.Put(MapDomainDnsMsg, []byte(dnsMsg.key:%v), value:%v", dnsMsg.key, value)
			tx.Put(MapDomainDnsMsg, []byte(dnsMsg.key), value, MaxDnsMsgAge)
			return nil
		})
		if err != nil {
			log.Errorln("handleDnsMsg db.Update(func(tx *nutsdb.Tx) error:%v", err)
		}
	}
}

func DnsMapOnEvict(key interface{}, value interface{}) {
	err := db.Update(func(tx *nutsdb.Tx) error {
		tx.Delete(MapDomainDnsMsg, []byte(key.(string)))
		return nil
	})
	if err != nil {
		log.Errorln("DnsMapOnEvict db.Update(func(tx *nutsdb.Tx) error %v", err)
	}
}

func loadToDnsMap(resolver *Resolver) {
	err := db.Update(func(tx *nutsdb.Tx) error {
		//db.Merge()
		entries, _ := tx.GetAll(MapDomainDnsMsg)
		for _, entry := range entries {
			//fmt.Println(string(entry.Key), string(entry.Value))
			//log.Debugln("loadToDnsMap entry.Key:%v entry.Value:%v", entry.Key, entry.Value)
			domainStr := string(entry.Key)
			dnsMsg, err := Byte2DnsMsg(entry.Value)
			if err != nil {
				continue
			}
			log.Infoln("loadToDnsMap SetWithExpire domainStr:%s| dnsMsg:%s| expire Time:%v| ttl:%d|", domainStr, dnsMsg, time.Second*time.Duration(3), 3)
			resolver.lruCache.SetWithExpire(domainStr, &dnsMsg, time.Now().Add(time.Second*time.Duration(3)))
		}
		return nil
	})
	if err != nil {
		log.Errorln("db.Update(func(tx *nutsdb.Tx) error %v", err)
	}

	resolver.lruCache.SetOnEvict(DnsMapOnEvict)
}

func initDB(resolver *Resolver, mapper *ResolverEnhancer) {
	if startED {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	startED = true

	gob.Register(&D.A{})
	gob.Register(&D.AAAA{})
	gob.Register(&D.PTR{})

	//dbpath := P.Join(C.Path.HomeDir(), "DNSSqlite.DB")
	dbpath := P.Join(C.Path.HomeDir(), "DNSNUTSDB")

	opt := nutsdb.DefaultOptions
	opt.Dir = dbpath
	newdb, err := nutsdb.Open(opt)
	if err != nil {
		log.Errorln("newdb, err := nutsdb.Open(opt) err:%v", err)
	}
	//newdb.Merge()
	//defer newdb.Close()

	db = newdb

	loadToIPDomainMap(mapper.mapping)
	loadToDnsMap(resolver)

	go handleDnsMsg()
	go processDnsMap()
}

func processDnsMap() {
	defer db.Close()
	queue := saveMapQueue

	for dnsMap := range queue {
		handleDnsMap(dnsMap)
	}
}
