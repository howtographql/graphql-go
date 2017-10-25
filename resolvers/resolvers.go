package resolvers

type Link struct {
	url         string
	description string
}

var allLinks = []*Link{
	{
		url:         "http://howtographql.com",
		description: "The best resource for learning GraphQL",
	},
	{
		url:         "https://golang.org/",
		description: "A language that makes it easy to build simple, reliable, and efficient software.",
	},
}

var linkData = make(map[string]*Link)

func init() {
	for _, l := range allLinks {
		linkData[l.url] = l
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

func (r *linkResolver) URL() string {
	return r.l.url
}

func (r *linkResolver) Description() string {
	return r.l.description
}
