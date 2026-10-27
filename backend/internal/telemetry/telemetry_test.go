package telemetry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk"
)

func TestNewResourceOnCloudRun(t *testing.T) {
	metadata := map[string]string{
		"/computeMetadata/v1/project/project-id": "test-project",
		"/computeMetadata/v1/instance/id":        "123456789",
		"/computeMetadata/v1/instance/region":    "projects/123/regions/europe-west4",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := metadata[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Metadata-Flavor", "Google")
		_, _ = fmt.Fprint(w, value)
	}))
	t.Cleanup(server.Close)

	// Exercise the real detector: outside GCP it returns no resource, hiding
	// schema conflicts with the SDK's built-in telemetry detector.
	t.Setenv("GCE_METADATA_HOST", strings.TrimPrefix(server.URL, "http://"))
	t.Setenv("K_CONFIGURATION", "cloud-run-service")
	t.Setenv("K_SERVICE", "cloud-run-service")
	t.Setenv("K_REVISION", "cloud-run-service-00001")

	res, err := newResource(t.Context(), "custom-service-name")
	if err != nil {
		t.Fatalf("creating Cloud Run resource: %v", err)
	}
	if res.SchemaURL() == "" {
		t.Fatal("expected a shared schema URL for GCP and SDK attributes")
	}

	for key, want := range map[string]string{
		"service.name":           "custom-service-name",
		"cloud.provider":         "gcp",
		"cloud.platform":         "gcp_cloud_run",
		"cloud.account.id":       "test-project",
		"cloud.region":           "europe-west4",
		"faas.name":              "cloud-run-service",
		"faas.version":           "cloud-run-service-00001",
		"faas.instance":          "123456789",
		"telemetry.sdk.name":     "opentelemetry",
		"telemetry.sdk.language": "go",
		"telemetry.sdk.version":  sdk.Version(),
	} {
		got, ok := res.Set().Value(attribute.Key(key))
		if !ok || got.AsString() != want {
			t.Errorf("resource attribute %s = %q, want %q", key, got.AsString(), want)
		}
	}
}
