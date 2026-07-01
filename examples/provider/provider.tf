terraform {
  required_providers {
    huggingface = {
      source = "gorgias-oss/huggingface"
    }
  }
}

provider "huggingface" {
  host      = "https://api.endpoints.huggingface.cloud/v2/endpoint"
  namespace = "some-namespace"
  token     = ""
}
