package cache

import "time"

type Item struct {
	Object     interface{} //真正的数据项
	Expiration int64       //生存时间
}



func (item Item) Expired() bool {
	if 0 == item.Expiration {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}