package stats

import (
	"fmt"
	"github.com/varyoo/albumpact/meta"
	"sort"
)

type StringTag struct {
	val  string
	seen uint8
}

func (t *StringTag) String() string {
	return fmt.Sprintf("%s (%d times)", t.val, t.seen)
}
func (t *StringTag) Val() string {
	return t.val
}

type StringTags map[string]*StringTag

func (ts StringTags) add(val string) {
	if val != "" {
		if t := ts[val]; t == nil {
			ts[val] = &StringTag{val: val, seen: 1}
		} else {
			ts[val].seen++
		}
	}
}
func (ts StringTags) Slice() []*StringTag {
	tags := make([]*StringTag, 0, len(ts))
	for _, v := range ts {
		tags = append(tags, v)
	}
	sort.Sort(tagSlice(tags))
	return tags
}
func mkTags() StringTags {
	return make(StringTags)
}

type AlbumStats struct {
	Albums       StringTags
	Dates        StringTags
	AlbumArtists StringTags
}
type TrackStats struct {
	Titles  StringTags
	Artists StringTags
}

func NewAlbumStats() AlbumStats {
	return AlbumStats{mkTags(), mkTags(), mkTags()}
}
func NewTrackStats() TrackStats {
	return TrackStats{mkTags(), mkTags()}
}
func (s AlbumStats) Add(t meta.Album) {
	s.Albums.add(t.AlbumTitle())
	s.Dates.add(t.Date())
	s.AlbumArtists.add(t.AlbumArtist())
}
func (s TrackStats) Add(t meta.Track) {
	s.Titles.add(t.Title())
	s.Artists.add(t.Artist())
}

type tagSlice []*StringTag

func (s tagSlice) Len() int           { return len(s) }
func (s tagSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s tagSlice) Less(i, j int) bool { return s[i].seen > s[j].seen }
