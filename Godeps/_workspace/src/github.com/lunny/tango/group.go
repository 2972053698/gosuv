// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

type groupRouter struct {
	methods []string
	url     string
	c       interface{}
}

type Group struct {
	routers  []groupRouter
	handlers []Handler
}

func NewGroup() *Group {
	return &Group{
		routers:  make([]groupRouter, 0),
		handlers: make([]Handler, 0),
	}
}

func (g *Group) Use(handlers ...Handler) {
	g.handlers = append(g.handlers, handlers...)
}

func (g *Group) Get(url string, c interface{}) {
	g.Route([]string{"GET", "HEAD"}, url, c)
}

func (g *Group) Post(url string, c interface{}) {
	g.Route([]string{"POST"}, url, c)
}

func (g *Group) Head(url string, c interface{}) {
	g.Route([]string{"HEAD"}, url, c)
}

func (g *Group) Options(url string, c interface{}) {
	g.Route([]string{"OPTIONS"}, url, c)
}

func (g *Group) Trace(url string, c interface{}) {
	g.Route([]string{"TRACE"}, url, c)
}

func (g *Group) Patch(url string, c interface{}) {
	g.Route([]string{"PATCH"}, url, c)
}

func (g *Group) Delete(url string, c interface{}) {
	g.Route([]string{"DELETE"}, url, c)
}

func (g *Group) Put(url string, c interface{}) {
	g.Route([]string{"PUT"}, url, c)
}

func (g *Group) Any(url string, c interface{}) {
	g.Route(SupportMethods, url, c)
}

func (g *Group) Route(methods []string, url string, c interface{}) {
	g.routers = append(g.routers, groupRouter{methods, url, c})
}

func (g *Group) Group(p string, o interface{}) {
	gr := getGroup(o)
	for _, gchild := range gr.routers {
		g.Route(gchild.methods, joinRoute(p, gchild.url), gchild.c)
	}
}

func getGroup(o interface{}) *Group {
	var g *Group
	var gf func(*Group)
	var ok bool
	if g, ok = o.(*Group); ok {
	} else if gf, ok = o.(func(*Group)); ok {
		g = NewGroup()
		gf(g)
	} else {
		panic("not allowed group parameter")
	}
	return g
}

func joinRoute(p, url string) string {
	if len(p) == 0 || p == "/" {
		return url
	}
	return p + url
}

func (t *Tango) addGroup(p string, g *Group) {
	for _, r := range g.routers {
		t.Route(r.methods, joinRoute(p, r.url), r.c)
	}
	for _, h := range g.handlers {
		t.Use(Prefix(p, h))
	}
}

func (t *Tango) Group(p string, o interface{}) {
	t.addGroup(p, getGroup(o))
}
