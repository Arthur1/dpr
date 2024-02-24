data "github_release" "dpr" {
  repository  = "dpr"
  owner       = "Arthur1"
  retrieve_by = "latest"
}

locals {
  download_url = [for asset in data.github_release.dpr.assets : asset if lookup(asset, "name") == "dpr-cleaner-eventbridge-lambda.zip"][0].browser_download_url
}

resource "terraform_data" "dpr_cleaner_eventbridge_lambda" {
  triggers_replace = timestamp()

  provisioner "local-exec" {
    command = "wget \"$download_url\" -nv -O dpr-cleaner-eventbridge-lambda.zip"
    environment = {
      download_url = local.download_url
    }
  }
}
