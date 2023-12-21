package ads

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_Ads(t *testing.T) {
	expected := Ad{
		ID:        0,
		Title:     "title",
		Text:      "text",
		AuthorID:  123,
		Published: false,
	}

	t.Run("new_ad", func(t *testing.T) {
		ad, _ := New("title", "text", 123)

		assert.Equal(t, ad.ID, expected.ID)
		assert.Equal(t, ad.Title, expected.Title)
		assert.Equal(t, ad.Text, expected.Text)
		assert.Equal(t, ad.AuthorID, expected.AuthorID)
		assert.Equal(t, ad.Published, expected.Published)
	})
}
