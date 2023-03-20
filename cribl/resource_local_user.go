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
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"first": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
		Username: d.Get("username").(string),
		First:    d.Get("first").(string),
		Last:     d.Get("last").(string),
		Email:    d.Get("email").(string),
		Roles:    d.Get("roles").([]string),
		Id:       d.Get("id").(string),
		Disabled: d.Get("disabled").(bool),
		Password: d.Get("password").(string),
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
		d.Set("username", key.Username)
		d.Set("first", key.First)
		d.Set("last", key.Last)
		d.Set("email", key.Email)
		d.Set("roles", key.Roles)
		d.Set("id", key.Id)
		d.Set("disabled", key.Disabled)
	}
	return diags
}

func resourceLocalUserupdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)

	user := &criblclient.PatchUser{
		Username: d.Get("username").(string),
		First:    d.Get("first").(string),
		Last:     d.Get("last").(string),
		Email:    d.Get("email").(string),
		Roles:    d.Get("roles").([]string),
		Id:       d.Get("id").(string),
		Disabled: d.Get("disabled").(bool),
		Password: d.Get("password").(string),
	}
	changedParams := []string{"username", "first", "last", "email", "id", "disabled", "password", "roles"}
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
