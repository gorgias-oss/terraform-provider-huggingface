resource "huggingface_endpoint" "example" {
  name = "product-identification-reran-soy"
  type = "private"

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
    framework  = "pytorch"
    repository = "sentence-transformers/stsb-roberta-base"
    task       = "text-classification"
    revision   = "main"
    env        = {}
    image = {
      vllm = {
        port                 = 8000
        url                  = "vllm/vllm-openai:gptoss"
        tensor_parallel_size = 1
        kv_cache_dtype       = "auto"
      }
    }
  }

  cloud = {
    region = "us-east-1"
    vendor = "aws"
  }
}
