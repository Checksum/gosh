package handlers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validHtml = `
<!DOCTYPE html>
<html>
    <head>
        <title>Valid HTML</title>
    </head>
    <body>
        <h1>Valid</h1>
    </body>
</html>
`

const invalidHtml = `
<h1>Invalid</h1>
`

func TestInvalidTransformResponse(t *testing.T) {
	content, err := transformResponse(bytes.NewReader([]byte(invalidHtml)))

	if err != nil {
		t.Errorf("%+v", err)
	}

	assert.NotEmpty(t, content)
}

func TestValidTransformResponse(t *testing.T) {
	content, err := transformResponse(bytes.NewReader([]byte(validHtml)))

	if err != nil {
		t.Errorf("%+v", err)
	}

	assert.NotEmpty(t, content)
}
