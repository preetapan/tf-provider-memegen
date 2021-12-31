resource "meme_generator" "my_meme" {
	template_id = var.template_id
	text = var.text
}

output "meme_url" {
  value = meme_generator.my_meme.page_url
}
