package gosql

import "testing"

func Test_sqlite3Dialect_GetName(t *testing.T) {
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
			want:   "sqlite3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := sqlite3Dialect{
				commonDialect: tt.fields.commonDialect,
			}
			if got := sq.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
