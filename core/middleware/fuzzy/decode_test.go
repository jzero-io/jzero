package fuzzy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CreateRequest struct {
	Base CreateBase `json:"base"` // 用户基本信息
}

type CreateBase struct {
	UserBase
	Username string `json:"username"` // 用户账号
	Password string `json:"password"` // 登陆密码
}

type UserBase struct {
	Nickname string `json:"nickname" olabel:"用户姓名"`                                            // 用户名称
	GroupId  int    `json:"groupId" olabel:"所属组织id"`                                           // 组织 id
	Phone    string `json:"phone" olabel:"手机号"`                                                // 手机号
	Email    string `json:"email" olabel:"邮箱"`                                                 // 邮箱
	Status   int8   `json:"status,optional" olabel:"状态" ovalue:"0=正常,1=已锁定,2=已失效,3=已注销,4=已禁用"` // 用户状态 0-正常 1-已锁定 2-已失效 3-已注销 4-已禁用
}

func TestFuzzyDecodeRequest(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		outputType  any
		expectError bool
	}{
		{
			name:       "simple struct with string field",
			input:      []byte(`{"name": "John"}`),
			outputType: &struct{ Name string }{},
		},
		{
			name:       "struct with int field from string",
			input:      []byte(`{"age": "18"}`),
			outputType: &struct{ Age int }{},
		},

		{
			name:       "struct with bool field from string",
			input:      []byte(`{"active": "true"}`),
			outputType: &struct{ Active bool }{},
		},
		{
			name:       "struct with bool field from bool",
			input:      []byte(`{"active": true}`),
			outputType: &struct{ Active bool }{},
		},
		{
			name:       "struct with pointer bool field from bool",
			input:      []byte(`{"active": true}`),
			outputType: &struct{ Active *bool }{},
		},
		{
			name:       "struct with *bool field from string",
			input:      []byte(`{"active": "true"}`),
			outputType: &struct{ Active *bool }{},
		},

		{
			name:       "struct with float field from string",
			input:      []byte(`{"price": "19.99"}`),
			outputType: &struct{ Price float64 }{},
		},
		{
			name:  "nested struct with fuzzy values",
			input: []byte(`{"user": {"id": "123", "premium": "1"}}`),
			outputType: &struct {
				User struct {
					ID      int
					Premium bool
				}
			}{},
		},
		{
			name:       "empty input",
			input:      []byte(`{}`),
			outputType: &struct{}{},
		},
		{
			name:       "pointer fields",
			input:      []byte(`{"ptr": "123"}`),
			outputType: &struct{ Ptr *int }{},
		},
		{
			name: "embeded struct",
			input: []byte(`{
  "base": {
    "username": "test12@myibc.net",
    "nickname": "",
    "groupId": 1,
    "phone": "",
    "email": "",
    "password": "0416f839f229cf221ebc4667c9839f90c687aef23cff6f885375f8ab6a506a19e191620e86ec517e61a5cee092298bf340934d7f70328b61e9559e0e9ec2849d8b801D21CB1459F0A1364377270A0BF8620732658BBB0958A1C0A1214EC5BD211ECB0757A3D5DB981090"
  }
}`),
			outputType: &CreateRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decode, err := Decode(tt.input, tt.outputType)
			if err != nil {
				t.Errorf("Decode error: %v", err.Error())
				return
			}
			t.Log(tt.outputType)
			err = json.Unmarshal(decode, tt.outputType)
			assert.NoError(t, err)
		})
	}
}

func TestFuzzyDecodeRequest_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		outputType any
		check      func(t *testing.T, output any)
	}{
		{
			name:       "null to zero value",
			input:      []byte(`{"value": null}`),
			outputType: &struct{ Value int }{},
			check: func(t *testing.T, output any) {
				t.Helper()
				assert.Equal(t, 0, output.(*struct{ Value int }).Value)
			},
		},
		{
			name:       "empty string to zero value",
			input:      []byte(`{"value": ""}`),
			outputType: &struct{ Value int }{},
			check: func(t *testing.T, output any) {
				t.Helper()
				assert.Equal(t, 0, output.(*struct{ Value int }).Value)
			},
		},
		{
			name:  "string true/false to bool",
			input: []byte(`{"trueVal": "true", "falseVal": "false"}`),
			outputType: &struct {
				TrueVal  bool
				FalseVal bool
			}{},
			check: func(t *testing.T, output any) {
				t.Helper()
				out := output.(*struct {
					TrueVal  bool
					FalseVal bool
				})
				assert.True(t, out.TrueVal)
				assert.False(t, out.FalseVal)
			},
		},
		{
			name:  "string numbers to bool",
			input: []byte(`{"trueVal": "1", "falseVal": "0"}`),
			outputType: &struct {
				TrueVal  bool
				FalseVal bool
			}{},
			check: func(t *testing.T, output any) {
				t.Helper()
				out := output.(*struct {
					TrueVal  bool
					FalseVal bool
				})
				assert.True(t, out.TrueVal)
				assert.False(t, out.FalseVal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.outputType
			_, err := Decode(tt.input, output)
			assert.NoError(t, err)
			tt.check(t, output)
		})
	}
}
