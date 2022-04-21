// These tests are supposed to run with `proxytest` build tag, and this way we can leverage the testing framework in "proxytest" package.
// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test -tags=proxytest ./...

//go:build proxytest

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestDuplicateBodyContext_OnHttpResponseBody(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	t.Run("can duplicate response body", func(t *testing.T) {
		// Create http context.
		id := host.InitializeHttpContext()

		// Call OnRequestHeaders.
		_ = host.CallOnRequestHeaders(id, [][2]string{
			{"x-duplicate", "true"},
		}, true)

		host.CallOnResponseHeaders(id, [][2]string{}, true)
		host.CallOnResponseBody(id, []byte(`body-data`), true)
		response := host.GetSentLocalResponse(id)
		require.Equal(t, `body-databody-data`, response)
	})

}
