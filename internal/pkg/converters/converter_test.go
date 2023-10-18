package converters_test

import (
	"errors"
	"testing"

	"github.com/arieffian/roman-alien-currency/internal/pkg/converters"
	"github.com/go-test/deep"
)

func TestRomanToArabic(t *testing.T) {

	converter := converters.NewConverter()

	type args struct {
		param string
	}

	type want struct {
		result int
		error  error
	}

	testcases := []struct {
		name       string
		args       args
		beforeEach func(*testing.T, *args)
		want       want
	}{
		{
			name: "when there is invalid roman numeral should return error",
			args: args{
				param: "IIII",
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: 0,
				error:  errors.New("invalid roman number"),
			},
		},
		{
			name: "when there is leftover string should return error",
			args: args{
				param: "IIIXM",
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: 0,
				error:  errors.New("invalid roman number"),
			},
		},
		{
			name: "when there is valid roman numeral should return success",
			args: args{
				param: "III",
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: 3,
				error:  nil,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result, err := converter.RomanToArabic(tc.args.param)

			if err != nil || tc.want.error != nil {
				if diff := deep.Equal(err.Error(), tc.want.error.Error()); diff != nil {
					t.Errorf("got unexpected error.\n expect: %v\n actual: %v\n diff: %v\n", tc.want.error, err, diff)
				}
			}

			if diff := deep.Equal(result, tc.want.result); diff != nil {
				t.Errorf("got unexpected result.\n expected: %v\n actual: %v\n diff: %v\n", tc.want.result, result, diff)
			}
		})

	}
}

func TestArabicToRoman(t *testing.T) {
	converter := converters.NewConverter()

	type args struct {
		param int
	}

	type want struct {
		result string
		error  error
	}

	testcases := []struct {
		name       string
		args       args
		beforeEach func(*testing.T, *args)
		want       want
	}{
		{
			name: "when there is number greater than 3999 should return error",
			args: args{
				param: 4000,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: "",
				error:  errors.New("number out of range"),
			},
		},
		{
			name: "when there is number less than 1 should return error",
			args: args{
				param: 0,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: "",
				error:  errors.New("number out of range"),
			},
		},
		{
			name: "when there is number between 1 and 3999 should return success",
			args: args{
				param: 30,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: "XXX",
				error:  nil,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result, err := converter.ArabicToRoman(tc.args.param)

			if err != nil || tc.want.error != nil {
				if diff := deep.Equal(err.Error(), tc.want.error.Error()); diff != nil {
					t.Errorf("got unexpected error.\n expect: %v\n actual: %v\n diff: %v\n", tc.want.error, err, diff)
				}
			}

			if diff := deep.Equal(result, tc.want.result); diff != nil {
				t.Errorf("got unexpected result.\n expected: %v\n actual: %v\n diff: %v\n", tc.want.result, result, diff)
			}
		})

	}
}

func TestAlienToRoman(t *testing.T) {
	converter := converters.NewConverter()

	alienDict := map[string]string{
		"glob": "I",
	}

	type args struct {
		dict        map[string]string
		alienNumber []string
	}

	type want struct {
		result string
		error  error
	}

	testcases := []struct {
		name       string
		args       args
		beforeEach func(*testing.T, *args)
		want       want
	}{
		{
			name: "when input is valid should return success",
			args: args{
				alienNumber: []string{"glob", "glob", "glob"},
				dict:        alienDict,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: "III",
				error:  nil,
			},
		},
		{
			name: "when input is invalid should return success",
			args: args{
				alienNumber: []string{"glob", "prok"},
				dict:        alienDict,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: "",
				error:  errors.New("invalid alien number"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result, err := converter.AlienToRoman(tc.args.dict, tc.args.alienNumber)

			if err != nil || tc.want.error != nil {
				if diff := deep.Equal(err.Error(), tc.want.error.Error()); diff != nil {
					t.Errorf("got unexpected error.\n expect: %v\n actual: %v\n diff: %v\n", tc.want.error, err, diff)
				}
			}

			if diff := deep.Equal(result, tc.want.result); diff != nil {
				t.Errorf("got unexpected result.\n expected: %v\n actual: %v\n diff: %v\n", tc.want.result, result, diff)
			}
		})

	}
}
