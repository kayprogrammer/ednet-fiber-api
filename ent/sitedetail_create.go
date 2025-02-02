// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/ent/sitedetail"
)

// SiteDetailCreate is the builder for creating a SiteDetail entity.
type SiteDetailCreate struct {
	config
	mutation *SiteDetailMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (sdc *SiteDetailCreate) SetCreatedAt(t time.Time) *SiteDetailCreate {
	sdc.mutation.SetCreatedAt(t)
	return sdc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableCreatedAt(t *time.Time) *SiteDetailCreate {
	if t != nil {
		sdc.SetCreatedAt(*t)
	}
	return sdc
}

// SetUpdatedAt sets the "updated_at" field.
func (sdc *SiteDetailCreate) SetUpdatedAt(t time.Time) *SiteDetailCreate {
	sdc.mutation.SetUpdatedAt(t)
	return sdc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableUpdatedAt(t *time.Time) *SiteDetailCreate {
	if t != nil {
		sdc.SetUpdatedAt(*t)
	}
	return sdc
}

// SetName sets the "name" field.
func (sdc *SiteDetailCreate) SetName(s string) *SiteDetailCreate {
	sdc.mutation.SetName(s)
	return sdc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableName(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetName(*s)
	}
	return sdc
}

// SetEmail sets the "email" field.
func (sdc *SiteDetailCreate) SetEmail(s string) *SiteDetailCreate {
	sdc.mutation.SetEmail(s)
	return sdc
}

// SetNillableEmail sets the "email" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableEmail(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetEmail(*s)
	}
	return sdc
}

// SetPhone sets the "phone" field.
func (sdc *SiteDetailCreate) SetPhone(s string) *SiteDetailCreate {
	sdc.mutation.SetPhone(s)
	return sdc
}

// SetNillablePhone sets the "phone" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillablePhone(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetPhone(*s)
	}
	return sdc
}

// SetAddress sets the "address" field.
func (sdc *SiteDetailCreate) SetAddress(s string) *SiteDetailCreate {
	sdc.mutation.SetAddress(s)
	return sdc
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableAddress(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetAddress(*s)
	}
	return sdc
}

// SetFb sets the "fb" field.
func (sdc *SiteDetailCreate) SetFb(s string) *SiteDetailCreate {
	sdc.mutation.SetFb(s)
	return sdc
}

// SetNillableFb sets the "fb" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableFb(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetFb(*s)
	}
	return sdc
}

// SetTw sets the "tw" field.
func (sdc *SiteDetailCreate) SetTw(s string) *SiteDetailCreate {
	sdc.mutation.SetTw(s)
	return sdc
}

// SetNillableTw sets the "tw" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableTw(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetTw(*s)
	}
	return sdc
}

// SetWh sets the "wh" field.
func (sdc *SiteDetailCreate) SetWh(s string) *SiteDetailCreate {
	sdc.mutation.SetWh(s)
	return sdc
}

// SetNillableWh sets the "wh" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableWh(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetWh(*s)
	}
	return sdc
}

// SetIg sets the "ig" field.
func (sdc *SiteDetailCreate) SetIg(s string) *SiteDetailCreate {
	sdc.mutation.SetIg(s)
	return sdc
}

// SetNillableIg sets the "ig" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableIg(s *string) *SiteDetailCreate {
	if s != nil {
		sdc.SetIg(*s)
	}
	return sdc
}

