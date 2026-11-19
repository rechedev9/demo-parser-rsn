package sendtables_test

import (
	"testing"

	"github.com/golang/geo/r3"
	"github.com/stretchr/testify/assert"

	"github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/constants"
	st "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/sendtables"
	stfake "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/sendtables/fake"
)

func TestPropertyValueWithoutValueReadsZero(t *testing.T) {
	v := st.PropertyValue{}

	assert.NotPanics(t, func() {
		assert.Zero(t, v.Int())
		assert.Zero(t, v.Int64())
		assert.Zero(t, v.UInt32())
		assert.Zero(t, v.UInt64())
		assert.Zero(t, v.Float())
		assert.False(t, v.BoolVal())
		assert.Empty(t, v.Str())
		assert.Nil(t, v.Array())
		assert.Equal(t, r3.Vector{}, v.R3Vec())
		assert.Nil(t, v.R3VecOrNil())
	})
	assert.Equal(t, uint64(constants.InvalidEntityHandleSource2), v.Handle())
}

func TestPropertyValueAcceptsEveryIntegerWidth(t *testing.T) {
	for _, any := range []any{int32(7), int64(7), uint32(7), uint64(7)} {
		v := st.PropertyValue{Any: any}

		assert.Equal(t, 7, v.Int(), "%T", any)
		assert.Equal(t, uint64(7), v.UInt64(), "%T", any)
		assert.Equal(t, uint64(7), v.Handle(), "%T", any)
		assert.True(t, v.BoolVal(), "%T", any)
	}

	assert.Equal(t, -3, st.PropertyValue{Any: int32(-3)}.Int())
	assert.Equal(t, int64(-3), st.PropertyValue{Any: int32(-3)}.Int64())
	assert.False(t, st.PropertyValue{Any: uint32(0)}.BoolVal())
}

func TestPropertyValueShortVectorReadsZero(t *testing.T) {
	v := st.PropertyValue{Any: []float32{1, 2}}

	assert.Equal(t, r3.Vector{}, v.R3Vec())
	assert.Nil(t, v.R3VecOrNil())
	assert.Equal(t, r3.Vector{X: 1, Y: 2, Z: 3}, st.PropertyValue{Any: []float32{1, 2, 3}}.R3Vec())
}

func TestValueOrZeroWithoutEntityOrProperty(t *testing.T) {
	assert.Equal(t, st.PropertyValue{}, st.ValueOrZero(nil, "m_iHealth"))

	entity := new(stfake.Entity)
	entity.On("PropertyValue", "m_iHealth").Return(st.PropertyValue{}, false)
	assert.Zero(t, st.ValueOrZero(entity, "m_iHealth").Int())
}
