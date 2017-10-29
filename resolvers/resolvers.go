package resolvers

import "github.com/neelance/graphql-go"

type Link struct {
	ID          graphql.ID
	URL         string
	Description string
}

var allLinks = []*Link{
	{
		ID:          "1",
		URL:         "http://howtographql.com",
		Description: "The best resource for learning GraphQL",
	},
	{
		ID:          "2",
		URL:         "https://golang.org/",
		Description: "A language that makes it easy to build simple, reliable, and efficient software.",
	},
}

var linkData = make(map[graphql.ID]*Link)

func init() {
	for _, l := range allLinks {
		linkData[l.ID] = l
	}
}

type Resolver struct{}

func (r *Resolver) AllLinks() *[]*linkResolver {
	var l []*linkResolver
	for _, link := range allLinks {
		l = append(l, &linkResolver{link})
	}
	return &l
}

type linkResolver struct {
	l *Link
}

func (r *linkResolver) ID() graphql.ID {
	return r.l.ID
}

func (r *linkResolver) URL() string {
	return r.l.URL
}

func (r *linkResolver) Description() string {
	return r.l.Description
}

func (r *Resolver) CreateLink(link *Link) *linkResolver {
	link.ID = "3"
	allLinks = append(allLinks, link)
	return &linkResolver{link}
}
