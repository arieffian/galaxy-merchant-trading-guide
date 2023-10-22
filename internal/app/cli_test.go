package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/arieffian/roman-alien-currency/internal/app"
	mockConverter "github.com/arieffian/roman-alien-currency/internal/pkg/converters/mocks"
	mockParser "github.com/arieffian/roman-alien-currency/internal/pkg/parsers/mocks"
	mockReader "github.com/arieffian/roman-alien-currency/internal/pkg/readers/mocks"
	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
)

func TestCLI(t *testing.T) {
	ctrl := gomock.NewController(t)
	converter := mockConverter.NewMockConverterService(ctrl)
	parser := mockParser.NewMockParserService(ctrl)
	fileReader := mockReader.NewMockFileService(ctrl)

	cli, _ := app.NewCli(app.NewCliParams{
		Converter:  converter,
		Parser:     parser,
		FileReader: fileReader,
	})

	ctx := context.Background()

	validResult := []string{
		"glob is I",
	}

	type args struct {
		param context.Context
	}

	type want struct {
		error error
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
				param: ctx,
			},
			beforeEach: func(t *testing.T, a *args) {
				fileReader.
					EXPECT().
					ReadFile(gomock.Any()).
					Return(validResult, nil)

				parser.
					EXPECT().
					FixTypo(gomock.Any()).
					Return("glob is I")

				parser.
					EXPECT().
					ParseCurrency(gomock.Any()).
					Return(false)

				parser.
					EXPECT().
					ParseMetal(gomock.Any()).
					Return(false, nil)

				parser.
					EXPECT().
					ProcessQuestion(gomock.Any()).
					Return([]string{}, nil)
			},
			want: want{
				error: nil,
			},
		},
		{
			name: "when filereader is error should return error",
			args: args{
				param: ctx,
			},
			beforeEach: func(t *testing.T, a *args) {
				fileReader.
					EXPECT().
					ReadFile("input").
					Return(nil, errors.New("error"))

			},
			want: want{
				error: errors.New("error"),
			},
		},
		{
			name: "when parser is error should return error",
			args: args{
				param: ctx,
			},
			beforeEach: func(t *testing.T, a *args) {
				fileReader.
					EXPECT().
					ReadFile(gomock.Any()).
					Return(validResult, nil)

				parser.
					EXPECT().
					FixTypo(gomock.Any()).
					Return("glob is I")

				parser.
					EXPECT().
					ParseCurrency(gomock.Any()).
					Return(false)

				parser.
					EXPECT().
					ParseMetal(gomock.Any()).
					Return(false, errors.New("error"))
			},
			want: want{
				error: errors.New("error"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tc.beforeEach(t, &tc.args)

			err := cli.Run(tc.args.param)

			if err != nil || tc.want.error != nil {
				if diff := deep.Equal(err.Error(), tc.want.error.Error()); diff != nil {
					t.Errorf("got unexpected error.\n expect: %v\n actual: %v\n diff: %v\n", tc.want.error, err, diff)
				}
			}
		})
	}
}
