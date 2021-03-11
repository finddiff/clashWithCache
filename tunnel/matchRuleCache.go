package tunnel

import (
	"fmt"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
	CC "github.com/karlseguin/ccache/v2"
	//CMAP "github.com/orcaman/concurrent-map"
	"time"
)

var (
	//Cm *concurrent_map.ConcurrentMap
	//Cm     = CMAP.New()
	Cm = CC.New(CC.Configure().MaxSize(1024 * 64).ItemsToPrune(500))
	//TimeCm = CMAP.New()
	//RulChan = make(chan string, 10000)
	//Cm = concurrent_map.CreateConcurrentMap(1024)
	//Bc *bigcache.BigCache
)

//TimeOut Cm rule
//func cmTimeOut() {
//	timeTickerChan := time.Tick(time.Second * 10)
//	for {
//		select {
//		case <-timeTickerChan:
//			for item := range TimeCm.Iter() {
//				value := item.Val.(int) + 1
//				if value < 360 {
//					//log.Debugln("cmTimeOut update %v:%v", item.Key, value)
//					TimeCm.Set(item.Key, item.Val.(int)+1)
//				} else {
//					//log.Debugln("cmTimeOut remove %v",item.Key)
//					Cm.Remove(item.Key)
//					TimeCm.Remove(item.Key)
//				}
//			}
//			//case key := <-RulChan:
//			//	if _, ok := TimeCm.Get(key); ok {
//			//		TimeCm.Set(key, 0)
//			//	}
//		}
//
//	}
//}

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
	//Cm.Set(key, value)
	//TimeCm.Set(key, 0)
}

func matchHashMap(metadata *C.Metadata) (adapter C.Proxy, hashRule C.Rule, err error) {
	//startT := time.Now()
	hoststr := fmt.Sprintf("%v", metadata)

	//全域名
	//domainList := strings.Split(hoststr, ".")
	domainStr := hoststr
	//domainListlen := len(domainList)

	//for dlevel := 0; dlevel < domainListlen; dlevel++ {
	//	//一共4级，全名，二级，三级 域名匹配
	//	if dlevel > 3 {
	//		break
	//	}
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

	//tc := time.Since(startT)	//计算耗时
	//log.Infoln("matchHashMap timecost = %v\n", tc)

	//	if len(metadata.Host) < 2 {
	//		break
	//	}
	//
	//	if dlevel != 0 {
	//		domainStr = strings.Join(domainList[domainListlen-1-dlevel:], ".")
	//		domainKey := domainList[domainListlen-dlevel-1]
	//		if hashValue, success := Cm.Get(domainKey); success {
	//			switch hashValue.(type) {
	//			case C.Rule:
	//				hashRule := hashValue.(C.Rule)
	//				if hashRule.Match(metadata) {
	//					adapter, ok := proxies[hashRule.Adapter()]
	//					if ok {
	//						setMatchHashMap(hoststr, hashRule)
	//						return adapter, hashRule, nil
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	adapter, hashRule, err = match(metadata)
	setMatchHashMap(hoststr, hashRule)
	log.Debugln("last match adapter=%v, hashRule=%v, err=%v ", adapter, hashRule, err)
	return
}
