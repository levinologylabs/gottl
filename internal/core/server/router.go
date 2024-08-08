package server

import (
	"net/http"
	"slices"
)

type Middlware func(http.Handler) http.Handler

type Mux struct {
	*http.ServeMux
	chain []Middlware
}

func NewServeMux(mx ...Middlware) *Mux {
	return &Mux{ServeMux: &http.ServeMux{}, chain: mx}
}

func (r *Mux) Use(mx ...Middlware) {
	r.chain = append(r.chain, mx...)
}

func (r *Mux) Group(fn func(r *Mux)) {
	fn(&Mux{ServeMux: r.ServeMux, chain: slices.Clone(r.chain)})
}

func (r *Mux) Handle(pattern string, fn http.Handler, mx ...Middlware) {
	r.ServeMux.Handle(pattern, r.wrap(fn.ServeHTTP, mx))
}

func (r *Mux) HandleFunc(pattern string, fn http.HandlerFunc, mx ...Middlware) {
	r.ServeMux.HandleFunc(pattern, r.wrap(fn, mx).ServeHTTP)
}

func (r *Mux) wrap(fn http.HandlerFunc, mx []Middlware) (out http.Handler) {
	out, mx = http.Handler(fn), append(slices.Clone(r.chain), mx...)

	slices.Reverse(mx)

	for _, m := range mx {
		out = m(out)
	}

	return
}
