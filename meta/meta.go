package meta

import (
	"fmt"
	"strconv"
	"strings"
)

type Tag string

const (
	AlbumTag       Tag = "album"
	DateTag            = "date"
	AlbumArtistTag     = "album artist"
	ArtistTag          = "artist"
	TitleTag           = "title"
	TrackNumberTag     = "track number"
)

type Album interface {
	AlbumTitle() string
	Date() string
	AlbumArtist() string
}

func AlbumString(a Album) string {
	return fmt.Sprintf("%s - %s (%s)", a.AlbumArtist(), a.AlbumTitle(), a.Date())
}

type Track interface {
	AlbumTitle() string
	Date() string
	AlbumArtist() string
	Artist() string
	Title() string
	TrackNumber() uint16
}

func TrackString(t Track) string {
	return fmt.Sprintf("%d. %s - %s", t.TrackNumber(), t.Artist(), t.Title())
}

type AlbumTags struct {
	AlbumTag       string
	AlbumArtistTag string
	DateTag        string
}

func AlbumEq(a, b Album) bool {
	return a.AlbumTitle() == b.AlbumTitle() && a.AlbumArtist() == b.AlbumArtist() &&
		a.Date() == b.Date()
}
func (t *AlbumTags) AlbumTitle() string {
	return t.AlbumTag
}
func (t *AlbumTags) Date() string {
	return t.DateTag
}
func (t *AlbumTags) AlbumArtist() string {
	return t.AlbumArtistTag
}

type TrackTags struct {
	Album
	TrackNumberTag uint16
	ArtistTag      string
	TitleTag       string
}

func TrackEq(a, b Track) bool {
	return a.Title() == b.Title() && a.TrackNumber() == b.TrackNumber() &&
		a.Artist() == b.Artist()
}
func (t *TrackTags) Title() string {
	return t.TitleTag
}
func (t *TrackTags) Artist() string {
	return t.ArtistTag
}
func (t *TrackTags) TrackNumber() uint16 {
	return t.TrackNumberTag
}
func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}
func TrackNumberPadder(trackCount uint16) func(uint16) string {
	zeroes := len([]rune(strconv.Itoa(int(trackCount))))
	return func(trackNumber uint16) string {
		return leftPad2Len(strconv.Itoa(int(trackNumber)), "0", zeroes)
	}
}
