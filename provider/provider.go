package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"net/http"
	"strings"
)

type MemeResponse struct {
	Success bool
	Data MemeData
}

type MemeData struct {
	Url string `json:"url"`
	PageURL string `json:"page_url"`
}

func Provider() func() *schema.Provider {
	return func() *schema.Provider {

		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"meme_generator": resourceServer(),
			},
			Schema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("MEMEGEN_USERNAME", nil),
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("MEMEGEN_PASSWORD", nil),
				},
			},
			ConfigureContextFunc: providerConfigure,
		}
		return p
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
    creds := make(map[string]string)
	if (username != "") && (password != "") {
		creds["username"] = username
		creds["password"] = password
	}
	return creds, diags
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: memeCreate,
		Read:   memeRead,
		Update: memeUpdate,
		Delete: memeDelete,

		Schema: map[string]*schema.Schema{
			"text": {
				Type:        schema.TypeString,
				Required:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
			},
			"page_url": {
				Type:        schema.TypeString,
				Computed: true,
			},
		},
	}
}

func memeCreate(d *schema.ResourceData, m interface{}) error {
	// create meme from api at https://api.imgflip.com/caption_image
	creds := m.(map[string]string)
	userName := creds["username"]
	passwd := creds["password"]
	templateID := d.Get("template_id").(string)
	text := d.Get("text").(string)
	resp, err := generateMeme(userName, passwd, templateID, text)
	if err == nil {
		d.SetId(resp.Data.Url)
		d.Set("page_url", resp.Data.PageURL)
	}
	return nil
}

func generateMeme(username string, passwd string, templateID string, message string) (MemeResponse, error) {
	var bodyStr bytes.Buffer
	bodyStr.WriteString("username=")
	bodyStr.WriteString(username)
	bodyStr.WriteString("&")
	bodyStr.WriteString("password=")
	bodyStr.WriteString(passwd)
	bodyStr.WriteString("&")
	bodyStr.WriteString("template_id=")
	bodyStr.WriteString(templateID)
	bodyStr.WriteString("&")
	bodyStr.WriteString("text0=")
	bodyStr.WriteString(message)
	bodyStr.WriteString("&")
	body := strings.NewReader(bodyStr.String())
	resp, err := http.Post("https://api.imgflip.com/caption_image", "application/x-www-form-urlencoded", body )
	var memeResponse MemeResponse
	if err != nil {
		fmt.Println("*****error invoking API *****", err)
		return memeResponse, err
	}
	defer resp.Body.Close()
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return memeResponse, err
	}
    err = json.Unmarshal(responseBytes, &memeResponse)
	if err != nil {
		return memeResponse, err
	}
	return memeResponse, nil
}

func memeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func memeUpdate(d *schema.ResourceData, m interface{}) error {
	return memeCreate(d, m)
}

func memeDelete(d *schema.ResourceData, m interface{}) error {
	// no delete API available
	return nil
}
