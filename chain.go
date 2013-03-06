package mux

import (
	"net/http"
)

// 过滤器链
type FilterChain struct {
	filters []Filter
	cur     int // 当前要执行的filter
}

// 实例化过滤器链
func NewFilterChain(filters ...Filter) *FilterChain {
	return &FilterChain{filters: filters}
}

// AddFilter 为过滤器链 增加过滤器，支持链式调用
func (this *FilterChain) AddFilter(filter Filter) *FilterChain {
	this.filters = append(this.filters, filter)
	return this
}

// 运行过滤器链
func (this *FilterChain) Run(handler HandlerFunc, rw http.ResponseWriter, req *http.Request) {
	if this.cur < len(this.filters) {
		i := this.cur
		this.cur++
		if this.filters[i].PreFilter(rw, req) {
			this.Run(handler, rw, req)
			this.filters[i].PostFilter(rw, req)
		}
	} else {
		// 复位
		this.cur = 0
		handler(rw, req)
	}
}
