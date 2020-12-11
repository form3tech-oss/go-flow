package payment_flow

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFlow(t *testing.T) {

	client := resty.New().SetHostURL(fmt.Sprintf("http://localhost:%d/v2", ServerPort))

	resp, err := client.R().
		EnableTrace().
		SetBody("Hello").
		Post("payments")

	require := require.New(t)
	require.NoError(err)

	assert.NotNil(t, resp)

}
