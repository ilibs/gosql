package gosql

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJsonObject(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				value: nil,
			},
			want:    emptyJSON,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				value: []byte(""),
			},
			want:    emptyJSON,
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				value: "",
			},
			want:    emptyJSON,
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				value: `{"other":1}`,
			},
			want:    []byte(`{"other":1}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonObject(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonValue() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestJSONText(t *testing.T) {
	j := JSONText(`{"foo": 1, "bar": 2}`)
	v, err := j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error")
	}
	err = (&j).Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error")
	}
	m := map[string]interface{}{}
	j.Unmarshal(&m)

	if m["foo"].(float64) != 1 || m["bar"].(float64) != 2 {
		t.Errorf("Expected valid json but got some garbage instead? %#v", m)
	}

	j = JSONText(`{"foo": 1, invalid, false}`)
	v, err = j.Value()
	if err == nil {
		t.Errorf("Was expecting invalid json to fail!")
	}

	j = JSONText("")
	v, err = j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error")
	}

	err = (&j).Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error")
	}

	j = JSONText(nil)
	v, err = j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error")
	}

	err = (&j).Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error")
	}

	t.Run("Binary", func(t *testing.T) {
		j := JSONText(`{"foo": 1, "bar": 2}`)
		v, err := j.MarshalBinary()
		if err != nil {
			t.Errorf("Was not expecting an error")
		}
		if string(v) != `{"foo": 1, "bar": 2}` {
			t.Errorf("MarshalBinary result error")
		}

		err = (&j).UnmarshalBinary(v)
		if err != nil {
			t.Errorf("Was not expecting an error")
		}
	})
}
