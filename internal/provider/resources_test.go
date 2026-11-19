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

func TestPreserveWriteOnlyVllmServerArgs(t *testing.T) {
	configuredArgs := []string{
		"--enable-auto-tool-choice",
		"--tool-call-parser",
		"hermes",
	}

	tests := map[string]struct {
		responseArgs []string
		expectedArgs []string
	}{
		"response omits server args": {
			expectedArgs: configuredArgs,
		},
		"response returns server args": {
			responseArgs: []string{"--new-arg"},
			expectedArgs: []string{"--new-arg"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			state := endpointResourceModel{
				Model: Model{Image: Image{Vllm: &Vllm{ServerArgs: test.responseArgs}}},
			}
			previous := endpointResourceModel{
				Model: Model{Image: Image{Vllm: &Vllm{ServerArgs: configuredArgs}}},
			}

			actual := preserveWriteOnlyVllmServerArgs(state, previous)
			if !reflect.DeepEqual(actual.Model.Image.Vllm.ServerArgs, test.expectedArgs) {
				t.Fatalf("server args are %#v, want %#v", actual.Model.Image.Vllm.ServerArgs, test.expectedArgs)
			}
		})
	}
}

func TestPreserveWriteOnlyVllmSettings(t *testing.T) {
	str := func(v string) *string { return &v }
	num := func(v int) *int { return &v }

	previous := endpointResourceModel{
		Model: Model{Image: Image{Vllm: &Vllm{
			KvCacheDtype:        str("fp8"),
			TensorParallelSize:  num(1),
			MaxModelLen:         num(4096),
			MaxNumSeqs:          num(8),
			MaxNumBatchedTokens: num(2048),
		}}},
	}

	t.Run("response omits settings", func(t *testing.T) {
		state := endpointResourceModel{Model: Model{Image: Image{Vllm: &Vllm{}}}}

		actual := preserveWriteOnlyVllmServerArgs(state, previous).Model.Image.Vllm
		if !reflect.DeepEqual(actual, previous.Model.Image.Vllm) {
			t.Fatalf("vllm settings are %#v, want %#v", actual, previous.Model.Image.Vllm)
		}
	})

	t.Run("response returns settings", func(t *testing.T) {
		state := endpointResourceModel{
			Model: Model{Image: Image{Vllm: &Vllm{
				KvCacheDtype:        str("auto"),
				TensorParallelSize:  num(2),
				MaxModelLen:         num(8192),
				MaxNumSeqs:          num(16),
				MaxNumBatchedTokens: num(4096),
			}}},
		}
		want := *state.Model.Image.Vllm

		actual := preserveWriteOnlyVllmServerArgs(state, previous).Model.Image.Vllm
		if !reflect.DeepEqual(*actual, want) {
			t.Fatalf("vllm settings are %#v, want %#v", *actual, want)
		}
	})
}
