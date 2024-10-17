package comet

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	RHEL8 = "rhel8"
)

var (
	cometFileLock sync.Mutex
)

type Comet struct {
	To   string      `json:"to" yaml:"to"`
	From []Component `json:"from" yaml:"from"`
}

type Component struct {
	Name string `json:"name" yaml:"name"`
	Repo string `json:"repo" yaml:"repo"`
}

type Guess struct {
	FilePath    string
	RHELVersion string
	Image       string
}

func GuessComet(g Guess) (*Comet, error) {
	cometFileLock.Lock()
	defer cometFileLock.Unlock()

	comets, err := loadComets(g.FilePath)
	if err != nil {
		return nil, err
	}
	image := imageName(g.Image)

	shortest := math.MaxInt
	idx := -1
	for i, comet := range comets {
		to := imageName(comet.To)
		if !strings.Contains(to, g.RHELVersion) {
			continue
		}
		for _, x := range comet.From {
			from := imageName(x.Repo)
			if dist := DistanceForStrings([]rune(image), []rune(from), DefaultOptions); dist < shortest {
				shortest = dist
				idx = i
			}
		}

		if dist := DistanceForStrings([]rune(to), []rune(image), DefaultOptions); dist < shortest {
			shortest = dist
			idx = i
		}
	}

	if idx < 0 {
		bytes, _ := yaml.Marshal(comets)
		return nil, fmt.Errorf("failed to find comet mapping in file %q: %v", g.FilePath, string(bytes))
	}

	return &comets[idx], nil
}

func Persist(filePath string, to string, from Component) error {
	cometFileLock.Lock()
	defer cometFileLock.Unlock()

	comets, err := loadComets(filePath)
	if err != nil {
		return err
	}

	found := false
	for i, comet := range comets {
		if comet.To == to {
			found = false
			for _, f := range comet.From {
				if f.Name == from.Name {
					found = true
					break
				}
			}
			if !found {
				found = true
				comets[i].From = append(comets[i].From, from)
			}
			break
		}
	}
	bytes, _ := yaml.Marshal(comets)
	if !found {
		return fmt.Errorf("comet 'to' mapping not found for %q in file %q, content:\n%v", to, filePath, string(bytes))
	}

	if err := os.WriteFile(filePath, bytes, 0666); err != nil {
		return fmt.Errorf("failed to persist comet to %v: %w", filePath, err)
	}

	return nil
}

func loadComets(filePath string) ([]Comet, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read comet.yaml: %w", err)
	}

	comets := make([]Comet, 0)
	if err := yaml.Unmarshal(bytes, &comets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comet.yaml: %w", err)
	}
	return comets, nil
}

func imageName(image string) string {
	idx := -1
	if i := strings.IndexRune(image, '/') + 1; i >= 0 {
		idx = i
	}
	return image[idx:]
}
