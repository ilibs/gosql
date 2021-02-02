package gosql

import "testing"

func Test_mysqlDialect_GetName(t *testing.T) {
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
			want:   "mysql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			my := mysqlDialect{
				commonDialect: tt.fields.commonDialect,
			}
			if got := my.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysqlDialect_Quote(t *testing.T) {
	type fields struct {
		commonDialect commonDialect
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "test",
			fields: fields{},
			args:   args{"status"},
			want:   "`status`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			my := mysqlDialect{
				commonDialect: tt.fields.commonDialect,
			}
			if got := my.Quote(tt.args.key); got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}
