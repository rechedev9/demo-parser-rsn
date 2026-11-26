package demoinfocs

import (
	"testing"

	"github.com/golang/geo/r3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	common "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/common"
	events "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/events"
	st "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/sendtables"
	stfake "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/sendtables/fake"
)

type fakeServerClass struct {
	st.ServerClass
	name string
}

func (sc fakeServerClass) Name() string { return sc.name }

type thrownProjectile struct {
	className string
	model     uint64
	values    map[string]any
}

func throwProjectile(t *testing.T, p *parser, nade thrownProjectile) {
	t.Helper()

	var onCreateFinished func()

	entity := new(stfake.Entity)
	entity.On("ID").Return(100)
	entity.On("ServerClass").Return(fakeServerClass{name: nade.className})
	entity.On("Position").Return(r3.Vector{})
	entity.On("OnDestroy", mock.Anything).Return()
	entity.On("OnCreateFinished", mock.Anything).Run(func(args mock.Arguments) {
		onCreateFinished = args.Get(0).(func())
	})
	entity.On("Property", mock.Anything).Return(nil)
	entity.On("PropertyValue", "CBodyComponent.m_hModel").Return(st.PropertyValue{Any: nade.model}, true)

	for name, value := range nade.values {
		entity.On("PropertyValue", name).Return(st.PropertyValue{Any: value}, true)
	}

	entity.On("PropertyValue", mock.Anything).Return(st.PropertyValue{}, false)

	p.bindGrenadeProjectiles(entity)
	onCreateFinished()
}

func TestBindGrenadeProjectiles_TypeFromModelOrClass(t *testing.T) {
	const learnedModel uint64 = 8528019637059176743

	tests := []struct {
		name string
		nade thrownProjectile
		want common.EquipmentType
	}{
		{"learned model", thrownProjectile{className: "CHEGrenadeProjectile", model: learnedModel}, common.EqSmoke},
		{"model 0", thrownProjectile{className: "CHEGrenadeProjectile", model: 0}, common.EqHE},
		{"model not seen yet", thrownProjectile{className: "CSmokeGrenadeProjectile", model: 12910967535364078964}, common.EqSmoke},
		{"flashbang", thrownProjectile{className: "CFlashbangProjectile", model: 1}, common.EqFlash},
		{"decoy", thrownProjectile{className: "CDecoyProjectile", model: 2531988892557392303}, common.EqDecoy},
		{"molotov", thrownProjectile{className: "CMolotovProjectile", model: 3, values: map[string]any{"m_bIsIncGrenade": false}}, common.EqMolotov},
		{"incendiary", thrownProjectile{className: "CMolotovProjectile", model: 0, values: map[string]any{"m_bIsIncGrenade": true}}, common.EqIncendiary},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := newParser()
			p.equipmentTypePerModel[learnedModel] = common.EqSmoke

			var (
				thrown   common.EquipmentType
				warnings []events.ParserWarn
			)

			p.RegisterEventHandler(func(e events.GrenadeProjectileThrow) { thrown = e.Projectile.WeaponInstance.Type })
			p.RegisterEventHandler(func(e events.ParserWarn) { warnings = append(warnings, e) })

			throwProjectile(t, p, tc.nade)

			assert.Equal(t, tc.want, thrown)
			assert.Empty(t, warnings)
		})
	}
}

func TestBindGrenadeProjectiles_UnknownClassWarns(t *testing.T) {
	p := newParser()

	var warnings []events.ParserWarn

	p.RegisterEventHandler(func(e events.ParserWarn) { warnings = append(warnings, e) })

	throwProjectile(t, p, thrownProjectile{className: "CNewGrenadeProjectile", model: 42})

	assert.Len(t, warnings, 1)
	assert.Equal(t, events.WarnType(events.WarnTypeUnknownGrenadeModel), warnings[0].Type)
}
