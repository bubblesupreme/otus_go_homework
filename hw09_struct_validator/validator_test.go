package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5|in:01.00,01.01"`
	}

	Token struct {
		Header byte `validate:"len:11"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	ResponseErrParseTag struct {
		Code int    `validate:"in:200,404"`
		Body string `json:"omitempty" validate:"len:11:10"`
	}

	ResponseErrItoa struct {
		Code int    `validate:"in:aa,404"`
		Body string `json:"omitempty"`
	}

	Nested struct {
		App App `validate:"nested"`
	}

	NestedNotTag struct {
		App App `validate:""`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			struct{}{},
			nil,
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566"}, json.RawMessage{}},
			nil,
		},
		{
			User{"01", "Ivan", 20, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"ID", arrayErrors{ErrStringValueLen}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 10, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Age", arrayErrors{ErrIntValueMin}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 51, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Age", arrayErrors{ErrIntValueMax}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@@\\gmail.com", UserRole("stuff"), []string{"89077876566"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Email", arrayErrors{ErrStringValueRegexp}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@gmail.com", UserRole("worker"), []string{"89077876566"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Role", arrayErrors{ErrStringValueSet}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566", "89077876566"}, json.RawMessage{}},
			nil,
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566788"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Phones", ValidationErrors{ValidationError{"0", arrayErrors{ErrStringValueLen}}}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566", "89077876566788"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Phones", ValidationErrors{ValidationError{"1", arrayErrors{ErrStringValueLen}}}}},
		},
		{
			User{"000000000000000000000000000000000001", "Ivan", 20, "ivan78@gmail.com", UserRole("stuff"), []string{"89077876566788", "89077876566788"}, json.RawMessage{}},
			ValidationErrors{ValidationError{"Phones", ValidationErrors{ValidationError{"0", arrayErrors{ErrStringValueLen}}, ValidationError{"1", arrayErrors{ErrStringValueLen}}}}},
		},
		{
			App{"01.01"},
			nil,
		},
		{
			App{"01.0"},
			ValidationErrors{ValidationError{"Version", arrayErrors{ErrStringValueLen, ErrStringValueSet}}},
		},
		{
			Response{300, ""},
			ValidationErrors{ValidationError{"Code", arrayErrors{ErrTagParseTwoInts}}},
		},
		{
			ResponseErrParseTag{300, ""},
			ValidationErrors{ValidationError{"Body", ErrParseTag}},
		},
		{
			ResponseErrItoa{300, ""},
			ValidationErrors{ValidationError{"Code", arrayErrors{errors.New("strconv.ParseInt: parsing \"aa\": invalid syntax")}}},
		},
		{
			Token{'0'},
			ValidationErrors{ValidationError{"Header", ErrTypeNotImplemented}},
		},
		{
			Nested{App{"01.01"}},
			nil,
		},
		{
			Nested{App{"01.02"}},
			ValidationErrors{ValidationError{"App", ValidationErrors{ValidationError{"Version", arrayErrors{ErrStringValueSet}}}}},
		},
		{
			NestedNotTag{App{"01.02"}},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
