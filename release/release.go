package release

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"os"
	"path/filepath"
	"strings"
)

var formats []format

type Release interface {
	AlbumTitle() string
	Date() string
	AlbumArtist() string
	Pack(path string) error
	Format() string
	Tracks() []Track
	Tags(t meta.Album)
}
type Track interface {
	AlbumTitle() string
	Date() string
	AlbumArtist() string
	Artist() string
	Title() string
	TrackNumber() uint16
	Secs() int
	Tags(t meta.Track)
}
type release struct {
	Release
}

func (r *release) Pack(path string) error {
	basename := fmt.Sprintf("%s - %s [%s]", r.AlbumArtist(), r.AlbumTitle(), r.Format())
	basename = Filename(basename)
	path = filepath.Join(path, basename)
	if err := os.Mkdir(path, 0755); err != nil {
		return errors.Wrap(err, "can't create release directory")
	}
	return r.Release.Pack(path)
}

type formatCreator func(path ...string) (Release, error)
type format struct {
	name   string
	create formatCreator
}

func (f *format) NewRelease(paths ...string) (Release, error) {
	if r, err := f.create(paths...); err != nil {
		return nil, err
	} else {
		return &release{r}, nil
	}
}
func RegisterFormat(name string, create formatCreator) {
	formats = append(formats, format{name, create})
}
func NewRelease(paths ...string) (Release, error) {
	status := errors.New("unsuported format")
	for _, format := range formats {
		if r, err := format.NewRelease(paths...); err == nil {
			autoCorrect(r)
			return r, nil
		} else {
			status = errors.Wrap(err, "flac")
		}
	}
	return nil, errors.Wrap(status, "format")
}

func autoCorrect(r Release) {
	at := &meta.AlbumTags{
		AlbumTag:       r.AlbumTitle(),
		AlbumArtistTag: r.AlbumArtist(),
		DateTag:        r.Date(),
	}
	tracks := r.Tracks()
	if r.AlbumArtist() == "" {
		var artist string
		sameArtist := true
		for _, t := range tracks {
			if artist == "" {
				artist = t.Artist()
			} else if t.Artist() != artist {
				sameArtist = false
				break
			}
		}
		if sameArtist {
			at.AlbumArtistTag = artist
		}

	}
	r.Tags(at)
}
func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
func Filename(s string) string {
	return stripchars(s, "<>:\"/\\|?*/")
}
