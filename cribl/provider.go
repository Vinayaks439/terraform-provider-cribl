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
				Required:    true,
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
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("API_TOKEN", nil),
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
	api_token := d.Get("api_token").(string)

	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if api_token != "" {
		return api_token, diags
	}

	if (username != "") && (password != "") {
		c, err := criblclient.NewClient(host, &username, &password)
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

	c, err := criblclient.NewClient(host, nil, nil)
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
