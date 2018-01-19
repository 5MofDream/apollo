package cache

import (
	"sort"
	"sync"
)

type issueZSet struct {
	Zm map[string]map[string]map[string]string
	sync.RWMutex
}

var IssueZSet *issueZSet

func init() {
	IssueZSet = new(issueZSet)
	IssueZSet.Zm = make(map[string]map[string]map[string]string)
}

func (obj *issueZSet) AddZSet(key string, item string, value map[string]string) {
	obj.Lock()
	defer obj.Unlock()
	if _, ok := obj.Zm[key]; !ok {
		obj.Zm[key] = make(map[string]map[string]string)
	}
	if _, ok := obj.Zm[key][item]; !ok {
		obj.Zm[key][item] = make(map[string]string)
	}
	obj.Zm[key][item] = value
}

//func (s *ZSet) Remove(item string) {
//	s.Lock()
//	defer s.Unlock()
//	delete(s.m, item)
//}
//
//func (s *ZSet) Has(item string) bool {
//	s.RLock()
//	defer s.RUnlock()
//	_, ok := s.m[item]
//	return ok
//}
//
//func (s *ZSet) Len() int {
//	return len(s.List())
//}

func (s *issueZSet) SortList(key string,skip int,take int) []map[string]map[string]string {
	list := []string{}
	for issue, _ := range s.Zm[key] {
		list = append(list, issue)
	}
	// 排序 取片
	//sort.Strings(list)
	sort.Sort(sort.Reverse(sort.StringSlice(list)))

	data := []map[string]map[string]string{}
	a := len(list)
	if a<skip || a == 0{
		return  data
	}
	keys := []string{}
	if skip<a && a<skip+take {
		keys = list[skip:a]
	}else{
		keys = list[skip:skip+take]
	}

	for _, issueId := range keys {
		item := make(map[string]map[string]string)
		item[issueId] = s.Zm[key][issueId]
		data = append(data, item)
	}
	return data
}
