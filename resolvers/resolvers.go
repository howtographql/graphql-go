package resolvers

import (
	"github.com/howtographql/graphql-go/db"
)

type Resolver struct{}

// AllLinks wraps each link formatting it to
// a resolver and appending it to the list of resolvers
func (r *Resolver) AllLinks() *[]*linkResolver {
	var l []*linkResolver
	wrapper := func(link *db.Link) {
		l = append(l, &linkResolver{link})
	}

	db.AllLinks(wrapper)

	return &l
}

type linkResolver struct {
	l *db.Link
}

func (r *linkResolver) ID() string {
	return r.l.ID
}

func (r *linkResolver) URL() string {
	return r.l.URL
}

func (r *linkResolver) Description() string {
	return r.l.Description
}

func (r *Resolver) CreateLink(link *db.Link) *linkResolver {
	db.CreateLink(link)
	return &linkResolver{link}
}
