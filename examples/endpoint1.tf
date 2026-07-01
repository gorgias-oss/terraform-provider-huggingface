resource "huggingface_endpoint" "product_identification_reran_soy" {
  name = "product-identification-reran-soy"

  compute = {
    accelerator   = "gpu"
    instance_size = "x1"
    instance_type = "nvidia-l4"
    scaling = {
      min_replica           = 0
      max_replica           = 1
      scale_to_zero_timeout = 60
      measure = {
        hardware_usage = 80.0
      }
    }
  }

  model = {
    framework = "pytorch"
    image = {
      vllm = {
        port                 = 8000
        url                  = "vllm/vllm-openai:gptoss"
        tensor_parallel_size = 1
        kv_cache_dtype       = "auto"
      }
    }
    env        = {}
    repository = "sentence-transformers/stsb-roberta-base"
    task       = "text-classification"
    revision   = "main"
  }

  cloud = {
    region = "us-east-1"
    vendor = "aws"
  }

  type = "private"
}

output "product_identification_reran_soy" {
  value = huggingface_endpoint.product_identification_reran_soy
}
