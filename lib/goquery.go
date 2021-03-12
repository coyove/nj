package lib

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/coyove/script"
	"golang.org/x/net/html"
)

func init() {
	script.AddGlobalValue("goquery", func(env *script.Env) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(env.Get(0).MustString("goquery", 0)))
		if err != nil {
			env.A = script.Interface(err)
			return
		}
		r := qResult{doc.Find(env.Get(1).StringDefault(""))}
		env.A = script.Interface(r)
	})
}

type qResult struct {
	sel *goquery.Selection
}

func (r qResult) Attrs() map[string]string {
	a := r.sel.Nodes[0].Attr
	m := make(map[string]string, len(a))
	for _, a := range a {
		m[a.Key] = a.Val
	}
	return m
}

func (r qResult) Find(q string) qResult {
	return qResult{r.sel.Find(q)}
}

func (r qResult) Attr(name string) script.Value {
	text, ok := r.sel.Attr(name)
	if !ok {
		return script.Value{}
	}
	return script.Interface(text)
}

func (r qResult) Text() string {
	return r.sel.Text()
}

func (r qResult) Html() (string, error) {
	return r.sel.Html()
}

func (r qResult) Next() qResult {
	return qResult{r.sel.Next()}
}

func (r qResult) NextFiltered(q string) qResult {
	return qResult{r.sel.NextFiltered(q)}
}

func (r qResult) NextAll() qResult {
	return qResult{r.sel.NextAll()}
}

func (r qResult) NextUntil(q string) qResult {
	return qResult{r.sel.NextUntil(q)}
}

func (r qResult) NextFilteredUntil(s, s2 string) qResult {
	return qResult{r.sel.NextFilteredUntil(s, s2)}
}

func (r qResult) Prev() qResult {
	return qResult{r.sel.Prev()}
}

func (r qResult) PrevFiltered(q string) qResult {
	return qResult{r.sel.PrevFiltered(q)}
}

func (r qResult) PrevAll() qResult {
	return qResult{r.sel.PrevAll()}
}

func (r qResult) PrevUntil(q string) qResult {
	return qResult{r.sel.PrevUntil(q)}
}

func (r qResult) PrevFilteredUntil(s, s2 string) qResult {
	return qResult{r.sel.PrevFilteredUntil(s, s2)}
}

func (r qResult) Nodes() []script.Value {
	x := make([]script.Value, len(r.sel.Nodes))
	for i, n := range r.sel.Nodes {
		s := *r.sel
		s.Nodes = []*html.Node{n}
		x[i] = script.Interface(qResult{&s})
	}
	return x
}
