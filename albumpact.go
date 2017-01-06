package albumpact

import (
	"github.com/pkg/errors"
	_ "github.com/varyoo/albumpact/flac"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/release"
	"github.com/varyoo/albumpact/source"
	_ "github.com/varyoo/albumpact/source/7digital"
	_ "github.com/varyoo/albumpact/source/discogs"
	_ "github.com/varyoo/albumpact/source/itunes"
	_ "github.com/varyoo/albumpact/source/qobuz"
	"github.com/varyoo/albumpact/stats"
)

type Release struct {
	release.Release
	stats      stats.AlbumStats
	trackStats map[uint16]stats.TrackStats
}

func NewRelease(paths ...string) (*Release, error) {
	eng, err := release.NewRelease(paths...)
	if err != nil {
		return nil, err
	}
	r := &Release{Release: eng, stats: stats.NewAlbumStats()}
	r.stats.Add(r)

	r.trackStats = make(map[uint16]stats.TrackStats)
	for _, t := range eng.Tracks() {
		ts := stats.NewTrackStats()
		ts.Add(t)
		r.trackStats[t.TrackNumber()] = ts
	}
	return r, nil
}
func (r *Release) Source(sources ...string) error {
	for _, s := range sources {
		a, l, err := source.Get(s)
		if err != nil {
			return errors.Wrap(err, "source")
		}
		r.stats.Add(a)

		for _, t := range l {
			r.trackStats[t.TrackNumber()].Add(t)
		}
	}
	return nil
}
func mergeTags(tag meta.Tag, tags stats.StringTags,
	ask func(meta.Tag, []*stats.StringTag) (string, error),
	set func(string)) error {
	values := tags.Slice()
	count := len(values)
	if count == 1 {
		set(values[0].Val())
	} else if count > 1 {
		if val, err := ask(tag, values); err != nil {
			return errors.Wrap(err, "manual tag fix")
		} else {
			set(val)
		}
	}
	return nil
}
func (r *Release) FixAlbum(ask func(meta.Tag, []*stats.StringTag) (string, error)) error {
	t := &meta.AlbumTags{
		AlbumTag:       r.AlbumTitle(),
		AlbumArtistTag: r.AlbumArtist(),
		DateTag:        r.Date(),
	}
	if err := mergeTags(meta.AlbumTag, r.stats.Albums, ask, func(v string) {
		t.AlbumTag = v
	}); err != nil {
		return errors.Wrap(err, "album")
	}
	if err := mergeTags(meta.AlbumArtistTag, r.stats.AlbumArtists, ask, func(v string) {
		t.AlbumArtistTag = v
	}); err != nil {
		return errors.Wrap(err, "album artist")
	}
	if err := mergeTags(meta.DateTag, r.stats.Dates, ask, func(v string) {
		t.DateTag = v
	}); err != nil {
		return errors.Wrap(err, "date")
	}
	r.Tags(t)
	return nil
}

type fixStringTag func(meta.Tag, []*stats.StringTag) (string, error)

func (r *Release) FixTrack(track release.Track, ask fixStringTag) error {
	tn := track.TrackNumber()
	t := &meta.TrackTags{
		Album:          r,
		TrackNumberTag: tn,
		TitleTag:       track.Title(),
		ArtistTag:      track.Artist(),
	}
	ts := r.trackStats[tn]
	if err := mergeTags(meta.TitleTag, ts.Titles, ask, func(v string) {
		t.TitleTag = v
	}); err != nil {
		return errors.Wrap(err, "title")
	}
	if err := mergeTags(meta.ArtistTag, ts.Artists, ask, func(v string) {
		t.ArtistTag = v
	}); err != nil {
		return errors.Wrap(err, "artist")
	}
	track.Tags(t)
	return nil
}
func (r *Release) FixTracks(ask fixStringTag) error {
	for _, t := range r.Tracks() {
		if err := r.FixTrack(t, ask); err != nil {
			return errors.Wrapf(err, "track %d", t.TrackNumber())
		}
	}
	return nil
}
