package render_test

import (
	"os"
	"testing"

	"github.com/vandordev/vxt"
	"github.com/vandordev/vxt/source"
)

func TestRenderSingleFileOriskinUsecaseTemplate(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-usecase.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/usecase.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"use_case_name": "CreateBooking",
		"receiver_name": "createBooking",
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if out == "" {
		t.Fatal("expected rendered usecase template")
	}
	if want := "type CreateBookingInput struct {"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := "type createBookingUseCase struct {"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func TestRenderSingleFileOriskinContextModuleTemplateWithConstructors(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-context-module.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/context_module.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"has_constructors": true,
		"context_package":  "booking",
		"import_path":      "github.com/oriskin/home-service/internal/core/contexts/booking/application/usecase",
		"constructors": []any{
			"NewCreateBookingUseCase",
			"NewCompleteBookingUseCase",
		},
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := "usecase.NewCreateBookingUseCase,"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := "usecase.NewCompleteBookingUseCase,"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func TestRenderSingleFileOriskinContextModuleTemplateWithoutConstructors(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-context-module-empty.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/context_module.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"has_constructors": false,
		"context_package":  "clinic",
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := "package clinic"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if contains(out, "non-empty") {
		t.Fatal("unexpected non-empty branch in output")
	}
}

func TestRenderSingleFileOriskinGatewayRegistryTemplateWithEntries(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-gateway-registry-module.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/gateway_registry_module.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"has_entries": true,
		"entries": []any{
			map[string]any{
				"package_name": "createbooking",
				"import_path":  "github.com/oriskin/home-service/internal/gateway/createbooking",
				"constructor":  "NewCreateBooking",
			},
			map[string]any{
				"package_name": "completebooking",
				"import_path":  "github.com/oriskin/home-service/internal/gateway/completebooking",
				"constructor":  "NewCompleteBooking",
			},
		},
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := `createbooking "github.com/oriskin/home-service/internal/gateway/createbooking"`; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := "completebooking.NewCompleteBooking,"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func TestRenderSingleFileOriskinQueryServiceTemplate(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-query-service.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/query_service.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"package_name":  "customer",
		"query_name":    "FindCustomer",
		"receiver_name": "findCustomer",
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := "type FindCustomerInput struct {"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := "type findCustomerService struct {"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func TestRenderSingleFileOriskinRepositoryModuleTemplateWithEntries(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-repository-module.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/repository_module.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"has_entries": true,
		"entries": []any{
			map[string]any{
				"package_name": "customerrepo",
				"import_path":  "github.com/oriskin/home-service/internal/infrastructure/repository/customerrepo",
			},
			map[string]any{
				"package_name": "bookingrepo",
				"import_path":  "github.com/oriskin/home-service/internal/infrastructure/repository/bookingrepo",
			},
		},
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := `customerrepo "github.com/oriskin/home-service/internal/infrastructure/repository/customerrepo"`; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := "bookingrepo.Module,"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func TestRenderSingleFileOriskinHTTPAdminHandlerTemplate(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-http-admin-handler.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/http_admin_handler.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"name":        "ListBookings",
		"receiver":    "listBookingsHandler",
		"method":      "GET",
		"path":        "/bookings/{booking_id}",
		"summary":     "List bookings",
		"tag":         "Bookings",
		"paginated":   true,
		"has_body":    false,
		"path_params": []any{map[string]any{"field_name": "BookingID", "raw_name": "booking_id"}},
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := `BookingID string ` + "`path:\"booking_id\"`"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := `Page httpcontract.OptionalParam[int] ` + "`query:\"page\"`"; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := `method.GET(h.api, "/bookings/{booking_id}", method.Operation{`; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func TestRenderSingleFileOriskinJobRegistryTemplateWithEntries(t *testing.T) {
	src := source.Source{
		ID:   "oriskin-job-module.vxt",
		Text: mustReadFixture(t, "testdata/oriskin/job_module.vxt"),
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{
		"has_entries": true,
		"entries": []any{
			map[string]any{
				"package_name": "bookingjobs",
				"import_path":  "github.com/oriskin/home-service/internal/core/jobs/booking",
				"job_key":      "booking.send_reminder",
				"job_name":     "SendReminder",
			},
			map[string]any{
				"package_name": "clinicjobs",
				"import_path":  "github.com/oriskin/home-service/internal/core/jobs/clinic",
				"job_key":      "clinic.sync_schedule",
				"job_name":     "SyncSchedule",
			},
		},
	})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if want := `bookingjobs "github.com/oriskin/home-service/internal/core/jobs/booking"`; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
	if want := `"clinic.sync_schedule": func() contracts.Job { return &clinicjobs.SyncSchedule{} },`; !contains(out, want) {
		t.Fatalf("expected %q in output", want)
	}
}

func contains(s, needle string) bool {
	return len(needle) > 0 && len(s) >= len(needle) && func() bool {
		for i := 0; i+len(needle) <= len(s); i++ {
			if s[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	}()
}

func mustReadFixture(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}

	return string(content)
}
