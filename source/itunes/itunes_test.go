package itunes

import (
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"testing"
)

func TestCeline(t *testing.T) {
	source.SourceTest(t, Get,
		"https://itunes.apple.com/fr/album/cover-girl-ep/id989283907",
		meta.AlbumTags{
			AlbumTag:       "Cover Girl - EP",
			AlbumArtistTag: "Céline Tolosa",
			DateTag:        "2015",
		}, []meta.TrackTags{
			{TrackNumberTag: 1, ArtistTag: "Céline Tolosa", TitleTag: "Cover Girl"},
			{TrackNumberTag: 2, ArtistTag: "Céline Tolosa", TitleTag: "Rue Mansart"},
			{TrackNumberTag: 3, ArtistTag: "Céline Tolosa", TitleTag: "Tu es fantastique"},
			{TrackNumberTag: 4, ArtistTag: "Céline Tolosa", TitleTag: "Fais-moi souffrir"},
		})
}
func TestProfilDeFace(t *testing.T) {
	source.SourceTest(t, Get,
		"https://itunes.apple.com/fr/album/cercle-subaquatique-profil/id1122147451?l=en",
		meta.AlbumTags{
			AlbumTag:       "Cercle subaquatique Profil de Face - EP",
			AlbumArtistTag: "Various Artists",
			DateTag:        "2016",
		}, []meta.TrackTags{
			{TrackNumberTag: 1, ArtistTag: "Février", TitleTag: "Château rouge"},
			{TrackNumberTag: 2, ArtistTag: "Lewis OfMan", TitleTag: "Yo bene"},
			{TrackNumberTag: 3, ArtistTag: "Magnüm", TitleTag: "L'épée à la main"},
			{TrackNumberTag: 4, ArtistTag: "Bleu Toucan", TitleTag: "Hanoï Café"},
			{TrackNumberTag: 5, ArtistTag: "Vendredi sur Mer",
				TitleTag: "La femme à la peau bleue (Chez toi)",
			},
		})
}
