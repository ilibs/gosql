package gosql

import "testing"

func Test_postgresDialect_GetName(t *testing.T) {
	type fields struct {
		commonDialect commonDialect
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test",
			fields: fields{},
			want:   "postgres",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			po := postgresDialect{
				commonDialect: tt.fields.commonDialect,
			}
			if got := po.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
