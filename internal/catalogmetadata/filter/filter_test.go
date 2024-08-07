package filter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/operator-framework/operator-registry/alpha/declcfg"

	"github.com/operator-framework/operator-controller/internal/catalogmetadata/filter"
)

func TestFilter(t *testing.T) {
	for _, tt := range []struct {
		name      string
		predicate filter.Predicate[declcfg.Bundle]
		want      []declcfg.Bundle
	}{
		{
			name: "simple filter with one predicate",
			predicate: func(bundle declcfg.Bundle) bool {
				return bundle.Name == "extension1.v1"
			},
			want: []declcfg.Bundle{
				{Name: "extension1.v1", Package: "extension1", Image: "fake1"},
			},
		},
		{
			name: "filter with Not predicate",
			predicate: filter.Not(func(bundle declcfg.Bundle) bool {
				return bundle.Name == "extension1.v1"
			}),
			want: []declcfg.Bundle{
				{Name: "extension1.v2", Package: "extension1", Image: "fake2"},
				{Name: "extension2.v1", Package: "extension2", Image: "fake1"},
			},
		},
		{
			name: "filter with And predicate",
			predicate: filter.And(
				func(bundle declcfg.Bundle) bool {
					return bundle.Name == "extension1.v1"
				},
				func(bundle declcfg.Bundle) bool {
					return bundle.Image == "fake1"
				},
			),
			want: []declcfg.Bundle{
				{Name: "extension1.v1", Package: "extension1", Image: "fake1"},
			},
		},
		{
			name: "filter with Or predicate",
			predicate: filter.Or(
				func(bundle declcfg.Bundle) bool {
					return bundle.Name == "extension1.v1"
				},
				func(bundle declcfg.Bundle) bool {
					return bundle.Image == "fake1"
				},
			),
			want: []declcfg.Bundle{
				{Name: "extension1.v1", Package: "extension1", Image: "fake1"},
				{Name: "extension2.v1", Package: "extension2", Image: "fake1"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			in := []declcfg.Bundle{
				{Name: "extension1.v1", Package: "extension1", Image: "fake1"},
				{Name: "extension1.v2", Package: "extension1", Image: "fake2"},
				{Name: "extension2.v1", Package: "extension2", Image: "fake1"},
			}

			actual := filter.Filter(in, tt.predicate)
			assert.Equal(t, tt.want, actual)
		})
	}
}
