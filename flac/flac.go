package flac

import (
	"fmt"
	lib "github.com/mewkiz/flac"
	libmeta "github.com/mewkiz/flac/meta"
	"github.com/mewspring/metautil"
	"github.com/pkg/errors"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/release"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func init() {
	release.RegisterFormat("flac", NewRelease)
}

const (
	tagAlbum       string = "ALBUM"
	tagDate               = "DATE"
	tagAlbumArtist        = "ALBUMARTIST"
	tagArtist             = "ARTIST"
	tagTitle              = "TITLE"
	tagTrackNumber        = "TRACKNUMBER"
)

func encodeTags(r *flac, f *file) [][2]string {
	return [][2]string{
		{tagAlbum, r.album},
		{tagDate, r.date},
		{tagAlbumArtist, r.albumArtist},
		{tagArtist, f.artist},
		{tagTitle, f.title},
		{tagTrackNumber, strconv.Itoa(int(f.trackNumber))},
	}
}

type file struct {
	album       string
	date        string
	albumArtist string
	artist      string
	title       string
	trackNumber uint16
	path        string
	stream      *lib.Stream
	vorbis      *libmeta.VorbisComment
}

func (f *file) Tags(t meta.Track) {
	f.artist = t.Artist()
	f.title = t.Title()
}
func (f *file) AlbumTitle() string {
	return f.album
}
func (f *file) Date() string {
	return f.date
}
func (f *file) AlbumArtist() string {
	return f.albumArtist
}
func (f *file) Artist() string {
	return f.artist
}
func (f *file) Title() string {
	return f.title
}
func (f *file) TrackNumber() uint16 {
	return f.trackNumber
}
func (f *file) rate() uint32 {
	return f.stream.Info.SampleRate
}
func (f *file) depth() uint8 {
	return f.stream.Info.BitsPerSample
}
func (f *file) Secs() int {
	return int(f.stream.Info.NSamples / uint64(f.rate()))
}

type flac struct {
	album       string
	date        string
	albumArtist string
	rate        uint32
	depth       uint8
	files       []*file
}

func (r *flac) Tags(t meta.Album) {
	r.album = t.AlbumTitle()
	r.date = t.Date()
	r.albumArtist = t.AlbumArtist()
}
func (r *flac) AlbumTitle() string {
	return r.album
}
func (r *flac) Date() string {
	return r.date
}
func (r *flac) AlbumArtist() string {
	return r.albumArtist
}
func (r *flac) Tracks() []release.Track {
	tracks := make([]release.Track, 0)
	for _, file := range r.files {
		tracks = append(tracks, file)
	}
	return tracks
}
func (r *flac) Format() string {
	return fmt.Sprintf("%d Hz %d bit", r.rate, r.depth)
}

func NewRelease(paths ...string) (release.Release, error) {
	count := len(paths)
	if count == 0 {
		return nil, errors.New("no tracks")
	}
	r := flac{files: make([]*file, count)}
	numbers := make([]uint16, count)
	for i, path := range paths {
		file, err := New(path)
		if err != nil {
			return nil, errors.Wrap(err, "can't open track")
		}
		r.files[i] = file
		if file.rate() > r.rate {
			r.rate = file.rate()
		}
		if file.depth() > r.depth {
			r.depth = file.depth()
		}
		numbers[i] = file.trackNumber
	}

	// album tags consistency check
	r.album = r.files[0].album
	r.date = r.files[0].date
	r.albumArtist = r.files[0].albumArtist
	for _, file := range r.files {
		if file.album != r.album ||
			file.date != r.date ||
			file.albumArtist != r.albumArtist {
			return nil, errors.New("this is not a single album")
		}
	}

	// album track numbers check
	sort.Sort(uint16Slice(numbers))
	prevNumber := uint16(0)
	for _, n := range numbers {
		if prevNumber+1 != n {
			return nil, errors.New("a track is missing from the album")
		}
		prevNumber++
	}
	return &r, nil
}

type uint16Slice []uint16

func (s uint16Slice) Len() int {
	return len(s)
}
func (s uint16Slice) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s uint16Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (r *flac) Pack(path string) error {
	pad := meta.TrackNumberPadder(uint16(len(r.files)))
	for _, file := range r.files {
		basename := fmt.Sprintf(
			"%s. %s.flac",
			pad(file.trackNumber),
			release.Filename(file.title),
		)
		w, err := os.Create(filepath.Join(path, basename))
		if err != nil {
			return errors.Wrap(err, "can't create target track file")
		}
		metautil.RemoveBlockType(file.stream, libmeta.TypePadding)
		metautil.RemoveBlockType(file.stream, libmeta.TypePicture)
		file.vorbis.Vendor = "mewkiz/flac"
		file.vorbis.Tags = encodeTags(r, file)
		if err := lib.Encode(w, file.stream); err != nil {
			return errors.Wrap(err, "can't encode track")
		}
	}
	return nil
}
func New(path string) (*file, error) {
	file := file{path: path}
	s, err := lib.ParseFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't open file")
	}
	file.stream = s
	var tagBlock *libmeta.Block
	for _, b := range s.Blocks {
		if b.Header.Type == libmeta.TypeVorbisComment {
			tagBlock = b
		}
	}
	if tagBlock == nil {
		return nil, errors.New("no tags")
	}
	file.vorbis = tagBlock.Body.(*libmeta.VorbisComment)
	m := make(map[string]string)
	for _, c := range file.vorbis.Tags {
		m[c[0]] = c[1]
	}
	file.artist = m[tagArtist]
	file.title = m[tagTitle]
	trackNumber, err := strconv.Atoi(m[tagTrackNumber])
	if err != nil || trackNumber <= 0 {
		return nil, errors.Wrap(err, "invalid track number")
	}
	file.trackNumber = uint16(trackNumber)
	file.album = m[tagAlbum]
	file.date = m[tagDate]
	file.albumArtist = m[tagAlbumArtist]
	return &file, nil
}
