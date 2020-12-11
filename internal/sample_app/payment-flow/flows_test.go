package payment_flow

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvalidPayloadReturnsError(t *testing.T) {

	client := resty.New().SetHostURL(fmt.Sprintf("http://localhost:%d/v2", ServerPort))

	_, err := client.R().
		EnableTrace().
		SetBody("Hello").
		Post("payments")

	require := require.New(t)
	require.Error(err)
}

func TestValidPayment(t *testing.T) {

	client := resty.New().SetHostURL(fmt.Sprintf("http://localhost:%d/v2", ServerPort))

	_, err := client.R().
		EnableTrace().
		SetBody("{}").
		Post("payments")

	require := require.New(t)
	require.NoError(err)
}
