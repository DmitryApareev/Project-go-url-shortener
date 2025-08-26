package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGenerator(t *testing.T) {
	links := []string{
		"https://docs.google.com/spreadsheets/d/16VDUCGNobAeqdh8JvwBLNwseyKfjKwOBze10LKqyZOo/edit?gid=0#gid=0",
		"https://yandex.ru/maps/213/moscow/?ll=37.599807%2C55.745088&mode=poi&poi%5Bpoint%5D=37.620027%2C55.741556&poi%5Buri%5D=ymapsbm1%3A%2F%2Forg%3Foid%3D21117108341&z=12",
		"https://vkvideo.ru/video-218359460_456241865?t=27m22s",
	}

	for _, link := range links {
		short := GenerateShortLink(link, UserId)

		// Проверяем, что короткая ссылка не пустая и имеет длину 8
		assert.NotEmpty(t, short, "short link should not be empty")
		assert.Len(t, short, 8, "short link should have length 8")
	}
}
