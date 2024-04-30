// Code generated by ent, DO NOT EDIT.

package intercept

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/vovanwin/template/internal/shared/store/gen"
	"github.com/vovanwin/template/internal/shared/store/gen/post"
	"github.com/vovanwin/template/internal/shared/store/gen/predicate"
	"github.com/vovanwin/template/internal/shared/store/gen/user"
)

// The Query interface represents an operation that queries a graph.
// By using this interface, users can write generic code that manipulates
// query builders of different types.
type Query interface {
	// Type returns the string representation of the query type.
	Type() string
	// Limit the number of records to be returned by this query.
	Limit(int)
	// Offset to start from.
	Offset(int)
	// Unique configures the query builder to filter duplicate records.
	Unique(bool)
	// Order specifies how the records should be ordered.
	Order(...func(*sql.Selector))
	// WhereP appends storage-level predicates to the query builder. Using this method, users
	// can use type-assertion to append predicates that do not depend on any generated package.
	WhereP(...func(*sql.Selector))
}

// The Func type is an adapter that allows ordinary functions to be used as interceptors.
// Unlike traversal functions, interceptors are skipped during graph traversals. Note that the
// implementation of Func is different from the one defined in entgo.io/ent.InterceptFunc.
type Func func(context.Context, Query) error

// Intercept calls f(ctx, q) and then applied the next Querier.
func (f Func) Intercept(next gen.Querier) gen.Querier {
	return gen.QuerierFunc(func(ctx context.Context, q gen.Query) (gen.Value, error) {
		query, err := NewQuery(q)
		if err != nil {
			return nil, err
		}
		if err := f(ctx, query); err != nil {
			return nil, err
		}
		return next.Query(ctx, q)
	})
}

// The TraverseFunc type is an adapter to allow the use of ordinary function as Traverser.
// If f is a function with the appropriate signature, TraverseFunc(f) is a Traverser that calls f.
type TraverseFunc func(context.Context, Query) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseFunc) Intercept(next gen.Querier) gen.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseFunc) Traverse(ctx context.Context, q gen.Query) error {
	query, err := NewQuery(q)
	if err != nil {
		return err
	}
	return f(ctx, query)
}

// The PostFunc type is an adapter to allow the use of ordinary function as a Querier.
type PostFunc func(context.Context, *gen.PostQuery) (gen.Value, error)

// Query calls f(ctx, q).
func (f PostFunc) Query(ctx context.Context, q gen.Query) (gen.Value, error) {
	if q, ok := q.(*gen.PostQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *gen.PostQuery", q)
}

// The TraversePost type is an adapter to allow the use of ordinary function as Traverser.
type TraversePost func(context.Context, *gen.PostQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraversePost) Intercept(next gen.Querier) gen.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraversePost) Traverse(ctx context.Context, q gen.Query) error {
	if q, ok := q.(*gen.PostQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *gen.PostQuery", q)
}

// The UserFunc type is an adapter to allow the use of ordinary function as a Querier.
type UserFunc func(context.Context, *gen.UserQuery) (gen.Value, error)

// Query calls f(ctx, q).
func (f UserFunc) Query(ctx context.Context, q gen.Query) (gen.Value, error) {
	if q, ok := q.(*gen.UserQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *gen.UserQuery", q)
}

// The TraverseUser type is an adapter to allow the use of ordinary function as Traverser.
type TraverseUser func(context.Context, *gen.UserQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseUser) Intercept(next gen.Querier) gen.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseUser) Traverse(ctx context.Context, q gen.Query) error {
	if q, ok := q.(*gen.UserQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *gen.UserQuery", q)
}

// NewQuery returns the generic Query interface for the given typed query.
func NewQuery(q gen.Query) (Query, error) {
	switch q := q.(type) {
	case *gen.PostQuery:
		return &query[*gen.PostQuery, predicate.Post, post.OrderOption]{typ: gen.TypePost, tq: q}, nil
	case *gen.UserQuery:
		return &query[*gen.UserQuery, predicate.User, user.OrderOption]{typ: gen.TypeUser, tq: q}, nil
	default:
		return nil, fmt.Errorf("unknown query type %T", q)
	}
}

type query[T any, P ~func(*sql.Selector), R ~func(*sql.Selector)] struct {
	typ string
	tq  interface {
		Limit(int) T
		Offset(int) T
		Unique(bool) T
		Order(...R) T
		Where(...P) T
	}
}

func (q query[T, P, R]) Type() string {
	return q.typ
}

func (q query[T, P, R]) Limit(limit int) {
	q.tq.Limit(limit)
}

func (q query[T, P, R]) Offset(offset int) {
	q.tq.Offset(offset)
}

func (q query[T, P, R]) Unique(unique bool) {
	q.tq.Unique(unique)
}

func (q query[T, P, R]) Order(orders ...func(*sql.Selector)) {
	rs := make([]R, len(orders))
	for i := range orders {
		rs[i] = orders[i]
	}
	q.tq.Order(rs...)
}

func (q query[T, P, R]) WhereP(ps ...func(*sql.Selector)) {
	p := make([]P, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	q.tq.Where(p...)
}