package sevenDigital

import (
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"testing"
)

func TestCeline(t *testing.T) {
	tracks := []meta.TrackTags{
		{TrackNumberTag: 1, ArtistTag: "Céline Tolosa", TitleTag: "Cover Girl"},
		{TrackNumberTag: 2, ArtistTag: "Céline Tolosa", TitleTag: "Rue Mansart"},
		{TrackNumberTag: 3, ArtistTag: "Céline Tolosa", TitleTag: "Tu es fantastique"},
		{TrackNumberTag: 4, ArtistTag: "Céline Tolosa", TitleTag: "Fais-moi souffrir"},
	}
	album := meta.AlbumTags{
		AlbumTag:       "Cover Girl - EP",
		AlbumArtistTag: "Céline Tolosa",
		DateTag:        "2015",
	}
	source.SourceTest(t, Get,
		"https://www.7digital.com/artist/celine-tolosa/release/cover-girl", album, tracks)
}
func TestCercleSubaquatique(t *testing.T) {
	album := meta.AlbumTags{
		AlbumTag:       "Cercle subaquatique Profil de Face - EP",
		AlbumArtistTag: "Profil de Face Records",
		DateTag:        "2016",
	}
	tracks := []meta.TrackTags{
		{TrackNumberTag: 1, ArtistTag: "Février", TitleTag: "Château rouge"},
		{TrackNumberTag: 2, ArtistTag: "Lewis OfMan", TitleTag: "Yo bene"},
		{TrackNumberTag: 3, ArtistTag: "Magnüm", TitleTag: "L'épée à la main"},
		{TrackNumberTag: 4, ArtistTag: "Bleu Toucan", TitleTag: "Hanoï Café"},
		{TrackNumberTag: 5, ArtistTag: "Vendredi sur Mer",
			TitleTag: "La femme à la peau bleue (Chez toi)",
		},
	}
	source.SourceTest(t, Get,
		"https://fr.7digital.com/artist/various-artists/release/cercle-subaquatique"+
			"-profil-de-face-ep-5472408", album, tracks)
}
