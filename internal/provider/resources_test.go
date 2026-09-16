package provider

import (
	"reflect"
	"testing"

	huggingface "github.com/gorgias-oss/huggingface-endpoints-client-go"
)

func TestVllmFieldsRoundTrip(t *testing.T) {
	serverArgs := []string{
		"--enable-auto-tool-choice",
		"--tool-call-parser",
		"hermes",
	}
	maxModelLen := 32768
	gpuMemoryUtilization := 0.95
	enforceEager := true
	blockSize := 16
	swapSpace := 4
	endpoint := huggingface.EndpointDetails{
		Model: huggingface.Model{
			Image: huggingface.Image{
				Vllm: &huggingface.Vllm{
					MaxModelLen:          &maxModelLen,
					GpuMemoryUtilization: &gpuMemoryUtilization,
					EnforceEager:         &enforceEager,
					BlockSize:            &blockSize,
					SwapSpace:            &swapSpace,
					ServerArgs:           serverArgs,
				},
			},
		},
	}

	providerEndpoint := clientEndpointToProviderEndpoint(endpoint)
	if !reflect.DeepEqual(providerEndpoint.Model.Image.Vllm, &Vllm{
		MaxModelLen:          &maxModelLen,
		GpuMemoryUtilization: &gpuMemoryUtilization,
		EnforceEager:         &enforceEager,
		BlockSize:            &blockSize,
		SwapSpace:            &swapSpace,
		ServerArgs:           serverArgs,
	}) {
		t.Fatalf("vLLM fields read as %#v", providerEndpoint.Model.Image.Vllm)
	}

	createRequest := providerEndpointToCreateEndpointRequest(providerEndpoint)
	if !reflect.DeepEqual(createRequest.Model.Image.Vllm, endpoint.Model.Image.Vllm) {
		t.Fatalf("create request vLLM fields are %#v", createRequest.Model.Image.Vllm)
	}

	updateRequest := providerEndpointToUpdateEndpointRequest(providerEndpoint)
	if !reflect.DeepEqual(updateRequest.Model.Image.Vllm, endpoint.Model.Image.Vllm) {
		t.Fatalf("update request vLLM fields are %#v", updateRequest.Model.Image.Vllm)
	}
}

func TestProviderIncludedInUpdateRequest(t *testing.T) {
	endpoint := endpointResourceModel{
		Cloud: Cloud{
			Region: "us-east-2",
			Vendor: "aws",
		},
	}

	updateRequest := providerEndpointToUpdateEndpointRequest(endpoint)
	if !reflect.DeepEqual(updateRequest.Provider, &huggingface.Provider{
		Region: "us-east-2",
		Vendor: "aws",
	}) {
		t.Fatalf("update request provider is %#v", updateRequest.Provider)
	}
}
