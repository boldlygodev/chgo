package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parse(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}

	tests := []struct {
		name    string
		args    args
		want    version
		wantErr bool
	}{
		{
			name:    "tip",
			args:    args{"gotip"},
			want:    version{},
			wantErr: false,
		},
		{
			name:    "minor",
			args:    args{"go1.2"},
			want:    version{Minor: 2, Stage: stageRelease},
			wantErr: false,
		},
		{
			name:    "patch",
			args:    args{"go1.2.3"},
			want:    version{Minor: 2, Patch: 3, Stage: stageRelease},
			wantErr: false,
		},
		{
			name:    "rc",
			args:    args{"go1.2rc4"},
			want:    version{Minor: 2, Stage: stageCandidate, Prerelease: 4},
			wantErr: false,
		},
		{
			name:    "beta",
			args:    args{"go1.2beta4"},
			want:    version{Minor: 2, Stage: stageBeta, Prerelease: 4},
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{},
			want:    version{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_version_String(t *testing.T) {
	t.Parallel()

	type fields struct {
		Minor      uint8
		Patch      uint8
		Stage      uint8
		Prerelease uint8
	}

	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name:   "tip",
			fields: fields{},
			want:   "gotip",
		},
		{
			name:   "minor",
			fields: fields{Minor: 2, Stage: stageRelease},
			want:   "go1.2",
		},
		{
			name:   "patch",
			fields: fields{Minor: 2, Patch: 3, Stage: stageRelease},
			want:   "go1.2.3",
		},
		{
			name:   "rc",
			fields: fields{Minor: 2, Stage: stageCandidate, Prerelease: 4},
			want:   "go1.2rc4",
		},
		{
			name:   "beta",
			fields: fields{Minor: 2, Stage: stageBeta, Prerelease: 4},
			want:   "go1.2beta4",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := version{
				Minor:      tt.fields.Minor,
				Patch:      tt.fields.Patch,
				Stage:      tt.fields.Stage,
				Prerelease: tt.fields.Prerelease,
			}

			got := v.String()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("version.String() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_version_MarshalText(t *testing.T) {
	t.Parallel()

	type fields struct {
		Minor      uint8
		Patch      uint8
		Stage      uint8
		Prerelease uint8
	}

	tests := []struct {
		name    string
		want    []byte
		fields  fields
		wantErr bool
	}{
		{
			name:    "tip",
			fields:  fields{},
			want:    []byte("gotip"),
			wantErr: false,
		},
		{
			name:    "minor",
			fields:  fields{Minor: 2, Stage: stageRelease},
			want:    []byte("go1.2"),
			wantErr: false,
		},
		{
			name:    "patch",
			fields:  fields{Minor: 2, Patch: 3, Stage: stageRelease},
			want:    []byte("go1.2.3"),
			wantErr: false,
		},
		{
			name:    "rc",
			fields:  fields{Minor: 2, Stage: stageCandidate, Prerelease: 4},
			want:    []byte("go1.2rc4"),
			wantErr: false,
		},
		{
			name:    "beta",
			fields:  fields{Minor: 2, Stage: stageBeta, Prerelease: 4},
			want:    []byte("go1.2beta4"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := version{
				Minor:      tt.fields.Minor,
				Patch:      tt.fields.Patch,
				Stage:      tt.fields.Stage,
				Prerelease: tt.fields.Prerelease,
			}

			got, err := v.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Fatalf("version.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("version.MarshalText() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_version_UnmarshalText(t *testing.T) {
	t.Parallel()

	type fields struct {
		Minor      uint8
		Patch      uint8
		Stage      uint8
		Prerelease uint8
	}

	type args struct {
		data []byte
	}

	tests := []struct {
		name    string
		args    args
		fields  fields
		want    version
		wantErr bool
	}{
		{
			name:    "tip",
			fields:  fields{},
			args:    args{[]byte("gotip")},
			want:    version{},
			wantErr: false,
		},
		{
			name:    "minor",
			fields:  fields{},
			args:    args{[]byte("go1.2")},
			want:    version{Minor: 2, Stage: stageRelease},
			wantErr: false,
		},
		{
			name:    "patch",
			fields:  fields{},
			args:    args{[]byte("go1.2.3")},
			want:    version{Minor: 2, Patch: 3, Stage: stageRelease},
			wantErr: false,
		},
		{
			name:    "rc",
			fields:  fields{},
			args:    args{[]byte("go1.2rc4")},
			want:    version{Minor: 2, Stage: stageCandidate, Prerelease: 4},
			wantErr: false,
		},
		{
			name:    "beta",
			fields:  fields{},
			args:    args{[]byte("go1.2beta4")},
			want:    version{Minor: 2, Stage: stageBeta, Prerelease: 4},
			wantErr: false,
		},
		{
			name:    "invalid",
			fields:  fields{},
			args:    args{},
			want:    version{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := version{
				Minor:      tt.fields.Minor,
				Patch:      tt.fields.Patch,
				Stage:      tt.fields.Stage,
				Prerelease: tt.fields.Patch,
			}
			if err := got.UnmarshalText(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("version.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("version.UnmarshalText() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_versions_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    versions
		want int
	}{
		{
			name: "empty",
			v:    versions{},
			want: 0,
		},
		{
			name: "one",
			v:    versions{{}},
			want: 1,
		},
		{
			name: "two",
			v:    versions{{}, {}},
			want: 2,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.v.Len()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("versions.Len() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_versions_Less(t *testing.T) {
	t.Parallel()

	type args struct {
		i int
		j int
	}

	tests := []struct {
		name string
		v    versions
		args args
		want bool
	}{
		{
			name: "tip same",
			v:    versions{{}},
			args: args{0, 0},
			want: false,
		},
		{
			name: "tip less",
			v:    versions{{}, {Minor: 1}},
			args: args{0, 1},
			want: false,
		},
		{
			name: "tip more",
			v:    versions{{}, {Minor: 1}},
			args: args{1, 0},
			want: true,
		},
		{
			name: "minor same",
			v:    versions{{Minor: 1}},
			args: args{0, 0},
			want: false,
		},
		{
			name: "minor less",
			v:    versions{{Minor: 1}, {Minor: 2}},
			args: args{0, 1},
			want: true,
		},
		{
			name: "minor more",
			v:    versions{{Minor: 1}, {Minor: 2}},
			args: args{1, 0},
			want: false,
		},
		{
			name: "patch same",
			v:    versions{{Patch: 1}},
			args: args{0, 0},
			want: false,
		},
		{
			name: "patch less",
			v:    versions{{Patch: 1}, {Patch: 2}},
			args: args{0, 1},
			want: true,
		},
		{
			name: "patch more",
			v:    versions{{Patch: 1}, {Patch: 2}},
			args: args{1, 0},
			want: false,
		},
		{
			name: "stage same",
			v:    versions{{Stage: 1}},
			args: args{0, 0},
			want: false,
		},
		{
			name: "stage less",
			v:    versions{{Stage: 1}, {Stage: 2}},
			args: args{0, 1},
			want: true,
		},
		{
			name: "stage more",
			v:    versions{{Stage: 1}, {Stage: 2}},
			args: args{1, 0},
			want: false,
		},
		{
			name: "prerelease same",
			v:    versions{{Prerelease: 1}},
			args: args{0, 0},
			want: false,
		},
		{
			name: "prerelease less",
			v:    versions{{Prerelease: 1}, {Prerelease: 2}},
			args: args{0, 1},
			want: true,
		},
		{
			name: "prerelease more",
			v:    versions{{Prerelease: 1}, {Prerelease: 2}},
			args: args{1, 0},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.v.Less(tt.args.i, tt.args.j)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("version.Less() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_versions_Swap(t *testing.T) {
	t.Parallel()

	type args struct {
		i int
		j int
	}

	tests := []struct {
		name string
		v    versions
		want versions
		args args
	}{
		{
			v:    versions{{Minor: 1}, {Minor: 2}, {Minor: 3}, {Minor: 4}},
			args: args{1, 2},
			want: versions{{Minor: 1}, {Minor: 3}, {Minor: 2}, {Minor: 4}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.v.Swap(tt.args.i, tt.args.j)

			got := tt.v

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("version.Swap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
