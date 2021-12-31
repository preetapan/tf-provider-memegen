#tf-provider-memegen
My experiment with writing a custom Terraform provider. 
Used the documented [Img FLip API](https://imgflip.com/api)
Learned that APIs that don't have CRUD HTTP endpoints are quite the challenge given how Terraform 
expects providers to work. Still, this was a fun side project

# Build and Install 
- Prerequisites - Terraform 0.12 or higher, Go 1.17
- Run `go mod init` and `go mod tidy` from project root
- `go build -o terraform-provider-meme` && `cp terraform-provider-meme ~/.terraform.d/plugins/github.com/preetapan/meme/1.0.0/<target_arch>`

# Steps 

- Sign up for an account at imgflip.com 
- Set your user name to the env var `MEMEGEN_USERNAME` and password to `MEMEGEN_USERNAME=`
- `terraform apply` . Example output with the generated meme's URL : 
```
Outputs:
meme_url = "https://imgflip.com/i/5zc8z1"
```
- The API takes a template_id and meme text. Find other meme ids at [popular meme IDs](https://imgflip.com/popular_meme_ids)

# Known issues 
This API doesn't have read or delete implemented, so terraform destroy doesn't actually do anything