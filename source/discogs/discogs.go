package discogs

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/source"
	"io"
	"net/http"
)

func init() {
	source.RegisterSource("discogs.com", Get)
}

type album struct {
	title   string
	artists []string
	tracks  []*track
	year    string
}
type track struct {
	*album
	artists []string
	title   string
	number  uint16
}

func (a *album) AlbumTitle() string {
	return a.title
}
func (a *album) AlbumArtist() string {
	return source.JoinArtists(a.artists)
}
func (a *album) Date() string {
	return a.year
}
func (t *track) Artist() string {
	return source.JoinArtists(t.artists)
}
func (t *track) Title() string {
	return t.title
}
func (t *track) TrackNumber() uint16 {
	return t.number
}
func parseArtists(a *album, s *goquery.Selection) error {
	artists := s.Find("[itemprop=name]")
	a.artists = make([]string, 0, artists.Size())
	var err error
	artists.EachWithBreak(func(i int, s *goquery.Selection) bool {
		title, exists := s.Attr("title")
		if !exists {
			err = errors.New("title")
			return false
		}
		a.artists = append(a.artists, title)
		return true
	})
	return errors.Wrap(err, "artist")
}

func parseProfile(a *album, s *goquery.Selection) error {
	artists := s.Find("#profile_title [itemprop=byArtist]")
	if artists.Size() == 0 {
		return errors.New("no artists")
	}
	if err := parseArtists(a, artists); err != nil {
		return errors.Wrap(err, "artists")
	}

	name := s.Find("#profile_title > [itemprop=name]")
	if name.Size() != 1 {
		return errors.New("no title")
	}
	a.title = source.TrimText(name.Text())

	props := s.Find(".content")
	if props.Size() != 6 {
		return errors.Errorf("%d proprieties", props.Size())
	}

	date := goquery.NewDocumentFromNode(props.Get(3)).Find("a")
	if date.Size() != 1 {
		return errors.New("date not found")
	}
	year, err := source.Year(date.Text())
	if err != nil {
		return errors.Wrap(err, "date")
	}
	a.year = year
	return nil
}
func parseTrackArtists(t *track, s *goquery.Selection) error {
	artists := s.Find("a")
	if artists.Size() == 0 {
		return errors.New("no artist")
	}
	t.artists = make([]string, 0, artists.Size())
	var err error
	artists.EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Size() != 1 {
			err = errors.New("text")
			return false
		}
		t.artists = append(t.artists, s.Text())
		return true
	})
	return errors.Wrap(err, "artist")
}
func parseTrack(t *track, s *goquery.Selection) error {
	artists := s.Find(".tracklist_track_artists")
	if s := artists.Size(); s == 1 {
		if err := parseTrackArtists(t, artists); err != nil {
			return errors.Wrap(err, "various artists")
		}
	} else if s != 0 {
		return errors.New("artists")
	}

	title := s.Find("[itemprop=name]")
	if title.Size() != 1 {
		return errors.New("title")
	}
	t.title = title.Text()
	return nil
}
func parseTracks(a *album, s *goquery.Selection) error {
	tracks := s.Find(".tracklist_track")
	count := tracks.Size()
	if count == 0 {
		return errors.New("no tracks")
	}
	a.tracks = make([]*track, count)
	var err error
	tracks.EachWithBreak(func(i int, s *goquery.Selection) bool {
		a.tracks[i] = &track{album: a, number: uint16(i + 1), artists: a.artists}
		err = parseTrack(a.tracks[i], s)
		return err == nil
	})
	return errors.Wrap(err, "track")
}
func parse(a *album, s *goquery.Selection) error {
	profile := s.Find(".profile")
	if profile.Size() != 1 {
		return errors.Errorf("profile nodes: %d", profile.Size())
	}
	if err := parseProfile(a, profile); err != nil {
		return errors.Wrap(err, "profile")
	}
	if err := parseTracks(a, s); err != nil {
		return errors.Wrap(err, "tracks")
	}
	return nil
}
func Get(url string) (meta.Album, []meta.Track, error) {
	req, err := http.NewRequest("GET", url, nil)

	// apparently required
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:50.0) Gecko/20100101 "+
		"Firefox/50.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	return Read(resp.Body)
}
func Read(reader io.Reader) (meta.Album, []meta.Track, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, nil, errors.Wrap(err, "document")
	}
	page := doc.Find("#page_content")
	if len(page.Nodes) != 1 {
		return nil, nil, errors.Errorf("page %d", len(page.Nodes))
	}
	a := album{}
	if err := parse(&a, page); err != nil {
		return nil, nil, err
	}

	tracks := make([]meta.Track, len(a.tracks))
	for i, track := range a.tracks {
		tracks[i] = track
	}
	return &a, tracks, nil
}
