package parsers_test

import (
	"errors"
	"testing"

	mockConverter "github.com/arieffian/roman-alien-currency/internal/pkg/converters/mocks"
	"github.com/arieffian/roman-alien-currency/internal/pkg/parsers"
	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
)

func TestParseCurrency(t *testing.T) {

	ctrl := gomock.NewController(t)
	converter := mockConverter.NewMockConverterService(ctrl)
	parser := parsers.NewParser(parsers.NewParserParams{
		Converter:       converter,
		AlienDictionary: map[string]string{},
		MetalValue:      map[string]float64{},
	})

	acceptedParams := []string{
		"glob", "is", "i",
	}

	notAcceptedParams := []string{
		"glob", "prok", "gold", "is", "57800", "credits",
	}

	type args struct {
		param []string
	}

	type want struct {
		result bool
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
				param: acceptedParams,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: true,
			},
		},
		{
			name: "when input is invalid should return success",
			args: args{
				param: notAcceptedParams,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: false,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result := parser.ParseCurrency(tc.args.param)

			if diff := deep.Equal(result, tc.want.result); diff != nil {
				t.Errorf("got unexpected result.\n expected: %v\n actual: %v\n diff: %v\n", tc.want.result, result, diff)
			}
		})

	}
}

func TestGetCurrencyValue(t *testing.T) {

	ctrl := gomock.NewController(t)
	converter := mockConverter.NewMockConverterService(ctrl)
	parser := parsers.NewParser(parsers.NewParserParams{
		Converter:       converter,
		AlienDictionary: map[string]string{},
		MetalValue:      map[string]float64{},
	})

	validParam := []string{
		"glob", "glob",
	}

	invalidRomanParam := []string{
		"glob", "glob", "glob", "glob",
	}

	invalidParam := []string{
		"glob", "prok",
	}

	type args struct {
		param []string
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
			name: "when input is valid should return success",
			args: args{
				param: validParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), validParam).
					Return("II", nil)

				converter.
					EXPECT().
					RomanToArabic("II").
					Return(2, nil)

			},
			want: want{
				result: 2,
				error:  nil,
			},
		},
		{
			name: "when input is invalid should return success",
			args: args{
				param: invalidParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), invalidParam).
					Return("", errors.New("invalid alien number"))
			},
			want: want{
				result: 0,
				error:  errors.New("invalid alien number"),
			},
		},
		{
			name: "when roman input is invalid should return success",
			args: args{
				param: invalidRomanParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), invalidRomanParam).
					Return("IIII", nil)

				converter.
					EXPECT().
					RomanToArabic("IIII").
					Return(0, errors.New("invalid roman number"))
			},
			want: want{
				result: 0,
				error:  errors.New("invalid roman number"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result, err := parser.GetCurrencyValue(tc.args.param)

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

func TestParserMetal(t *testing.T) {

	ctrl := gomock.NewController(t)
	converter := mockConverter.NewMockConverterService(ctrl)
	parser := parsers.NewParser(parsers.NewParserParams{
		Converter: converter,
		AlienDictionary: map[string]string{
			"glob": "i",
		},
		MetalValue: map[string]float64{},
	})

	validParam := []string{
		"glob", "gold", "is", "57800", "credits",
	}

	invalidAlienCurrencyParam := []string{
		"glob", "prok", "gold", "is", "57800", "credits",
	}

	invalidCreditsParam := []string{
		"glob", "gold", "is", "57800a", "credits",
	}

	invalidRomanParam := []string{
		"glob", "glob", "glob", "glob", "gold", "is", "57800", "credits",
	}

	type args struct {
		param []string
	}

	type want struct {
		result bool
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
				param: validParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("I", nil)

				converter.
					EXPECT().
					RomanToArabic(gomock.Any()).
					Return(1, nil)

			},
			want: want{
				result: true,
				error:  nil,
			},
		},
		{
			name: "when roman is invalid should return error",
			args: args{
				param: invalidRomanParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("IIII", nil)

				converter.
					EXPECT().
					RomanToArabic("IIII").
					Return(0, errors.New("invalid roman number"))

			},
			want: want{
				result: false,
				error:  errors.New("invalid roman number"),
			},
		},
		{
			name: "when alien currency is invalid should return error",
			args: args{
				param: invalidAlienCurrencyParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("", errors.New("invalid alien number"))

			},
			want: want{
				result: false,
				error:  errors.New("invalid alien number"),
			},
		},
		{
			name: "when credit is invalid should return error",
			args: args{
				param: invalidCreditsParam,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: false,
				error:  errors.New(`strconv.Atoi: parsing "57800a": invalid syntax`),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result, err := parser.ParseMetal(tc.args.param)

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

func TestProcessQuestion(t *testing.T) {

	ctrl := gomock.NewController(t)
	converter := mockConverter.NewMockConverterService(ctrl)
	parser := parsers.NewParser(parsers.NewParserParams{
		Converter: converter,
		AlienDictionary: map[string]string{
			"glob": "I",
		},
		MetalValue: map[string]float64{
			"Gold": 100,
		},
	})

	validHowMuchParam := []string{
		"how much is glob glob ?",
	}

	invalidHowMuchParam := []string{
		"how much is glob glob prok ?",
	}

	validHowManyParam := []string{
		"how many Credits is glob glob Gold ?",
	}

	invalidHowManyParam := []string{
		"how many Credits is glob prok Gold ?",
	}

	validDoesParam := []string{
		"does glob glob Gold has more Credits than glob Gold ?",
	}

	validIsParam := []string{
		"is glob larger than glob glob ?",
	}

	type args struct {
		param []string
	}

	type want struct {
		result []string
		error  error
	}

	testcases := []struct {
		name       string
		args       args
		beforeEach func(*testing.T, *args)
		want       want
	}{
		{
			name: "when how much input is valid should return success",
			args: args{
				param: validHowMuchParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("II", nil)

				converter.
					EXPECT().
					RomanToArabic("II").
					Return(2, nil)
			},
			want: want{
				result: []string{"glob glob is 2"},
				error:  nil,
			},
		},
		{
			name: "when how much input is invalid should return error",
			args: args{
				param: invalidHowMuchParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("", errors.New("invalid alien number"))
			},
			want: want{
				result: []string{"I have no idea what you are talking about"},
				error:  nil,
			},
		},
		{
			name: "when how many input is valid should return success",
			args: args{
				param: validHowManyParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("II", nil)

				converter.
					EXPECT().
					RomanToArabic("II").
					Return(2, nil)
			},
			want: want{
				result: []string{"glob glob Gold is 200 Credits"},
				error:  nil,
			},
		},
		{
			name: "when how many input is invalid should return error",
			args: args{
				param: invalidHowManyParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("", errors.New("invalid alien number"))
			},
			want: want{
				result: []string{"I have no idea what you are talking about"},
				error:  nil,
			},
		},
		{
			name: "when does input is valid should return success",
			args: args{
				param: validDoesParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				firstAlienToRoman := converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("II", nil)

				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("I", nil).After(firstAlienToRoman)

				firstRomanToArabic := converter.
					EXPECT().
					RomanToArabic("II").
					Return(2, nil)

				converter.
					EXPECT().
					RomanToArabic("I").
					Return(1, nil).After(firstRomanToArabic)
			},
			want: want{
				result: []string{"glob glob Gold has more Credits than glob Gold"},
				error:  nil,
			},
		},
		{
			name: "when is input is valid should return success",
			args: args{
				param: validIsParam,
			},
			beforeEach: func(t *testing.T, a *args) {
				firstAlienToRoman := converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("I", nil)

				converter.
					EXPECT().
					AlienToRoman(gomock.Any(), gomock.Any()).
					Return("II", nil).After(firstAlienToRoman)

				firstRomanToArabic := converter.
					EXPECT().
					RomanToArabic("I").
					Return(1, nil)

				converter.
					EXPECT().
					RomanToArabic("II").
					Return(2, nil).After(firstRomanToArabic)
			},
			want: want{
				result: []string{"glob is smaller than glob glob"},
				error:  nil,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result, err := parser.ProcessQuestion(tc.args.param)

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

func TestFixTypo(t *testing.T) {

	ctrl := gomock.NewController(t)
	converter := mockConverter.NewMockConverterService(ctrl)
	parser := parsers.NewParser(parsers.NewParserParams{
		Converter:       converter,
		AlienDictionary: map[string]string{},
		MetalValue:      map[string]float64{},
	})

	acceptedParams := "istegj glob glob smaller than glob prok?"

	// notAcceptedParams := []string{
	// 	"glob", "prok", "gold", "is", "57800", "credits",
	// }

	type args struct {
		param string
	}

	type want struct {
		result string
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
				param: acceptedParams,
			},
			beforeEach: func(t *testing.T, a *args) {},
			want: want{
				result: "is tegj glob glob smaller than glob prok ?",
			},
		},
		// {
		// 	name: "when input is invalid should return success",
		// 	args: args{
		// 		param: notAcceptedParams,
		// 	},
		// 	beforeEach: func(t *testing.T, a *args) {},
		// 	want: want{
		// 		result: false,
		// 	},
		// },
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			result := parser.FixTypo(tc.args.param)

			if diff := deep.Equal(result, tc.want.result); diff != nil {
				t.Errorf("got unexpected result.\n expected: %v\n actual: %v\n diff: %v\n", tc.want.result, result, diff)
			}
		})

	}
}
