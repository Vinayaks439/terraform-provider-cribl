package cribl

import (
	"context"

	"github.com/Vinayaks439/terraform-provider-cribl/criblclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CRIBL_HOST", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CRIBL_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CRIBL_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cribl_local_user": ResourceLocalUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"hashicups_coffees":     dataSourceCoffees(),
			//"hashicups_order":       dataSourceOrder(),
			//"hashicups_ingredients": dataSourceIngredients(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	authDetails := &criblclient.Auth{
		Username: username,
		Password: password,
	}

	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	var diags diag.Diagnostics

	if (username != "") && (password != "") {

		c, err := criblclient.AuthLogin(authDetails, *host)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Cribl client",
				Detail:   "Unable to authenticate user for authenticated Cribl client",
			})

			return nil, diags
		}

		return c, diags
	}

	c, err := criblclient.AuthLogin(nil, *host)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Cribl client",
			Detail:   "Unable to create anonymous Cribl client",
		})
		return nil, diags
	}

	return c, diags
}
