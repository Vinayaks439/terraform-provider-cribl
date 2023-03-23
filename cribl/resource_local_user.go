package cribl

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Vinayaks439/terraform-provider-cribl/criblclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLocalUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"unique_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"items": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: `The username of the User.`,
						},
						"first": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: `The first name of the User.`,
						},
						"last": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: `The last name of the User.`,
						},
						"email": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: `The email of the User.`,
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Computed:    false,
							Description: `The id of the User.`,
						},
						"roles": &schema.Schema{
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: `The list of roles attached to the User.`,
						},
						"password": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `The password attached to the User.`,
						},
						"disabled": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Make it true to disable the user.`,
						},
					},
				},
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
	items := d.Get("items").([]interface{})
	users := []criblclient.CreateUser{}
	var diags diag.Diagnostics
	for _, item := range items {
		i := item.(map[string]interface{})
		r := i["roles"].(*schema.Set).List()
		s := make([]string, len(r))
		for i, role := range r {
			s[i] = role.(string)
		}
		d := i["disabled"].(string)
		disable, err := strconv.ParseBool(d)
		if err != nil {
			log.Fatalln("Error Converting string to bool for disbaled")
		}
		user := criblclient.CreateUser{
			Username: i["username"].(string),
			First:    i["first"].(string),
			Last:     i["last"].(string),
			Email:    i["email"].(string),
			Id:       i["id"].(string),
			Roles:    s,
			Password: i["password"].(string),
			Disabled: disable,
		}

		users = append(users, user)
	}

	// user := &criblclient.CreateUser{
	// 	Username: d.Get("username").(string),
	// 	First:    d.Get("first").(string),
	// 	Last:     d.Get("last").(string),
	// 	Email:    d.Get("email").(string),
	// 	Roles:    roles,
	// 	Id:       d.Get("id").(string),
	// 	Disabled: d.Get("disabled").(bool),
	// 	Password: d.Get("password").(string),
	// }

	_, err := c.CreateUser(users)
	if err != nil {
		log.Fatalln("Error Creating the user : ", err)
		return diag.FromErr(err)
	}
	// for _, key := range createUser.Items {
	// 	d.SetId(key.Id)
	// }
	d.SetId(d.Get("unique_id").(string))
	//resourceLocalUserRead(ctx, d, m)
	return diags
}

func resourceLocalUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)

	var diags diag.Diagnostics
	getUserbyID, err := c.GetUserByID(d.Get("id").(string))
	if err != nil {
		log.Fatalln("Error Fetching the user by iD : "+d.Get("id").(string)+"Error is : ", err)
		return diag.FromErr(err)
	}
	// for _, key := range getUserbyID.Items {
	// 	d.Set("username", key.Username)
	// 	d.Set("first", key.First)
	// 	d.Set("last", key.Last)
	// 	d.Set("email", key.Email)
	// 	if key.Roles != nil && len(key.Roles) > 0 {
	// 		roles := make([]interface{}, len(key.Roles))
	// 		for i, v := range key.Roles {
	// 			roles[i] = v
	// 		}
	// 		if err := d.Set("roles", schema.NewSet(schema.HashString, roles)); err != nil {
	// 			diag.FromErr(err)
	// 		}
	// 	}
	// 	d.Set("id", key.Id)
	// 	d.Set("disabled", key.Disabled)
	// }
	userItems := flattencreateUsers(&getUserbyID.Items)
	if err := d.Set("items", userItems); err != nil {
		log.Fatalln("Error flattening the users : ", err)
		return diag.FromErr(err)
	}

	return diags
}

func resourceLocalUserupdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*criblclient.Client)
	var roles []string
	items := d.Get("items").([]interface{})
	users := []criblclient.PatchUser{}
	for _, item := range items {
		i := item.(map[string]interface{})
		for _, role := range i["roles"].([]interface{}) {
			roles = role.([]string)
		}
		user := criblclient.PatchUser{
			Username: i["username"].(string),
			First:    i["first"].(string),
			Last:     i["last"].(string),
			Email:    i["email"].(string),
			Id:       i["id"].(string),
			Roles:    roles,
			Password: i["password"].(string),
			Disabled: i["disabled"].(bool),
		}

		users = append(users, user)
	}
	changedParams := []string{"username", "first", "last", "email", "id", "disabled", "password", "roles"}
	for _, key := range changedParams {
		if d.HasChange(key) {
			_, err := c.PatchUserInfo(d.Id(), users)
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
	_, err := c.DeleteUserbyID(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func flattencreateUsers(listUser *[]criblclient.GetUser) []interface{} {
	if listUser != nil {
		ois := make([]interface{}, len(*listUser), len(*listUser))

		for i, listuser := range *listUser {
			oi := make(map[string]interface{})
			oi["username"] = listuser.Username
			oi["first"] = listuser.First
			oi["last"] = listuser.Last
			oi["email"] = listuser.Email
			oi["disabled"] = listuser.Disabled
			oi["roles"] = listuser.Roles
			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
