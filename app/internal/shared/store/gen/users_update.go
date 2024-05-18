// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/vovanwin/template/internal/shared/store/gen/predicate"
	"github.com/vovanwin/template/internal/shared/store/gen/users"
)

// UsersUpdate is the builder for updating Users entities.
type UsersUpdate struct {
	config
	hooks    []Hook
	mutation *UsersMutation
}

// Where appends a list predicates to the UsersUpdate builder.
func (uu *UsersUpdate) Where(ps ...predicate.Users) *UsersUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetLogin sets the "login" field.
func (uu *UsersUpdate) SetLogin(s string) *UsersUpdate {
	uu.mutation.SetLogin(s)
	return uu
}

// SetNillableLogin sets the "login" field if the given value is not nil.
func (uu *UsersUpdate) SetNillableLogin(s *string) *UsersUpdate {
	if s != nil {
		uu.SetLogin(*s)
	}
	return uu
}

// SetPassword sets the "password" field.
func (uu *UsersUpdate) SetPassword(s string) *UsersUpdate {
	uu.mutation.SetPassword(s)
	return uu
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (uu *UsersUpdate) SetNillablePassword(s *string) *UsersUpdate {
	if s != nil {
		uu.SetPassword(*s)
	}
	return uu
}

// SetDeletedAt sets the "deleted_at" field.
func (uu *UsersUpdate) SetDeletedAt(t time.Time) *UsersUpdate {
	uu.mutation.SetDeletedAt(t)
	return uu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (uu *UsersUpdate) SetNillableDeletedAt(t *time.Time) *UsersUpdate {
	if t != nil {
		uu.SetDeletedAt(*t)
	}
	return uu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (uu *UsersUpdate) ClearDeletedAt() *UsersUpdate {
	uu.mutation.ClearDeletedAt()
	return uu
}

// SetCreatedAt sets the "created_at" field.
func (uu *UsersUpdate) SetCreatedAt(t time.Time) *UsersUpdate {
	uu.mutation.SetCreatedAt(t)
	return uu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (uu *UsersUpdate) SetNillableCreatedAt(t *time.Time) *UsersUpdate {
	if t != nil {
		uu.SetCreatedAt(*t)
	}
	return uu
}

// SetUpdatedAt sets the "updated_at" field.
func (uu *UsersUpdate) SetUpdatedAt(t time.Time) *UsersUpdate {
	uu.mutation.SetUpdatedAt(t)
	return uu
}

// Mutation returns the UsersMutation object of the builder.
func (uu *UsersUpdate) Mutation() *UsersMutation {
	return uu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UsersUpdate) Save(ctx context.Context) (int, error) {
	uu.defaults()
	return withHooks(ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UsersUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UsersUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UsersUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uu *UsersUpdate) defaults() {
	if _, ok := uu.mutation.UpdatedAt(); !ok {
		v := users.UpdateDefaultUpdatedAt()
		uu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uu *UsersUpdate) check() error {
	if v, ok := uu.mutation.Password(); ok {
		if err := users.PasswordValidator(v); err != nil {
			return &ValidationError{Name: "password", err: fmt.Errorf(`gen: validator failed for field "Users.password": %w`, err)}
		}
	}
	return nil
}

func (uu *UsersUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := uu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(users.Table, users.Columns, sqlgraph.NewFieldSpec(users.FieldID, field.TypeUUID))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.Login(); ok {
		_spec.SetField(users.FieldLogin, field.TypeString, value)
	}
	if value, ok := uu.mutation.Password(); ok {
		_spec.SetField(users.FieldPassword, field.TypeString, value)
	}
	if value, ok := uu.mutation.DeletedAt(); ok {
		_spec.SetField(users.FieldDeletedAt, field.TypeTime, value)
	}
	if uu.mutation.DeletedAtCleared() {
		_spec.ClearField(users.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := uu.mutation.CreatedAt(); ok {
		_spec.SetField(users.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := uu.mutation.UpdatedAt(); ok {
		_spec.SetField(users.FieldUpdatedAt, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{users.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UsersUpdateOne is the builder for updating a single Users entity.
type UsersUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UsersMutation
}

// SetLogin sets the "login" field.
func (uuo *UsersUpdateOne) SetLogin(s string) *UsersUpdateOne {
	uuo.mutation.SetLogin(s)
	return uuo
}

// SetNillableLogin sets the "login" field if the given value is not nil.
func (uuo *UsersUpdateOne) SetNillableLogin(s *string) *UsersUpdateOne {
	if s != nil {
		uuo.SetLogin(*s)
	}
	return uuo
}

// SetPassword sets the "password" field.
func (uuo *UsersUpdateOne) SetPassword(s string) *UsersUpdateOne {
	uuo.mutation.SetPassword(s)
	return uuo
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (uuo *UsersUpdateOne) SetNillablePassword(s *string) *UsersUpdateOne {
	if s != nil {
		uuo.SetPassword(*s)
	}
	return uuo
}

// SetDeletedAt sets the "deleted_at" field.
func (uuo *UsersUpdateOne) SetDeletedAt(t time.Time) *UsersUpdateOne {
	uuo.mutation.SetDeletedAt(t)
	return uuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (uuo *UsersUpdateOne) SetNillableDeletedAt(t *time.Time) *UsersUpdateOne {
	if t != nil {
		uuo.SetDeletedAt(*t)
	}
	return uuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (uuo *UsersUpdateOne) ClearDeletedAt() *UsersUpdateOne {
	uuo.mutation.ClearDeletedAt()
	return uuo
}

// SetCreatedAt sets the "created_at" field.
func (uuo *UsersUpdateOne) SetCreatedAt(t time.Time) *UsersUpdateOne {
	uuo.mutation.SetCreatedAt(t)
	return uuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (uuo *UsersUpdateOne) SetNillableCreatedAt(t *time.Time) *UsersUpdateOne {
	if t != nil {
		uuo.SetCreatedAt(*t)
	}
	return uuo
}

// SetUpdatedAt sets the "updated_at" field.
func (uuo *UsersUpdateOne) SetUpdatedAt(t time.Time) *UsersUpdateOne {
	uuo.mutation.SetUpdatedAt(t)
	return uuo
}

// Mutation returns the UsersMutation object of the builder.
func (uuo *UsersUpdateOne) Mutation() *UsersMutation {
	return uuo.mutation
}

// Where appends a list predicates to the UsersUpdate builder.
func (uuo *UsersUpdateOne) Where(ps ...predicate.Users) *UsersUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UsersUpdateOne) Select(field string, fields ...string) *UsersUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated Users entity.
func (uuo *UsersUpdateOne) Save(ctx context.Context) (*Users, error) {
	uuo.defaults()
	return withHooks(ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UsersUpdateOne) SaveX(ctx context.Context) *Users {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UsersUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UsersUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uuo *UsersUpdateOne) defaults() {
	if _, ok := uuo.mutation.UpdatedAt(); !ok {
		v := users.UpdateDefaultUpdatedAt()
		uuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UsersUpdateOne) check() error {
	if v, ok := uuo.mutation.Password(); ok {
		if err := users.PasswordValidator(v); err != nil {
			return &ValidationError{Name: "password", err: fmt.Errorf(`gen: validator failed for field "Users.password": %w`, err)}
		}
	}
	return nil
}

func (uuo *UsersUpdateOne) sqlSave(ctx context.Context) (_node *Users, err error) {
	if err := uuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(users.Table, users.Columns, sqlgraph.NewFieldSpec(users.FieldID, field.TypeUUID))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`gen: missing "Users.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, users.FieldID)
		for _, f := range fields {
			if !users.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("gen: invalid field %q for query", f)}
			}
			if f != users.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.Login(); ok {
		_spec.SetField(users.FieldLogin, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Password(); ok {
		_spec.SetField(users.FieldPassword, field.TypeString, value)
	}
	if value, ok := uuo.mutation.DeletedAt(); ok {
		_spec.SetField(users.FieldDeletedAt, field.TypeTime, value)
	}
	if uuo.mutation.DeletedAtCleared() {
		_spec.ClearField(users.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := uuo.mutation.CreatedAt(); ok {
		_spec.SetField(users.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := uuo.mutation.UpdatedAt(); ok {
		_spec.SetField(users.FieldUpdatedAt, field.TypeTime, value)
	}
	_node = &Users{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{users.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}