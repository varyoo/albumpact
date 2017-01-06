package qobuz

import (
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"testing"
)

func TestCercle(t *testing.T) {
	source.SourceTest(t, Get,
		"http://www.qobuz.com/fr-fr/album/cercle-subaquatique-profil-de-face-ep"+
			"-various-artists/3663729010001",
		meta.AlbumTags{
			AlbumTag:       "Cercle subaquatique Profil de Face - EP",
			AlbumArtistTag: "Various Artists",
			DateTag:        "2016",
		}, []meta.TrackTags{
			{TrackNumberTag: 1, ArtistTag: "Février", TitleTag: "Château rouge"},
			{TrackNumberTag: 2, ArtistTag: "Lewis OfMan", TitleTag: "Yo bene"},
			{TrackNumberTag: 3, ArtistTag: "Magnum", TitleTag: "L'épée à la main"},
			{TrackNumberTag: 4, ArtistTag: "Bleu Toucan", TitleTag: "Hanoï Café"},
			{TrackNumberTag: 5, ArtistTag: "Vendredi sur Mer",
				TitleTag: "La femme à la peau bleue (Chez toi)",
			},
		})
}
