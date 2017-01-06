package main

import (
	"bufio"
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
	"os"
	"sort"
	"strings"
)

var printSection func(string, ...interface{}) = color.New(color.Bold).PrintfFunc()
var printTrack func(string, ...interface{}) = color.New(color.Bold, color.Italic).PrintfFunc()

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
	printSection("Album: %s\n", meta.AlbumString(r))
	if err := r.FixAlbum(fixTag); err != nil {
		return errors.Wrap(err, "album tag fix")
	}
	for _, t := range r.Tracks() {
		printTrack("Track: %s\n", meta.TrackString(t))
		if err := r.FixTrack(t, fixTag); err != nil {
			return errors.Wrapf(err, "track %d fix", t.TrackNumber())
		}
	}
	if list {
		printSection("Album Description\n")
		trackList(r)
		sourceList(sources)
	}
	return r.Pack(dest)
}
func yesNo() (bool, error) {
	fmt.Printf("[y/N] ")
	reader := bufio.NewReader(os.Stdin)
	c, err := reader.ReadByte()
	if err != nil {
		return false, errors.Wrap(err, "input")
	}
	return c == []byte("Y")[0] || c == []byte("y")[0], nil
}
func main() {
	if err := try(); err != nil {
		log.Fatalln("failure:", err)
	}
}
