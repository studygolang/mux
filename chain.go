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

// Append 将两个过滤器链合并（不排除重复的过滤器）
func (this *FilterChain) Append(filterChain *FilterChain) *FilterChain {
	if len(this.filters) == 0 {
		this.filters = filterChain.filters
	} else {
		this.filters = append(this.filters, filterChain.filters...)
	}
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
		} else {
			// 错误处理中，过滤器链不应该往下执行了。
			this.cur = 0
			// 执行错误处理
			this.filters[i].PreErrorHandle(rw, req)
		}
	} else {
		// 复位
		this.cur = 0
		handler(rw, req)
	}
}
