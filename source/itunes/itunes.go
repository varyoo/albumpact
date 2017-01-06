package itunes

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"regexp"
)

func init() {
	source.RegisterSource("apple.com", Get)
}
func track(t *meta.TrackTags, s *goquery.Selection) error {
	name := s.Find(".name .text")
	if len(name.Nodes) != 1 {
		return errors.New("can't find title")
	}
	t.TitleTag = name.Text()

	artist := s.Find(".artist .text")
	if len(artist.Nodes) != 1 {
		return errors.New("can't find artist")
	}
	t.ArtistTag = artist.Text()
	return nil
}
func album(t *meta.AlbumTags, s *goquery.Selection) error {
	title := s.Find("#title")
	h1 := title.Find("h1")
	if len(h1.Nodes) != 1 {
		return errors.New("can't find album title")
	}
	t.AlbumTag = h1.Text()

	h2 := title.Find("h2")
	if len(h2.Nodes) != 1 {
		return errors.New("can't find album artist")
	}
	t.AlbumArtistTag = h2.Text()

	date := s.Find(".release-date [itemprop=dateCreated]")
	if len(date.Nodes) != 1 {
		return errors.New("can't find release date")
	}
	r := regexp.MustCompile("([0-9]{4})")
	m := r.FindStringSubmatch(date.Text())
	if len(m) != 2 {
		return errors.New("can't parse release date")
	}
	t.DateTag = m[1]
	return nil
}
func Get(url string) (meta.Album, []meta.Track, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't open page")
	}
	content := doc.Find("#content")
	if len(content.Nodes) != 1 {
		return nil, nil, errors.New("can't find content")
	}
	a := &meta.AlbumTags{}
	if err := album(a, content); err != nil {
		return nil, nil, errors.Wrap(err, "album")
	}
	tracks := make([]meta.Track, 0)
	var trackErr error
	doc.Find(".track-list.album.music table.tracklist-table tr.song").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			t := &meta.TrackTags{
				Album:          a,
				TrackNumberTag: uint16(i + 1),
			}
			if trackErr = track(t, s); trackErr != nil {
				return false
			}
			tracks = append(tracks, t)
			return true
		})
	return a, tracks, errors.Wrap(trackErr, "track")
}
