package source

import (
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func VASplit(s, sep string) string {
	rawArtists := strings.Split(s, sep)
	artists := make([]string, len(rawArtists))
	for i, a := range rawArtists {
		a = strings.Trim(a, " ")
		artists[i] = a
	}
	return JoinArtists(artists)
}
func JoinArtists(artists []string) string {
	sort.Sort(sort.StringSlice(artists))
	return strings.Join(artists, " & ")
}
func SourceTest(t *testing.T,
	get func(url string) (meta.Album, []meta.Track, error), url string,
	album meta.AlbumTags, tracks []meta.TrackTags) {
	a, l, err := get(url)
	if err != nil {
		t.Fatal(err)
	}
	if !meta.AlbumEq(&album, a) {
		t.Error(meta.AlbumString(a))
	}
	if len(tracks) != len(l) {
		t.Errorf("bad track count, expected %d but have %d", len(tracks), len(l))
	}
	for i, li := range l {
		if !meta.TrackEq(&tracks[i], li) {
			t.Error(meta.TrackString(li))
		}
	}
}
func TrimText(s string) string {
	return strings.Trim(strings.Replace(s, "\n", " ", -1), " ")
}
func Year(date string) (string, error) {
	r := regexp.MustCompile("([0-9]{4})")
	m := r.FindStringSubmatch(date)
	if len(m) != 2 {
		return "", errors.Errorf("not found in %s", date)
	}
	return m[1], nil
}
