package handler

import "net/http"

// Constructor представляет функцию-конструктор для обертывания HTTP-обработчика.
type Constructor func(http.Handler) http.Handler

// Chain представляет цепочку функций-конструкторов для обертывания HTTP-обработчика.
type Chain struct {
	constructors []Constructor
}

// aliceNew создает новую цепочку функций-конструкторов.
func aliceNew(constructors ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), constructors...)}
}

// Then связывает цепочку функций-конструкторов с заданным HTTP-обработчиком.
func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h)
	}
	return h
}

// ThenFunc связывает цепочку функций-конструкторов с заданным HTTP-обработчиком, представленным как http.HandlerFunc.
func (c Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}

// Append добавляет функции-конструкторы к существующей цепочке.
func (c Chain) Append(constructors ...Constructor) Chain {
	newCons := make([]Constructor, 0, len(c.constructors)+len(constructors))
	newCons = append(newCons, c.constructors...)
	newCons = append(newCons, constructors...)

	return Chain{newCons}
}

// Extend расширяет текущую цепочку функциями-конструкторами другой цепочки.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}
