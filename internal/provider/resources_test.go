package provider

import (
	"reflect"
	"testing"

	huggingface "github.com/gorgias-oss/huggingface-endpoints-client-go"
)

func TestVllmServerArgsRoundTrip(t *testing.T) {
	serverArgs := []string{
		"--enable-auto-tool-choice",
		"--tool-call-parser",
		"hermes",
	}
	endpoint := huggingface.EndpointDetails{
		Model: huggingface.Model{
			Image: huggingface.Image{
				Vllm: &huggingface.Vllm{ServerArgs: serverArgs},
			},
		},
	}

	providerEndpoint := clientEndpointToProviderEndpoint(endpoint)
	if !reflect.DeepEqual(providerEndpoint.Model.Image.Vllm.ServerArgs, serverArgs) {
		t.Fatalf("server args read as %v, want %v", providerEndpoint.Model.Image.Vllm.ServerArgs, serverArgs)
	}

	createRequest := providerEndpointToCreateEndpointRequest(providerEndpoint)
	if !reflect.DeepEqual(createRequest.Model.Image.Vllm.ServerArgs, serverArgs) {
		t.Fatalf("create request server args are %v, want %v", createRequest.Model.Image.Vllm.ServerArgs, serverArgs)
	}

	updateRequest := providerEndpointToUpdateEndpointRequest(providerEndpoint)
	if !reflect.DeepEqual(updateRequest.Model.Image.Vllm.ServerArgs, serverArgs) {
		t.Fatalf("update request server args are %v, want %v", updateRequest.Model.Image.Vllm.ServerArgs, serverArgs)
	}
}
