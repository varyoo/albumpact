package sevenDigital

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"regexp"
)

func init() {
	source.RegisterSource("7digital.com", Get)
}
func album(a *meta.AlbumTags, s *goquery.Selection) error {
	title := s.Find("h1.release-info-title")
	if len(title.Nodes) != 1 {
		return errors.New("can't find title")
	}
	a.AlbumTag = title.Text()

	artist := s.Find("h2.release-info-artist [itemprop=byArtist] [itemprop=name]")
	if len(artist.Nodes) != 1 {
		return errors.New("can't find artist")
	}
	artistText, exists := artist.Attr("content")
	if !exists {
		return errors.New("can't parse artist")
	}
	if artistText == "Various Artists" {
		label := s.Find(".release-label-info .release-data-info")
		if len(label.Nodes) != 1 {
			return errors.New("can't find label")
		}
		artistText = label.Text()
	}
	a.AlbumArtistTag = artistText

	date := s.Find(".release-date-info .release-data-info")
	if len(date.Nodes) != 1 {
		return errors.New("can't find date")
	}
	r := regexp.MustCompile("([0-9]{4})")
	m := r.FindStringSubmatch(date.Text())
	if len(m) != 2 {
		return errors.Errorf("can't extract date: %s", date.Text())
	}
	a.DateTag = m[1]
	return nil
}
func track(t *meta.TrackTags, s *goquery.Selection) error {
	name := s.Find(".release-track-name [itemprop=name]")
	if len(name.Nodes) != 1 {
		return errors.New("can't find title")
	}
	if nameContent, exists := name.Attr("content"); !exists {
		return errors.New("can't parse title")
	} else {
		t.TitleTag = nameContent
	}
	va := s.Find(".release-track-list-additional a")
	if len(va.Nodes) == 1 {
		// various artists
		t.ArtistTag = source.VASplit(va.Text(), ",")
	} else {
		t.ArtistTag = t.AlbumArtist()
	}
	return nil
}
func Get(url string) (meta.Album, []meta.Track, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't open document")
	}
	s := doc.Find(".release")
	if len(s.Nodes) != 1 {
		return nil, nil, errors.New("can't find release")
	}
	releaseInfo := s.Find(".release-info")
	if len(releaseInfo.Nodes) != 1 {
		return nil, nil, errors.New("can't find release info")
	}
	a := &meta.AlbumTags{}
	if err := album(a, releaseInfo); err != nil {
		return nil, nil, errors.Wrap(err, "album")
	}
	var trackErr error
	tracks := make([]meta.Track, 0)
	s.Find(".release-track-list .release-track").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			t := &meta.TrackTags{
				Album:          a,
				TrackNumberTag: uint16(i + 1),
			}
			if err := track(t, s); err != nil {
				trackErr = errors.Wrap(err, "track")
				return false
			}
			tracks = append(tracks, t)
			return true
		})
	if trackErr != nil {
		return nil, nil, trackErr
	}
	return a, tracks, nil
}
