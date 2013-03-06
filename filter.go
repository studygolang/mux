package mux

import (
	"net/http"
)

// 自定义HandlerFunc，以便应用过滤器
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 执行当前Route的FilterChain
	filterChain := CurrentRoute(req).FilterChain
	if filterChain != nil {
	    filterChain.Run(f, rw, req)
	    return
	}
	// 没有设置FilterChain时，直接执行Handler
	f(rw, req)
	
}

// 过滤器接口
type Filter interface {
	// 在Handler执行之前 执行
	PreFilter(http.ResponseWriter, *http.Request) bool
	// 在Handler执行之后 执行
	PostFilter(http.ResponseWriter, *http.Request) bool
}