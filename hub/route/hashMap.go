package route

import (
	"context"
	"fmt"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
	CC "github.com/karlseguin/ccache/v2"
	"net"
	"net/http"
	"strings"

	R "github.com/finddiff/clashWithCache/rules"
	"github.com/finddiff/clashWithCache/tunnel"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type HashMap struct {
	Key   string `json:"key"`
	Value string `json:"proxyName"`
}

func hashMapRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getHashMaps)
	r.Post("/", updateHashMap)

	r.Route("/{key}", func(r chi.Router) {
		r.Use(parseHashMapKey, findHashMapByKey)
		r.Get("/", getHashMap)
	})
	return r
}

func parseHashMapKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := getEscapeParam(r, "key")
		ctx := context.WithValue(r.Context(), CtxKeyHashMapKey, name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func findHashMapByKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Context().Value(CtxKeyHashMapKey).(string)
		value := tunnel.Cm.Get(key)
		if value == nil {
			//value, exist := tunnel.Cm.Get(key)
			//if !exist {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyHashMapValue, HashMap{
			Key: key,
			//Value: value.(string),
			Value: value.Value().(string),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getHashMap(w http.ResponseWriter, r *http.Request) {
	hasmap := r.Context().Value(CtxKeyHashMapValue)
	render.JSON(w, r, hasmap)
}

func getHashMaps(w http.ResponseWriter, r *http.Request) {
	items := []HashMap{}
	//for item := range tunnel.Cm.Iter() {
	//	items = append(items, HashMap{
	//		Key:   item.Key,
	//		Value: fmt.Sprintf("%v", item.Val),
	//	})
	//}

	tunnel.Cm.ForEachFunc(func(key string, i *CC.Item) bool {
		items = append(items, HashMap{
			Key:   key,
			Value: fmt.Sprintf("%v", i.Value()),
		})
		return true
	})

	render.JSON(w, r, render.M{
		"hashMap": items,
	})
}

func clearHashMap() {
	tunnel.Cm.Clear()
}

func updateHashMap(w http.ResponseWriter, r *http.Request) {
	req := HashMap{}
	//req :=
	log.Debugln("updateHashMap r.Body:%v", r.Body)
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		log.Debugln("updateHashMap err:%v", ErrBadRequest)
		return
	}

	if strings.Contains(req.Key, ":") {
		req.Key = req.Key[:strings.Index(req.Key, ":")]
	}

	//tunnel.Cm.Set(req.Key, req.Value)
	if req.Value == "DELETE" {
		var newRule = []C.Rule{}
		if net.ParseIP(req.Key) != nil {
			req.Key = req.Key + "/32"
		}
		for _, rule := range tunnel.Rules() {
			if rule.Payload() != req.Key {
				newRule = append(newRule, rule)
			}
		}
		tunnel.UpdateRules(newRule)
	} else {
		if net.ParseIP(req.Key) != nil {
			newRule, err := R.NewIPCIDR(req.Key+"/32", req.Value, R.WithIPCIDRNoResolve(true))
			if err == nil {
				tunnel.UpdateRules(append([]C.Rule{newRule}, tunnel.Rules()...))
			}
		} else {
			tunnel.UpdateRules(append([]C.Rule{R.NewDomainKeyword(req.Key, req.Value)}, tunnel.Rules()...))
		}
	}
	tunnel.Cm.Clear()
	log.Debugln("updateHashMap set %v:%v", req.Key, req.Value)
	render.NoContent(w, r)
}
