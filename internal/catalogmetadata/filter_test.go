package catalogmetadata_test

import (
	"testing"

	"github.com/operator-framework/operator-registry/alpha/declcfg"
	"github.com/stretchr/testify/assert"

	"github.com/operator-framework/operator-controller/internal/catalogmetadata"
)

func TestFilter(t *testing.T) {
	in := []*catalogmetadata.Bundle{
		{Bundle: declcfg.Bundle{Name: "operator1.v1", Package: "operator1", Image: "fake1"}},
		{Bundle: declcfg.Bundle{Name: "operator1.v2", Package: "operator1", Image: "fake2"}},
		{Bundle: declcfg.Bundle{Name: "operator2.v1", Package: "operator2", Image: "fake1"}},
	}

	for _, tt := range []struct {
		name      string
		predicate catalogmetadata.Predicate[catalogmetadata.Bundle]
		want      []*catalogmetadata.Bundle
	}{
		{
			name: "simple filter with one predicate",
			predicate: func(bundle *catalogmetadata.Bundle) bool {
				return bundle.Name == "operator1.v1"
			},
			want: []*catalogmetadata.Bundle{
				{Bundle: declcfg.Bundle{Name: "operator1.v1", Package: "operator1", Image: "fake1"}},
			},
		},
		{
			name: "filter with Not predicate",
			predicate: catalogmetadata.Not(func(bundle *catalogmetadata.Bundle) bool {
				return bundle.Name == "operator1.v1"
			}),
			want: []*catalogmetadata.Bundle{
				{Bundle: declcfg.Bundle{Name: "operator1.v2", Package: "operator1", Image: "fake2"}},
				{Bundle: declcfg.Bundle{Name: "operator2.v1", Package: "operator2", Image: "fake1"}},
			},
		},
		{
			name: "filter with And predicate",
			predicate: catalogmetadata.And(
				func(bundle *catalogmetadata.Bundle) bool {
					return bundle.Name == "operator1.v1"
				},
				func(bundle *catalogmetadata.Bundle) bool {
					return bundle.Image == "fake1"
				},
			),
			want: []*catalogmetadata.Bundle{
				{Bundle: declcfg.Bundle{Name: "operator1.v1", Package: "operator1", Image: "fake1"}},
			},
		},
		{
			name: "filter with Or predicate",
			predicate: catalogmetadata.Or(
				func(bundle *catalogmetadata.Bundle) bool {
					return bundle.Name == "operator1.v1"
				},
				func(bundle *catalogmetadata.Bundle) bool {
					return bundle.Image == "fake1"
				},
			),
			want: []*catalogmetadata.Bundle{
				{Bundle: declcfg.Bundle{Name: "operator1.v1", Package: "operator1", Image: "fake1"}},
				{Bundle: declcfg.Bundle{Name: "operator2.v1", Package: "operator2", Image: "fake1"}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := catalogmetadata.Filter(in, tt.predicate)
			assert.Equal(t, tt.want, actual)
		})
	}
}
