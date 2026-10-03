package nerovasdk_test

import (
	"context"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/authentication"
	httpadapter "github.com/microsoft/kiota-http-go"

	"github.com/nerova-systems/nerova-sdk-go/generated"
)

// TestClientBuildsOffline constructs the client on the net/http request adapter and
// inspects a request without touching the network.
func TestClientBuildsOffline(t *testing.T) {
	authenticationProvider, err := authentication.NewApiKeyAuthenticationProvider(
		"Bearer nrv_live_smoke_test_value_is_not_a_key", "Authorization", authentication.HEADER_KEYLOCATION,
	)
	if err != nil {
		t.Fatalf("building the authentication provider failed: %v", err)
	}

	requestAdapter, err := httpadapter.NewNetHttpRequestAdapter(authenticationProvider)
	if err != nil {
		t.Fatalf("building the request adapter failed: %v", err)
	}

	client := generated.NewNerovaPartnerClient(requestAdapter)
	if baseUrl := requestAdapter.GetBaseUrl(); baseUrl != "https://api.nerovasystems.com" {
		t.Fatalf("unexpected default base URL: %s", baseUrl)
	}

	requestInformation, err := client.Api().V1().Status().ToGetRequestInformation(context.Background(), nil)
	if err != nil {
		t.Fatalf("building the request failed: %v", err)
	}
	if requestInformation.UrlTemplate != "{+baseurl}/api/v1/status" {
		t.Fatalf("unexpected URL template: %s", requestInformation.UrlTemplate)
	}
	if requestInformation.PathParameters["baseurl"] != "https://api.nerovasystems.com" {
		t.Fatalf("base URL missing from path parameters: %v", requestInformation.PathParameters)
	}
}
