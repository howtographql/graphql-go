package resolvers

type Resolver struct {}

func (r *Resolver) Hello() string {
	return "world"
}
