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
	// 在PreFilter返回false时，执行PreErrorHandle处理错误的情况
	PreErrorHandle(http.ResponseWriter, *http.Request)
	// 在Handler执行之后 执行
	PostFilter(http.ResponseWriter, *http.Request) bool
}

// Filter接口的空实现
// 由于不少filter都不需要实现Filter中所有的方法，因此，提供一个空实现，具体filter可以“继承”该空实现来选择实现某些Filter接口方法
type EmptyFilter struct{}

func (this *EmptyFilter) PreFilter(http.ResponseWriter, *http.Request) bool { return true }

func (this *EmptyFilter) PreErrorHandle(http.ResponseWriter, *http.Request) {}

func (this *EmptyFilter) PostFilter(http.ResponseWriter, *http.Request) bool { return true }
