package action_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ocv1 "github.com/operator-framework/operator-controller/api/v1"
	"github.com/operator-framework/operator-controller/internal/operator-controller/action"
)

func Test_SyntheticUserRestConfigMapper_UsesDefaultConfigMapper(t *testing.T) {
	isDefaultRequestMapperUsed := false
	defaultServiceMapper := func(ctx context.Context, o client.Object, c *rest.Config) (*rest.Config, error) {
		isDefaultRequestMapperUsed = true
		return c, nil
	}
	syntheticAuthServiceMapper := action.SyntheticUserRestConfigMapper(defaultServiceMapper)
	obj := &ocv1.ClusterExtension{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-clusterextension",
		},
		Spec: ocv1.ClusterExtensionSpec{
			ServiceAccount: ocv1.ServiceAccountReference{
				Name: "my-service-account",
			},
			Namespace: "my-namespace",
		},
	}
	actualCfg, err := syntheticAuthServiceMapper(context.Background(), obj, &rest.Config{})
	require.NoError(t, err)
	require.NotNil(t, actualCfg)
	require.True(t, isDefaultRequestMapperUsed)
}

func Test_SyntheticUserRestConfigMapper_UsesSyntheticAuthMapper(t *testing.T) {
	syntheticAuthServiceMapper := action.SyntheticUserRestConfigMapper(func(ctx context.Context, o client.Object, c *rest.Config) (*rest.Config, error) {
		return c, nil
	})
	obj := &ocv1.ClusterExtension{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-clusterextension",
		},
		Spec: ocv1.ClusterExtensionSpec{
			ServiceAccount: ocv1.ServiceAccountReference{
				Name: "olm.synthetic-user",
			},
			Namespace: "my-namespace",
		},
	}
	actualCfg, err := syntheticAuthServiceMapper(context.Background(), obj, &rest.Config{})
	require.NoError(t, err)
	require.NotNil(t, actualCfg)

	// test that the impersonation headers are appropriately injected into the request
	// by wrapping a fake round tripper around the returned configurations transport
	// nolint:bodyclose
	_, _ = actualCfg.WrapTransport(fakeRoundTripper(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "olm:clusterextension:my-clusterextension", req.Header.Get("Impersonate-User"))
		require.Equal(t, "olm:clusterextensions", req.Header.Get("Impersonate-Group"))
		return &http.Response{}, nil
	})).RoundTrip(&http.Request{})
}

type fakeRoundTripper func(req *http.Request) (*http.Response, error)

func (f fakeRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
