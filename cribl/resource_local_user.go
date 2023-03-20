package cribl

import (
	"context"
	"time"

	"github.com/Vinayaks439/terraform-provider-cribl/criblclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLocalUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"Username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"First": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"Last": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"Email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"Roles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
			},
			"Id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"Disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Default:  true,
			},
			"Password": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
		},
		CreateContext: resourceLocalUserCreate,
		ReadContext:   resourceLocalUserRead,
		UpdateContext: resourceLocalUserupdate,
		DeleteContext: resourceLocalUserDelete,
	}
}

func resourceLocalUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)

	var diags diag.Diagnostics

	user := &criblclient.CreateUser{
		Username: d.Get("Username").(string),
		First:    d.Get("First").(string),
		Last:     d.Get("Last").(string),
		Email:    d.Get("Email").(string),
		Roles:    d.Get("Roles").([]string),
		Id:       d.Get("Id").(string),
		Disabled: d.Get("Disabled").(bool),
		Password: d.Get("Password").(string),
	}

	createUser, err := c.CreateUser(user)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, key := range createUser.Items {
		d.SetId(key.Id)
	}
	resourceLocalUserRead(ctx, d, m)
	return diags
}

func resourceLocalUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)

	var diags diag.Diagnostics
	getUserbyID, err := c.GetUserByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	for _, key := range getUserbyID.Items {
		d.Set("Username", key.Username)
		d.Set("First", key.First)
		d.Set("Last", key.Last)
		d.Set("Email", key.Email)
		d.Set("Roles", key.Roles)
		d.Set("Id", key.Id)
		d.Set("Disabled", key.Disabled)
	}
	return diags
}

func resourceLocalUserupdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)

	user := &criblclient.PatchUser{
		Username: d.Get("Username").(string),
		First:    d.Get("First").(string),
		Last:     d.Get("Last").(string),
		Email:    d.Get("Email").(string),
		Roles:    d.Get("Roles").([]string),
		Id:       d.Get("Id").(string),
		Disabled: d.Get("Disabled").(bool),
		Password: d.Get("Password").(string),
	}
	changedParams := []string{"Username", "First", "Last", "Email", "Id", "Disabled", "Password", "Roles"}
	for _, key := range changedParams {
		if d.HasChange(key) {
			_, err := c.PatchUserInfo(d.Id(), user)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourceLocalUserRead(ctx, d, m)
}

func resourceLocalUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)
	var diags diag.Diagnostics
	_, err := c.DeleteUserbyID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
