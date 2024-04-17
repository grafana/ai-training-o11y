package testutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/dskit/user"
)

func NewHTTPClient(tenant string) *http.Client {
	return &http.Client{
		Transport: NewTenantRoundTripper(http.DefaultTransport, tenant),
		Timeout:   time.Second * 5,
	}
}

type tenantRoundTripper struct {
	next   http.RoundTripper
	tenant string
}

// NewTenantRoundTripper creates a RoundTripper that adds the tenant to each request.
func NewTenantRoundTripper(base http.RoundTripper, tenant string) http.RoundTripper {
	return &tenantRoundTripper{
		next:   base,
		tenant: tenant,
	}
}

func (rt *tenantRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := user.InjectOrgID(req.Context(), rt.tenant)
	req = req.WithContext(ctx)
	err := user.InjectOrgIDIntoHTTPRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not inject org id: %w", err)
	}
	return rt.next.RoundTrip(req)
}
