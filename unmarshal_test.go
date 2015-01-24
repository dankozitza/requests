package requests

import (
	"bytes"
	"testing"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Emote string

func (e *Emote) Receive(stoic interface{}) error {
	*e = Emote(stoic.(string) + "! :D")
	return nil
}

func TestReceiver(t *testing.T) {
	var greeting Emote
	greeting.Receive("hello")

	assert.Equal(t, greeting, "hello! :D")
}

type testTarget struct {
    Foo int64
    Bar float64
    Baz string
	 Qux *Emote
}

func createTestTarget(request *http.Request) (*testTarget, error) {
    target := new(testTarget)
    if err := New(request).Unmarshal(target); err != nil {
        if inputErrs, ok := err.(InputErrors); ok {
            // inputErrs is a map of input names to error
            // messages, so send them to a function to turn
            // them into a proper user-friendly error message.
            return nil, inputErrs
        }
        return nil, err
    }
    return target, nil
}

func TestUnmarshal_UrlEncoded(t *testing.T) {
	body := bytes.NewBufferString(`foo=1&bar=2.7&baz=taz&qux=welcome`)
	httpRequest, err := http.NewRequest("POST", "/", body)
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(t, err)
	example, err := createTestTarget(httpRequest)
	require.NoError(t, err)

	var expectEmote Emote
	expectEmote.Receive("hello")

	assert.Equal(t, int64(1), example.Foo)
	assert.Equal(t, float64(2.7), example.Bar)
	assert.Equal(t, "taz", example.Baz)
	assert.Equal(t, "welcome! :D", string(*example.Qux))
}

func TestUnmarshal_Errors(t *testing.T) {
	body := bytes.NewBufferString(`foo=1&bar=2.7&baz=taz&qux=welcome`)
	httpRequest, _ := http.NewRequest("POST", "/", body)
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var target Unmarshaller
	err := New(httpRequest).Unmarshal(target)
	assert.Error(t, err)
}
