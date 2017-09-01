package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/robbiev/dilemma"
	"github.com/varyoo/albumpact"
	"github.com/varyoo/albumpact/meta"
	"github.com/varyoo/albumpact/release"
	"github.com/varyoo/albumpact/stats"
	"log"
	"sort"
	"strings"
)

type byNumber []release.Track

func (s byNumber) Len() int {
	return len(s)
}
func (s byNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byNumber) Less(i, j int) bool {
	return s[i].TrackNumber() < s[j].TrackNumber()
}

func trackList(r release.Release) {
	tracks := r.Tracks()
	sort.Sort(byNumber(tracks))
	for _, track := range tracks {
		secs := track.Secs()
		mm := secs / 60
		ss := secs - (mm * 60)
		fmt.Printf("[#]%s [i](%02d:%02d)[/i]\n", track.Title(), mm, ss)
	}
}
func sourceList(sources []string) {
	for _, s := range sources {
		if s != "" {
			fmt.Printf("[*]%s\n", s)
		}
	}
}
func fixTag(tag meta.Tag, tags []*stats.StringTag) (string, error) {
	options := make([]string, 0, len(tags))
	m := make(map[string]string)
	for _, t := range tags {
		opt := t.String()
		options = append(options, opt)
		m[opt] = t.Val()
	}
	selected, exitKey, err := dilemma.Prompt(dilemma.Config{
		Title:   fmt.Sprintf("Multiple values for tag: %s, choose one:", tag),
		Help:    "Use arrow up and down, then enter to select.",
		Options: options,
	})
	if err != nil {
		return "", errors.Wrap(err, "input")
	}
	if exitKey != dilemma.Empty {
		return "", errors.New("no choice given")
	}
	return m[selected], nil
}

var noVal = color.New(color.Italic).SprintfFunc()

type artist string
type title string
type date string

func (a artist) String() string {
	return valOrMsg(string(a), "Unknown artist")
}
func (a title) String() string {
	return valOrMsg(string(a), "Untitled")
}
func (a date) String() string {
	return valOrMsg(string(a), "Unknown date")
}
func valOrMsg(val, msg string) string {
	if val == "" {
		return noVal(msg)
	} else {
		return val
	}
}

type track struct {
	meta.Track
}

func (t track) String() string {
	return fmt.Sprintf("%d. %s — %s", t.TrackNumber(), artist(t.Artist()), title(t.Title()))
}

type album struct {
	meta.Album
}

func (a album) String() string {
	return fmt.Sprintf("%s — %s (%s)", artist(a.AlbumArtist()), title(a.AlbumTitle()), date(a.Date()))
}

var bullet = color.New(color.FgBlue, color.Bold).SprintfFunc()

func try() error {
	var dest string
	var list bool
	var rawSources string
	flag.StringVar(&dest, "o", "./", "destination directory")
	flag.BoolVar(&list, "l", false, "list tracks")
	flag.StringVar(&rawSources, "s", "", "itunes.com/...,qobuz.com/...")
	flag.Parse()
	paths := flag.Args()
	r, err := albumpact.NewRelease(paths...)
	if err != nil {
		return errors.Wrap(err, "can't load the release")
	}
	sources := make([]string, 0)
	if rawSources != "" {
		sources = strings.Split(rawSources, ",")
		if err := r.Source(sources...); err != nil {
			return err
		}
	}
	fmt.Println(bullet("==>"), album{r})
	if err := r.FixAlbum(fixTag); err != nil {
		return errors.Wrap(err, "album tag fix")
	}
	for _, t := range r.Tracks() {
		fmt.Println(bullet("  ->"), track{t})
		if err := r.FixTrack(t, fixTag); err != nil {
			return errors.Wrapf(err, "track %d fix", t.TrackNumber())
		}
	}
	if list {
		trackList(r)
		sourceList(sources)
	}
	return r.Pack(dest)
}
func main() {
	if err := try(); err != nil {
		log.Fatalln("failure:", err)
	}
}
