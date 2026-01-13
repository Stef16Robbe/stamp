package adr

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Store struct {
	Directory string
}

func NewStore(directory string) *Store {
	return &Store{Directory: directory}
}

var filenameRegex = regexp.MustCompile(`^(\d{4})-(.+)\.md$`)

func (s *Store) List() ([]*ADR, error) {
	entries, err := os.ReadDir(s.Directory)
	if err != nil {
		return nil, err
	}

	var adrs []*ADR
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !filenameRegex.MatchString(entry.Name()) {
			continue
		}

		adr, err := s.Load(entry.Name())
		if err != nil {
			continue
		}

		adrs = append(adrs, adr)
	}

	sort.Slice(adrs, func(i, j int) bool {
		return adrs[i].Number < adrs[j].Number
	})

	return adrs, nil
}

func (s *Store) Load(filename string) (*ADR, error) {
	path := filepath.Join(s.Directory, filename)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	adr, err := ParseMarkdown(string(content))
	if err != nil {
		return nil, err
	}

	adr.Filename = filename
	return adr, nil
}

func (s *Store) Save(adr *ADR) error {
	if adr.Filename == "" {
		adr.Filename = FormatFilename(adr.Number, adr.Title)
	}

	path := filepath.Join(s.Directory, adr.Filename)
	content := adr.ToMarkdown()

	return os.WriteFile(path, []byte(content), 0644)
}

func (s *Store) NextNumber() (int, error) {
	entries, err := os.ReadDir(s.Directory)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, err
	}

	maxNum := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		match := filenameRegex.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}

		num, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}

		if num > maxNum {
			maxNum = num
		}
	}

	return maxNum + 1, nil
}

func (s *Store) FindByNumber(number int) (*ADR, error) {
	adrs, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, adr := range adrs {
		if adr.Number == number {
			return adr, nil
		}
	}

	return nil, os.ErrNotExist
}
