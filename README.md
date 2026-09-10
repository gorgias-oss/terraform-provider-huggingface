# Terraform Provider for Hugging Face

A Terraform provider for managing Hugging Face Inference Endpoints, enabling Infrastructure as Code (IaC) for ML model deployments.

## Features

- **Full Lifecycle Management**: Create, read, update, and delete Hugging Face inference endpoints
- **Multiple Inference Frameworks**: Support for TGI, TEI, vLLM, Llama.cpp, and custom Docker images
- **Auto-scaling Configuration**: Configure min/max replicas, scale-to-zero, and scaling metrics
- **Cloud Provider Support**: Deploy to AWS and other cloud vendors
- **Environment Variables**: Configure model-specific environment variables
- **Private and Public Endpoints**: Support for both endpoint types

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (for development)
- A Hugging Face account with API access
- A Hugging Face API token

## Installation

### Using Terraform Registry

```hcl
terraform {
  required_providers {
    huggingface = {
      source  = "gorgias-oss/huggingface"
      version = "~> 0.1"
    }
  }
}
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/gorgias-oss/terraform-provider-huggingface.git
cd terraform-provider-huggingface
```

2. Build the provider:
```bash
go build -o terraform-provider-huggingface
```

3. Install the provider locally:
```bash
mkdir -p ~/.terraform.d/plugins/github.com/gorgias-oss/huggingface/0.1/darwin_amd64
mv terraform-provider-huggingface ~/.terraform.d/plugins/github.com/gorgias-oss/huggingface/0.1/darwin_amd64/
```

## Usage

### Provider Configuration

```hcl
provider "huggingface" {
  host      = "https://api.endpoints.huggingface.cloud/v2/endpoint"
  namespace = "your-namespace"  # Your Hugging Face organization or username
  token     = var.huggingface_token  # Your Hugging Face API token
}
```

### Creating an Inference Endpoint

#### Basic Example

```hcl
resource "huggingface_endpoint" "example" {
  name = "my-inference-endpoint"
  type = "private"
  
  compute = {
    accelerator   = "gpu"
    instance_size = "x1"
    instance_type = "nvidia-l4"
  }
  
  model = {
    framework  = "pytorch"
    repository = "bert-base-uncased"
    task       = "text-classification"
    revision   = "main"
  }
  
  cloud = {
    vendor = "aws"
    region = "us-east-1"
  }
}
```

#### Advanced Example with Auto-scaling

```hcl
resource "huggingface_endpoint" "advanced" {
  name = "scalable-llm-endpoint"
  type = "private"
  
  compute = {
    accelerator   = "gpu"
    instance_size = "x4"
    instance_type = "nvidia-a100"
    scaling = {
      min_replica           = 0
      max_replica           = 4
      scale_to_zero_timeout = 300  # 5 minutes
      measure = {
        hardware_usage = 80.0      # Scale when hardware usage > 80%
        pending_requests = 10.0    # Or when pending requests > 10
      }
    }
  }
  
  model = {
    framework  = "pytorch"
    repository = "meta-llama/Llama-2-7b-chat-hf"
    task       = "text-generation"
    revision   = "main"
    
    env = {
      MAX_BATCH_SIZE = "8"
      MAX_SEQ_LENGTH = "2048"
    }
    
    image = {
      tgi = {
        url                     = "ghcr.io/huggingface/text-generation-inference:latest"
        port                    = 8080
        health_route            = "/health"
        env = {
          MAX_INPUT_TOKENS      = "2048"
          MAX_TOTAL_TOKENS      = "4096"
        }
      }
    }
  }
  
  cloud = {
    vendor = "aws"
    region = "us-west-2"
  }
}
```

#### Using vLLM

```hcl
resource "huggingface_endpoint" "vllm_example" {
  name = "vllm-endpoint"
  type = "private"
  
  compute = {
    accelerator   = "gpu"
    instance_size = "x2"
    instance_type = "nvidia-l4"
  }
  
  model = {
    framework  = "pytorch"
    repository = "mistralai/Mistral-7B-v0.1"
    task       = "text-generation"
    revision   = "main"
    
    image = {
      vllm = {
        url                  = "vllm/vllm-openai:latest"
        port                 = 8000
        tensor_parallel_size = 2
        env = {
          VLLM_API_KEY = var.vllm_api_key
        }
      }
    }
  }
  
  cloud = {
    vendor = "aws"
    region = "eu-west-1"
  }
}
```

## Resource Reference

### `huggingface_endpoint`

#### Arguments

- `name` - (Required) The name of the endpoint
- `type` - (Required) The endpoint type ("private" or "public")
- `compute` - (Required) Compute configuration block
  - `accelerator` - (Required) Hardware accelerator type ("cpu" or "gpu")
  - `instance_size` - (Required) Instance size (e.g., "x1", "x2", "x4")
  - `instance_type` - (Required) Instance type (e.g., "nvidia-l4", "nvidia-a100")
  - `scaling` - (Optional) Auto-scaling configuration
    - `min_replica` - (Optional) Minimum number of replicas
    - `max_replica` - (Optional) Maximum number of replicas
    - `scale_to_zero_timeout` - (Optional) Seconds before scaling to zero
    - `measure` - (Optional) Scaling metrics
      - `hardware_usage` - (Optional) Hardware usage percentage threshold
      - `pending_requests` - (Optional) Pending requests threshold
- `model` - (Required) Model configuration block
  - `framework` - (Required) ML framework ("pytorch", "tensorflow", etc.)
  - `repository` - (Required) Hugging Face model repository
  - `task` - (Required) Task type (e.g., "text-generation", "text-classification")
  - `revision` - (Optional) Model revision/branch
  - `env` - (Optional) Environment variables (map)
  - `image` - (Optional) Custom image configuration
    - `huggingface` - Hugging Face native image config
    - `tgi` - Text Generation Inference config
    - `tgi_neuron` - TGI Neuron config
    - `tei` - Text Embeddings Inference config
    - `vllm` - vLLM config
    - `llamacpp` - Llama.cpp config
    - `custom` - Custom Docker image config
- `cloud` - (Required) Cloud deployment configuration
  - `vendor` - (Required) Cloud vendor ("aws", etc.)
  - `region` - (Required) Deployment region

#### Attributes

- `id` - The endpoint ID
- `status` - Current endpoint status
- `url` - The inference endpoint URL

## Development

### Building from Source

```bash
go build -o terraform-provider-huggingface
```

### Running Tests

```bash
go test ./...
```

### Generating Documentation

```bash
go generate ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- Uses the [Hugging Face Endpoints Client for Go](https://github.com/gorgias-oss/huggingface-endpoints-client-go)
