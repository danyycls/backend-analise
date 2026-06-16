package testkit

import (
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
)

type MockConfig struct {
	Target string
	Method string
	Args   []any
	Return []any
	Times  int
}

type TestCase struct {
	Name         string
	Input        any
	WantOutput   any
	Err          bool
	WantStatus   int
	MockConfigs  []MockConfig
	CustomAssert func(t *testing.T, got any, err error)
}

type MockFactory func(target string, ctrl *gomock.Controller) any

type MockGetter func(target string) any

type ExecFunc[I any] func(input I, getMock MockGetter) (any, error)

func RunTestCase[I any](t *testing.T, tc TestCase, factory MockFactory, exec ExecFunc[I]) {
	t.Run(tc.Name, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mocks := map[string]any{}

		for _, mc := range tc.MockConfigs {
			mock := factory(mc.Target, ctrl)
			require.NotNil(t, mock, "mock %s nao encontrado na factory", mc.Target)
			mocks[mc.Target] = mock
			setupMock(t, mock, mc)
		}

		getter := func(target string) any {
			if m, ok := mocks[target]; ok {
				return m
			}
			m := factory(target, ctrl)
			require.NotNil(t, m, "mock %s nao encontrado na factory", target)
			mocks[target] = m
			return m
		}

		var input I
		if tc.Input != nil {
			var ok bool
			input, ok = tc.Input.(I)
			require.True(t, ok, "tipo de Input nao corresponde: %T", tc.Input)
		}

		got, err := exec(input, getter)

		if tc.CustomAssert != nil {
			tc.CustomAssert(t, got, err)
			return
		}
		if tc.Err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		if tc.WantOutput != nil {
			require.Equal(t, tc.WantOutput, got)
		}
	})
}

func setupMock(t *testing.T, mock any, mc MockConfig) {
	t.Helper()

	method := reflect.ValueOf(mock).MethodByName(mc.Method)
	require.True(t, method.IsValid(), "metodo %s nao encontrado no mock %T", mc.Method, mock)

	expectMethod := reflect.ValueOf(mock).MethodByName("EXPECT")
	require.True(t, expectMethod.IsValid(), "EXPECT nao encontrado em %T", mock)

	expect := expectMethod.Call(nil)[0]
	call := expect.MethodByName(mc.Method)
	require.True(t, call.IsValid(), "EXPECT().%s nao encontrado", mc.Method)

	in := make([]reflect.Value, len(mc.Args))
	for i, a := range mc.Args {
		in[i] = reflect.ValueOf(a)
	}
	ret := call.Call(in)[0]

	retMethod := ret.MethodByName("Return")
	require.True(t, retMethod.IsValid(), "Return nao encontrado em %T", ret.Interface())

	var ifaceTyp = reflect.TypeOf((*interface{})(nil)).Elem()

	retIn := make([]reflect.Value, len(mc.Return))
	for i, r := range mc.Return {
		if r == nil {
			retIn[i] = reflect.Zero(ifaceTyp)
		} else {
			retIn[i] = reflect.ValueOf(r)
		}
	}
	retMethod.Call(retIn)

	if mc.Times > 0 {
		times := ret.MethodByName("Times")
		require.True(t, times.IsValid(), "Times nao encontrado")
		times.Call([]reflect.Value{reflect.ValueOf(mc.Times)})
	} else {
		anyTimes := ret.MethodByName("AnyTimes")
		require.True(t, anyTimes.IsValid(), "AnyTimes nao encontrado")
		anyTimes.Call(nil)
	}
}
