package qobuz

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"regexp"
	"strings"
)

func init() {
	source.RegisterSource("qobuz.com", Get)
}
func track(t *meta.TrackTags, s *goquery.Selection) error {
	title := s.Find(".title .track-title")
	if len(title.Nodes) != 1 {
		return errors.New("can't find title")
	}
	t.TitleTag = source.TrimText(title.Text())

	details := s.Find(".track-details")
	if len(details.Nodes) != 1 {
		return errors.New("can't find details")
	}
	artist, err := findTrackArtistSpan(details)
	if err != nil {
		return errors.Wrap(err, "artist")
	}
	artists := strings.Split(artist.Text(), ",")
	if len(artists) == 0 {
		return errors.New("can't parse artist")
	}
	t.ArtistTag = artists[0]
	return nil
}
func findTrackArtistSpan(s *goquery.Selection) (*goquery.Selection, error) {
	var artist *goquery.Selection
	s.Find("span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.HasClass("copyright") {
			return true
		} else {
			artist = s
			return false
		}
	})
	if artist == nil {
		return nil, errors.New("not found")
	}
	if len(artist.Nodes) != 1 {
		return nil, errors.New("can't find")
	}
	return artist, nil
}
func album(t *meta.AlbumTags, s *goquery.Selection) error {
	h1 := s.Find("h1")
	if len(h1.Nodes) != 1 {
		return errors.New("can't find title")
	}
	t.AlbumTag = h1.Text()

	h2 := s.Find("h2")
	if len(h2.Nodes) != 1 {
		return errors.New("can't find artist")
	}
	t.AlbumArtistTag = h2.Text()

	date := s.Find("p.txt04 span")
	if len(date.Nodes) != 1 {
		return errors.New("can't find date")
	}
	r := regexp.MustCompile("([0-9]{4})")
	m := r.FindStringSubmatch(date.Text())
	if len(m) != 2 {
		return errors.New("can't parse date")
	}
	t.DateTag = m[1]
	return nil
}
func Get(url string) (meta.Album, []meta.Track, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't open page")
	}
	info := doc.Find("#info .meta")
	if len(info.Nodes) != 1 {
		return nil, nil, errors.New("can't find album info area")
	}
	a := &meta.AlbumTags{}
	if err := album(a, info); err != nil {
		return nil, nil, errors.Wrap(err, "album")
	}
	tracks := make([]meta.Track, 0)
	var trackErr error
	doc.Find("#tracklisting ol.tracks .track").
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
