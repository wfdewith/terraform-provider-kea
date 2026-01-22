package errors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewProviderDataTypeError(value any) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Unexpected ProviderData Type",
		fmt.Sprintf(
			"Expected clients.KeaClients, got: %T. "+
				"Please report this issue to the provider developers.",
			value,
		),
	)
}

func NewUnconfiguredClientError(client string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Unconfigured Client Error",
		fmt.Sprintf(
			"Expected Kea %s client, but it was not configured. "+
				"Please configure the '%s' block in the provider configuration.",
			strings.ToUpper(client),
			client,
		),
	)
}
