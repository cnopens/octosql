package execution

import (
	"context"
	"testing"
	"time"

	"github.com/cube2222/octosql"
	"github.com/cube2222/octosql/docs"
)

func TestOrderBy_Get(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	type args struct {
		stream      RecordStream
		expressions []Expression
		directions  []OrderDirection
	}
	tests := []struct {
		name    string
		args    args
		want    RecordStream
		wantErr bool
	}{
		{
			name: "simple order - one column int ascending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"id", "age"},
						[]interface{}{1, 7}),
					NewRecordFromSliceWithNormalize(
						[]string{"id", "age"},
						[]interface{}{2, 10}),
					NewRecordFromSliceWithNormalize(
						[]string{"id", "age"},
						[]interface{}{3, 2}),
				}),
				expressions: []Expression{NewVariable(octosql.NewVariableName("age"))},
				directions:  []OrderDirection{Ascending},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]string{"id", "age"},
					[]interface{}{3, 2}),
				NewRecordFromSliceWithNormalize(
					[]string{"id", "age"},
					[]interface{}{1, 7}),
				NewRecordFromSliceWithNormalize(
					[]string{"id", "age"},
					[]interface{}{2, 10}),
			}),
			wantErr: false,
		},
		{
			name: "simple order - one column string descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"b", 10}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"c", 7}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", 2}),
				}),
				expressions: []Expression{NewVariable(octosql.NewVariableName("name"))},
				directions:  []OrderDirection{Descending},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"c", 7}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"b", 10}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"a", 2}),
			}),
			wantErr: false,
		},
		{
			name: "simple order - one column time descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"name", "birth"},
						[]interface{}{"b", now}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "birth"},
						[]interface{}{"c", now.Add(time.Hour)}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "birth"},
						[]interface{}{"a", now.Add(-1 * time.Hour)}),
				}),
				expressions: []Expression{NewVariable(octosql.NewVariableName("birth"))},
				directions:  []OrderDirection{Descending},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]string{"name", "birth"},
					[]interface{}{"c", now.Add(time.Hour)}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "birth"},
					[]interface{}{"b", now}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "birth"},
					[]interface{}{"a", now.Add(-1 * time.Hour)}),
			}),
			wantErr: false,
		},
		{
			name: "complex order - string ascending then int descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", 7}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"d", 19}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", -2}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"c", 1}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"d", 17}),
				}),
				expressions: []Expression{
					NewVariable(octosql.NewVariableName("name")),
					NewVariable(octosql.NewVariableName("age")),
				},
				directions: []OrderDirection{
					Ascending,
					Descending,
				},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"a", 7}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"a", -2}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"c", 1}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"d", 19}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"d", 17}),
			}),
			wantErr: false,
		},
		{
			name: "complex order - string ascending then int descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", 7.3}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"d", 19.02}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", -2.248}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"c", 1.123}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"d", 17.918}),
				}),
				expressions: []Expression{
					NewVariable(octosql.NewVariableName("name")),
					NewFunctionExpression(&FuncIdentity, []Expression{NewVariable(octosql.NewVariableName("age"))}),
				},
				directions: []OrderDirection{
					Ascending,
					Descending,
				},
			},

			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"a", 7.3}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"a", -2.248}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"c", 1.123}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"d", 19.02}),
				NewRecordFromSliceWithNormalize(
					[]string{"name", "age"},
					[]interface{}{"d", 17.918}),
			}),
			wantErr: false,
		},
		{
			name: "failed - missing field",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", 7}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age?"},
						[]interface{}{"d", 19}),
				}),
				expressions: []Expression{NewVariable(octosql.NewVariableName("age"))},
				directions:  []OrderDirection{Descending},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "failed - type mismatch",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"a", 7}),
					NewRecordFromSliceWithNormalize(
						[]string{"name", "age"},
						[]interface{}{"d", 19.5}),
				}),
				expressions: []Expression{NewVariable(octosql.NewVariableName("age"))},
				directions:  []OrderDirection{Descending},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordered, err := createOrderedStream(ctx, tt.args.expressions, tt.args.directions, octosql.NoVariables(), tt.args.stream)
			if err != nil && !tt.wantErr {
				t.Errorf("Error in create stream: %v", err)
				return
			} else if err != nil {
				return
			}

			equal, err := AreStreamsEqual(context.Background(), tt.want, ordered)
			if err != nil {
				t.Errorf("Error in AreStreamsEqual(): %v", err)
				return
			}

			if !equal {
				t.Errorf("Streams don't match")
				return
			}
		})
	}
}

type AnyOk struct {
}

func (*AnyOk) Document() docs.Documentation {
	panic("implement me")
}
func (*AnyOk) Validate(args ...octosql.Value) error {
	return nil
}

var FuncIdentity = Function{
	Validator: &AnyOk{},
	Logic: func(args ...octosql.Value) (octosql.Value, error) {
		return args[0], nil
	},
}
