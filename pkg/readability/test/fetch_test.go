package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/readability"
)

func TestFetchNonReadablePage(t *testing.T) {
	ctx := context.TODO()
	article, err := readability.FetchArticle(ctx, "https://about.readflow.app/")
	assert.NotNil(t, err)
	assert.Equal(t, "unable to extract content from HTML page", err.Error())
	assert.NotNil(t, article)
	assert.Equal(t, "readflow", article.Title)
	assert.Equal(t, "Read your Internet article flow in one place with complete peace of mind and freedom", *article.Text)
	assert.Equal(t, "https://about.readflow.app/images/readflow.png", *article.Image)
}
