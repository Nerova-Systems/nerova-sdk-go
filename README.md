# `nerova-sdk-go`

Server-side Go client for Nerova's dedicated stable `tenant-v1` API. The typed
client is built from the published `tenant-v1` OpenAPI contract, so every path,
parameter, and model exists because the contract says so.

This module is server-side only: no browser credential flow, no UI, no
persistent credential storage. Never ship an API key to a browser or a mobile
app.

```bash
go get github.com/nerova-systems/nerova-sdk-go@v0.3.0-preview.2
```

The module path is `github.com/nerova-systems/nerova-sdk-go`, the client lives in the
`generated` package (`github.com/nerova-systems/nerova-sdk-go/generated`), and it needs Go 1.25
or later. `go get` also resolves the module's HTTP and serialization runtime libraries
(the net/http request adapter and the JSON, text, form, and multipart serializers).

## Versioning

The SDK shares one version line with `@nerova/sdk` and `Nerova.Sdk`
(currently `0.3.0-preview.2`). Go tags are semver, so the line is used unchanged:
`0.3.0-preview.2` is the tag `v0.3.0-preview.2`, and `0.3.0-preview.3` will be
`v0.3.0-preview.3`. A stable release is `v0.3.0`. Go treats pre-releases as unstable and
only picks one for `@latest` while no stable release exists. `v0.3.0-preview.1` is retracted
in favor of `v0.3.0-preview.2`. The API surface may change before 1.0, so pin the exact
version.

## Quickstart

The API has one public host, `https://api.nerovasystems.com`, which is the
client's default base URL. Every API key is a Live key (`nrv_live_`) that reaches
Production; approved platforms start on a free testing allowance. Read the key from the
environment and send it as `Authorization: Bearer <key>`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/microsoft/kiota-abstractions-go/authentication"
	httpadapter "github.com/microsoft/kiota-http-go"

	"github.com/nerova-systems/nerova-sdk-go/generated"
)

func main() {
	apiKey := os.Getenv("NEROVA_API_KEY")
	tenantID := os.Getenv("TENANT_ID")
	if apiKey == "" || tenantID == "" {
		log.Fatal("NEROVA_API_KEY and TENANT_ID are required")
	}

	authenticationProvider, err := authentication.NewApiKeyAuthenticationProvider(
		"Bearer "+apiKey, "Authorization", authentication.HEADER_KEYLOCATION,
	)
	if err != nil {
		log.Fatal(err)
	}
	requestAdapter, err := httpadapter.NewNetHttpRequestAdapter(authenticationProvider)
	if err != nil {
		log.Fatal(err)
	}
	client := generated.NewNerovaPartnerClient(requestAdapter)

	manifest, err := client.Api().V1().Tenants().ByTenantId(tenantID).Activation().Manifest().Get(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(manifest.GetState(), manifest.GetNextAllowedActions())
}
```

The client mirrors the URL structure of the API (`client.Api().V1().Status()`,
`client.Api().V1().Tenants().ByTenantId(id)`, and so on); request and response types live in
`github.com/nerova-systems/nerova-sdk-go/generated/models`.

## Mutations and errors

Every mutation requires a caller-owned `Idempotency-Key` header supplied through
the request configuration; the client performs no automatic retries and never
generates a hidden key:

```go
headers := abstractions.NewRequestHeaders()
headers.Add("Idempotency-Key", "provision-tenant-your-crm-id-123")
configuration := &abstractions.RequestConfiguration[abstractions.DefaultQueryParameters]{Headers: headers}

displayName, externalReference := "Demo Salon", "your-crm-id-123"
body := models.NewCreateTenantV1Request()
body.SetDisplayName(&displayName)
body.SetExternalReference(&externalReference)

created, err := client.Api().V1().Tenants().Post(ctx, body, configuration)
```

(`abstractions` is `github.com/microsoft/kiota-abstractions-go`.) Errors follow RFC 9457: a failed
call returns the deserialized `*models.ProblemDetails`, which also carries `ResponseStatusCode`.
The stable `code` and the `correlationId` arrive in `GetAdditionalData()`; unwrap it with
`errors.As(err, &problem)`. Enums are open string sets, so treat unknown values as forward
compatibility rather than errors. List endpoints paginate with `unixms|id` cursors: pass `Limit` and
`Cursor` in the query parameters and follow `GetNextCursor()` until the server stops returning one.

## Documentation and support

Guides, the API reference, and support contacts are on
[docs.nerovasystems.com](https://docs.nerovasystems.com). Start with the
[quickstart](https://docs.nerovasystems.com/getting-started/quickstart) and the
[SDK overview](https://docs.nerovasystems.com/sdks).

This repository is a read-only distribution of the module so that the Go module proxy can
serve it. Do not open issues or pull requests here, they are not monitored; contact Nerova
support through the documentation site instead.
