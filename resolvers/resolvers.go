package resolvers

import (
	"github.com/howtographql/graphql-go/db"
	"github.com/neelance/graphql-go"
	"context"
	"log"
	"time"
)

type Resolver struct{}

// AllLinks wraps each link formatting it to
// a resolver and appending it to the list of resolvers
func (r *Resolver) AllLinks(ctx context.Context) *[]*linkResolver {
	if fetchUserFromAuthorizationToken(ctx) == nil {
		return nil
	}
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

func (r *linkResolver) PostedBy() *userResolver {
	user := db.FindUserByID(r.l.PostedBy)
	return &userResolver{user}
}

func (r *linkResolver) Votes(ctx context.Context) []*voteResolver {
	if fetchUserFromAuthorizationToken(ctx) == nil {
		return nil
	}
	var v []*voteResolver
	wrapper := func(vote *db.Vote) {
		v = append(v, &voteResolver{vote})
	}

	db.FindVotesByLinkID(r.l.ID, wrapper)

	return v
}

func (r *Resolver) CreateLink(ctx context.Context, link *db.Link) *linkResolver {
	user := fetchUserFromAuthorizationToken(ctx)
	if user == nil {
		return nil
	}
	link.PostedBy = user.ID
	db.CreateLink(link)
	return &linkResolver{link}
}

type userResolver struct {
	u *db.User
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(r.u.ID)
}

func (r *userResolver) Name() string {
	return r.u.Name
}

func (r *userResolver) Email() *string {
	return &r.u.Email
}

func (r *userResolver) Password() *string {
	return &r.u.Password
}

func (r *userResolver) Votes(ctx context.Context) []*voteResolver {
	if fetchUserFromAuthorizationToken(ctx) == nil {
		return nil
	}
	var v []*voteResolver
	wrapper := func(vote *db.Vote) {
		v = append(v, &voteResolver{vote})
	}

	db.FindVotesByUserID(r.u.ID, wrapper)

	return v
}

func (r *Resolver) CreateUser(args *struct {
	Name         string
	AuthProvider *db.AuthData
}) *userResolver {
	user := &db.User{Name: args.Name, Email: args.AuthProvider.Email, Password: args.AuthProvider.Password}
	db.CreateUser(user)
	return &userResolver{user}
}

type authDataResolver struct {
	a *db.AuthData
}

func (r *authDataResolver) Email() string {
	return r.a.Email
}

func (r *authDataResolver) Password() string {
	return r.a.Password
}

type SigninPayload struct {
	Token string
	User  db.User
}

type signinPayloadResolver struct {
	s *SigninPayload
}

func (r signinPayloadResolver) Token() *string {
	return &r.s.Token
}

func (r signinPayloadResolver) User() *userResolver {
	return &userResolver{&r.s.User}
}

func (r *Resolver) SigninUser(args *struct {
	Auth *db.AuthData
}) *signinPayloadResolver {
	user := *db.FindUserByEmail(args.Auth.Email)
	signinPayload := &SigninPayload{}
	if user.Password == args.Auth.Password {
		signinPayload.Token = user.ID // encrypt here irl
		signinPayload.User = user
	}
	return &signinPayloadResolver{signinPayload}
}

func fetchUserFromAuthorizationToken(ctx context.Context) *db.User {
	token, ok := ctx.Value("AuthorizationToken").(string)
	if !ok {
		log.Println("Error occurred while parsing authorization token")
		return nil
	}

	return db.FindUserByID(token)
}

type voteResolver struct {
	v *db.Vote
}

func (r *Resolver) CreateVote(args *struct {
	LinkID *string
	UserID *string
}) *voteResolver {
	vote := &db.Vote{LinkID: *args.LinkID, UserID: *args.UserID}
	db.CreateVote(vote)
	return &voteResolver{vote}
}

func (r *voteResolver) ID() graphql.ID {
	return graphql.ID(r.v.ID)
}

func (r *voteResolver) CreatedAt() graphql.Time {
	t, err := time.Parse(time.RubyDate, r.v.CreatedAt)
	if err != nil {
		log.Println("Error occurred while parsing time", err)
	}
	return graphql.Time{t}
}

func (r *voteResolver) User() *userResolver {
	user := db.FindUserByID(r.v.UserID)
	return &userResolver{user}
}

func (r *voteResolver) Link() *linkResolver {
	link := db.FindLinkByID(r.v.LinkID)
	return &linkResolver{link}
}
