// Package sendtables contains code related to decoding sendtables.
// Mostly used internally but can be interesting for direct access to server-classes and entities.
package sendtables

import (
	"fmt"

	"github.com/golang/geo/r3"

	"github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/constants"
)

// PropertyValueType specifies the type of PropertyValue
type PropertyValueType int

// Possible types of property values.
// See Property.Bind()
const (
	ValTypeInt PropertyValueType = iota
	ValTypeFloat32
	ValTypeFloat64 // Like ValTypeFloat32 but with additional cast to float64
	ValTypeString
	ValTypeVector
	ValTypeArray
	ValTypeBoolInt // Int that is treated as bool (1 -> true, != 1 -> false)
)

// PropertyValue stores parsed & decoded send-table values.
// For instance player health, location etc.
type PropertyValue struct {
	Any any
}

// The getters below never panic: a property without a value reads as its zero value,
// and integer getters accept every integer width a decoder can produce.

func (v PropertyValue) R3Vec() r3.Vector {
	if vec := v.R3VecOrNil(); vec != nil {
		return *vec
	}

	return r3.Vector{}
}

func (v PropertyValue) R3VecOrNil() *r3.Vector {
	switch fs := v.Any.(type) {
	case [3]float32:
		return &r3.Vector{X: float64(fs[0]), Y: float64(fs[1]), Z: float64(fs[2])}
	case []float32:
		if len(fs) < 3 {
			return nil
		}

		return &r3.Vector{X: float64(fs[0]), Y: float64(fs[1]), Z: float64(fs[2])}
	default:
		return nil
	}
}

func (v PropertyValue) integer() (uint64, bool) {
	switch n := v.Any.(type) {
	case int32:
		return uint64(n), true
	case int64:
		return uint64(n), true
	case uint32:
		return uint64(n), true
	case uint64:
		return n, true
	case bool:
		if n {
			return 1, true
		}

		return 0, true
	default:
		return 0, false
	}
}

func (v PropertyValue) Int() int {
	if n, ok := v.Any.(int32); ok {
		return int(n)
	}

	n, _ := v.integer()

	return int(n)
}

func (v PropertyValue) Int64() int64 {
	if n, ok := v.Any.(int32); ok {
		return int64(n)
	}

	n, _ := v.integer()

	return int64(n)
}

func (v PropertyValue) UInt64() uint64 {
	n, _ := v.integer()

	return n
}

func (v PropertyValue) Array() []any {
	arr, _ := v.Any.([]any)

	return arr
}

func (v PropertyValue) UInt32() uint32 {
	n, _ := v.integer()

	return uint32(n)
}

// Handle returns constants.InvalidEntityHandleSource2 when the property has no value.
func (v PropertyValue) Handle() uint64 {
	n, ok := v.integer()
	if !ok {
		return constants.InvalidEntityHandleSource2
	}

	return n
}

func (v PropertyValue) Float() float32 {
	switch f := v.Any.(type) {
	case float32:
		return f
	case float64:
		return float32(f)
	default:
		return 0
	}
}

func (v PropertyValue) Str() string {
	s, _ := v.Any.(string)

	return s
}

func (v PropertyValue) String() string {
	return fmt.Sprint(v.Any)
}

// BoolVal returns the bool value, or whether an integer value is not 0.
func (v PropertyValue) BoolVal() bool {
	if b, ok := v.Any.(bool); ok {
		return b
	}

	n, _ := v.integer()

	return n != 0
}

// PropertyUpdateHandler is the interface for handlers that are interested in property changes.
type PropertyUpdateHandler func(PropertyValue)

type PropertyEntry struct {
	Name    string
	IsArray bool
}

// EntityCreatedHandler is the interface for handlers that are interested in EntityCreatedEvents.
type EntityCreatedHandler func(Entity)

// ServerClasses is a searchable list of ServerClasses.
type ServerClasses interface {
	All() []ServerClass
	FindByName(name string) ServerClass
}

// EntityOp is a bitmask representing the type of operation performed on an Entity
type EntityOp int

const (
	EntityOpNone           EntityOp = 0x00
	EntityOpCreated        EntityOp = 0x01
	EntityOpUpdated        EntityOp = 0x02
	EntityOpDeleted        EntityOp = 0x04
	EntityOpEntered        EntityOp = 0x08
	EntityOpLeft           EntityOp = 0x10
	EntityOpCreatedEntered EntityOp = EntityOpCreated | EntityOpEntered
	EntityOpUpdatedEntered EntityOp = EntityOpUpdated | EntityOpEntered
	EntityOpDeletedLeft    EntityOp = EntityOpDeleted | EntityOpLeft
)

var entityOpNames = map[EntityOp]string{
	EntityOpNone:           "None",
	EntityOpCreated:        "Created",
	EntityOpUpdated:        "Updated",
	EntityOpDeleted:        "Deleted",
	EntityOpEntered:        "Entered",
	EntityOpLeft:           "Left",
	EntityOpCreatedEntered: "Created+Entered",
	EntityOpUpdatedEntered: "Updated+Entered",
	EntityOpDeletedLeft:    "Deleted+Left",
}

// Flag determines whether an EntityOp includes another. This is primarily
// offered to prevent bit flag errors in downstream clients.
func (o EntityOp) Flag(p EntityOp) bool {
	return o&p != 0
}

// String returns a human identifiable string for the EntityOp
func (o EntityOp) String() string {
	return entityOpNames[o]
}

// EntityHandler is a function that receives Entity updates
type EntityHandler func(Entity, EntityOp) error

// ValueOrZero returns the property's value, or the zero value when the entity is nil or its
// class has no such property (CS2 renames, POV demos, a stale handle naming another class).
func ValueOrZero(entity Entity, name string) PropertyValue {
	if entity == nil {
		return PropertyValue{}
	}

	v, _ := entity.PropertyValue(name)

	return v
}
