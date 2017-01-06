package discogs

import (
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"testing"
)

func TestCocoon(t *testing.T) {
	source.SourceTest(t, Get,
		"https://www.discogs.com/Dj-SneakPhil-WeeksJoss-Moog-Various/release/1997730",
		meta.AlbumTags{
			AlbumTag:       "Various",
			AlbumArtistTag: "DJ Sneak & Joss Moog & Phil Weeks",
			DateTag:        "2009",
		}, []meta.TrackTags{
			{TrackNumberTag: 1, ArtistTag: "DJ Sneak", TitleTag: "Hit Em Where It Hurts"},
			{TrackNumberTag: 2, ArtistTag: "DJ Sneak",
				TitleTag: "Hit Em Where It Hurts (PW & JM Remix)"},
			{TrackNumberTag: 3, ArtistTag: "Joss Moog & Phil Weeks",
				TitleTag: "Back In Effect (DJ Sneak Remix)"},
			{TrackNumberTag: 4, ArtistTag: "DJ Sneak", TitleTag: "Just A Good Party"},
		})
}
func TestLili(t *testing.T) {
	source.SourceTest(t, Get,
		"https://www.discogs.com/Lilimarche-Chansons-Polaroids/release/8305596",
		meta.AlbumTags{
			AlbumTag:       "Chansons Polaroids",
			AlbumArtistTag: "Lilimarche",
			DateTag:        "2016",
		}, []meta.TrackTags{
			{TrackNumberTag: 1, ArtistTag: "Lilimarche", TitleTag: "Camille"},
			{TrackNumberTag: 2, ArtistTag: "Lilimarche", TitleTag: "Amour D'Ete"},
			{TrackNumberTag: 3, ArtistTag: "Lilimarche", TitleTag: "Flashball"},
			{TrackNumberTag: 4, ArtistTag: "Lilimarche", TitleTag: "Francois"},
		})
}