// SetID sets the "id" field.
func (sdc *SiteDetailCreate) SetID(u uuid.UUID) *SiteDetailCreate {
	sdc.mutation.SetID(u)
	return sdc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (sdc *SiteDetailCreate) SetNillableID(u *uuid.UUID) *SiteDetailCreate {
	if u != nil {
		sdc.SetID(*u)
	}
	return sdc
}

// Mutation returns the SiteDetailMutation object of the builder.
func (sdc *SiteDetailCreate) Mutation() *SiteDetailMutation {
	return sdc.mutation
}

// Save creates the SiteDetail in the database.
func (sdc *SiteDetailCreate) Save(ctx context.Context) (*SiteDetail, error) {
	sdc.defaults()
	return withHooks(ctx, sdc.sqlSave, sdc.mutation, sdc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sdc *SiteDetailCreate) SaveX(ctx context.Context) *SiteDetail {
	v, err := sdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sdc *SiteDetailCreate) Exec(ctx context.Context) error {
	_, err := sdc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sdc *SiteDetailCreate) ExecX(ctx context.Context) {
	if err := sdc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sdc *SiteDetailCreate) defaults() {
	if _, ok := sdc.mutation.CreatedAt(); !ok {
		v := sitedetail.DefaultCreatedAt()
		sdc.mutation.SetCreatedAt(v)
	}
	if _, ok := sdc.mutation.UpdatedAt(); !ok {
		v := sitedetail.DefaultUpdatedAt()
		sdc.mutation.SetUpdatedAt(v)
	}
	if _, ok := sdc.mutation.Name(); !ok {
		v := sitedetail.DefaultName
		sdc.mutation.SetName(v)
	}
	if _, ok := sdc.mutation.Email(); !ok {
		v := sitedetail.DefaultEmail
		sdc.mutation.SetEmail(v)
	}
	if _, ok := sdc.mutation.Phone(); !ok {
		v := sitedetail.DefaultPhone
		sdc.mutation.SetPhone(v)
	}
	if _, ok := sdc.mutation.Address(); !ok {
		v := sitedetail.DefaultAddress
		sdc.mutation.SetAddress(v)
	}
	if _, ok := sdc.mutation.Fb(); !ok {
		v := sitedetail.DefaultFb
		sdc.mutation.SetFb(v)
	}
	if _, ok := sdc.mutation.Tw(); !ok {
		v := sitedetail.DefaultTw
		sdc.mutation.SetTw(v)
	}
	if _, ok := sdc.mutation.Wh(); !ok {
		v := sitedetail.DefaultWh
		sdc.mutation.SetWh(v)
	}
	if _, ok := sdc.mutation.Ig(); !ok {
		v := sitedetail.DefaultIg
		sdc.mutation.SetIg(v)
	}
	if _, ok := sdc.mutation.ID(); !ok {
		v := sitedetail.DefaultID()
		sdc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sdc *SiteDetailCreate) check() error {
	if _, ok := sdc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "SiteDetail.created_at"`)}
	}
	if _, ok := sdc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "SiteDetail.updated_at"`)}
	}
	if _, ok := sdc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "SiteDetail.name"`)}
	}
	if _, ok := sdc.mutation.Email(); !ok {
		return &ValidationError{Name: "email", err: errors.New(`ent: missing required field "SiteDetail.email"`)}
	}
	if _, ok := sdc.mutation.Phone(); !ok {
		return &ValidationError{Name: "phone", err: errors.New(`ent: missing required field "SiteDetail.phone"`)}
	}
	if _, ok := sdc.mutation.Address(); !ok {
		return &ValidationError{Name: "address", err: errors.New(`ent: missing required field "SiteDetail.address"`)}
	}
	if _, ok := sdc.mutation.Fb(); !ok {
		return &ValidationError{Name: "fb", err: errors.New(`ent: missing required field "SiteDetail.fb"`)}
	}
	if _, ok := sdc.mutation.Tw(); !ok {
		return &ValidationError{Name: "tw", err: errors.New(`ent: missing required field "SiteDetail.tw"`)}
	}
	if _, ok := sdc.mutation.Wh(); !ok {
		return &ValidationError{Name: "wh", err: errors.New(`ent: missing required field "SiteDetail.wh"`)}
	}
	if _, ok := sdc.mutation.Ig(); !ok {
		return &ValidationError{Name: "ig", err: errors.New(`ent: missing required field "SiteDetail.ig"`)}
	}
	return nil
}

func (sdc *SiteDetailCreate) sqlSave(ctx context.Context) (*SiteDetail, error) {
	if err := sdc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sdc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sdc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	sdc.mutation.id = &_node.ID
	sdc.mutation.done = true
	return _node, nil
}

func (sdc *SiteDetailCreate) createSpec() (*SiteDetail, *sqlgraph.CreateSpec) {
	var (
		_node = &SiteDetail{config: sdc.config}
		_spec = sqlgraph.NewCreateSpec(sitedetail.Table, sqlgraph.NewFieldSpec(sitedetail.FieldID, field.TypeUUID))
	)
	if id, ok := sdc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := sdc.mutation.CreatedAt(); ok {
		_spec.SetField(sitedetail.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := sdc.mutation.UpdatedAt(); ok {
		_spec.SetField(sitedetail.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := sdc.mutation.Name(); ok {
		_spec.SetField(sitedetail.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sdc.mutation.Email(); ok {
		_spec.SetField(sitedetail.FieldEmail, field.TypeString, value)
		_node.Email = value
	}
	if value, ok := sdc.mutation.Phone(); ok {
		_spec.SetField(sitedetail.FieldPhone, field.TypeString, value)
		_node.Phone = value
	}
	if value, ok := sdc.mutation.Address(); ok {
		_spec.SetField(sitedetail.FieldAddress, field.TypeString, value)
		_node.Address = value
	}
	if value, ok := sdc.mutation.Fb(); ok {
		_spec.SetField(sitedetail.FieldFb, field.TypeString, value)
		_node.Fb = value
	}
	if value, ok := sdc.mutation.Tw(); ok {
		_spec.SetField(sitedetail.FieldTw, field.TypeString, value)
		_node.Tw = value
	}
	if value, ok := sdc.mutation.Wh(); ok {
		_spec.SetField(sitedetail.FieldWh, field.TypeString, value)
		_node.Wh = value
	}
	if value, ok := sdc.mutation.Ig(); ok {
		_spec.SetField(sitedetail.FieldIg, field.TypeString, value)
		_node.Ig = value
	}
	return _node, _spec
}

// SiteDetailCreateBulk is the builder for creating many SiteDetail entities in bulk.
type SiteDetailCreateBulk struct {
	config
	err      error
	builders []*SiteDetailCreate
}

// Save creates the SiteDetail entities in the database.
func (sdcb *SiteDetailCreateBulk) Save(ctx context.Context) ([]*SiteDetail, error) {
	if sdcb.err != nil {
		return nil, sdcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(sdcb.builders))
	nodes := make([]*SiteDetail, len(sdcb.builders))
	mutators := make([]Mutator, len(sdcb.builders))
	for i := range sdcb.builders {
		func(i int, root context.Context) {
			builder := sdcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SiteDetailMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, sdcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, sdcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, sdcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (sdcb *SiteDetailCreateBulk) SaveX(ctx context.Context) []*SiteDetail {
	v, err := sdcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sdcb *SiteDetailCreateBulk) Exec(ctx context.Context) error {
	_, err := sdcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sdcb *SiteDetailCreateBulk) ExecX(ctx context.Context) {
	if err := sdcb.Exec(ctx); err != nil {
		panic(err)
	}
}
