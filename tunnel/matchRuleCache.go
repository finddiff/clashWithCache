package tunnel

import (
	"fmt"
	HS "github.com/cornelk/hashmap"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
	CC "github.com/karlseguin/ccache/v2"
	"golang.org/x/sync/syncmap"
	"time"
)

var (
	//Cm *concurrent_map.ConcurrentMap
	//Cm     = CMAP.New()
	Cm = CC.New(CC.Configure().MaxSize(1024 * 64).ItemsToPrune(500))
	Am = &HS.HashMap{}
	Bm = syncmap.Map{}
	//Am.Set("amount", 123)
	//TimeCm = CMAP.New()
	//RulChan = make(chan string, 10000)
	//Cm = concurrent_map.CreateConcurrentMap(1024)
)

func DnsPreCache(domain string) {
	adapter, hashRule, err := matchHashMap(&C.Metadata{
		NetWork:  C.TCP,
		Type:     C.SOCKS,
		SrcIP:    nil,
		DstIP:    nil,
		SrcPort:  "65535",
		DstPort:  "65535",
		AddrType: C.AtypDomainName,
		Host:     domain,
	})
	//TimeCm.Set(domain, 0)
	log.Debugln("DnsPreCache call return adapter:%v,hashRule:%v,err:%v", adapter, hashRule, err)
}

func setMatchHashMap(key string, value interface{}) {
	Cm.Set(key, value, time.Minute*60)
	//Bm.Store(key, value)
	//Am.Set(key, value)
	//Cm.Set(key, value)
	//TimeCm.Set(key, 0)
}

func matchHashMap(metadata *C.Metadata) (adapter C.Proxy, hashRule C.Rule, err error) {
	domainStr := fmt.Sprintf("%v", metadata)

	//if hashValue, ok := Am.Get(domainStr); ok {
	if item := Cm.Get(domainStr); item != nil {
		hashValue := item.Value()
		//if hashValue, success := Cm.Get(domainStr); success {
		//log.Debugln("Cm.Get time cost = " + time.Since(startT).String())
		switch hashValue.(type) {
		case C.Rule:
			//log.Debugln("case C.Rule time cost = " + time.Since(startT).String())
			hashRule := hashValue.(C.Rule)
			//if dlevel == 0 {
			//	return proxies[hashRule.Adapter()], hashRule, nil
			//}
			if hashRule.Match(metadata) {
				adapter, ok := proxies[hashRule.Adapter()]
				if ok {
					//setMatchHashMap(hoststr, hashRule)
					return adapter, hashRule, nil
				}
			}
		case string:
			if proxy, ok := proxies[hashValue.(string)]; ok {
				return proxy, nil, nil
			} else {
				break
			}
		}
	}

	adapter, hashRule, err = match(metadata)
	setMatchHashMap(domainStr, hashRule)
	log.Debugln("last match adapter=%v, hashRule=%v, err=%v ", adapter, hashRule, err)
	return
}
