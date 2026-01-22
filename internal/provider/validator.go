package provider

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = keaURLValidator{}

type keaURLValidator struct{}

func (v keaURLValidator) Description(ctx context.Context) string {
	return "string must be a valid URL with http, https, or unix scheme"
}

func (v keaURLValidator) MarkdownDescription(ctx context.Context) string {
	return "string must be a valid URL with `http`, `https`, or `unix` scheme"
}

func (v keaURLValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	u, err := url.Parse(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			"The value must be a valid URL: "+err.Error(),
		)
		return
	}

	switch u.Scheme {
	case "http", "https", "unix":
	default:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL Scheme",
			"The scheme must be one of: http, https, unix",
		)
	}
}

func IsValidKeaURL() validator.String {
	return keaURLValidator{}
}
