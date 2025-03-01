package domain

import (
	"slices"
	"strings"
)

const (
	DungeonsFile = "dungeons.json"
	QuestsFile   = "quests.json"
	SitemapFile  = "sitemap.json"
	TitlesFile   = "titles.json"
)

type TypeKey string

const (
	QuestsKey   TypeKey = "quest"
	DungeonsKey TypeKey = "dungeon"
)

type ManualFile struct {
	Dungeons ManualKey `yaml:"dungeons"`
	Quests   ManualKey `yaml:"quests"`
}

type ManualKey struct {
	UnmappedId     []int                 `yaml:"unmapped_id"`
	RemappedId     map[int]string        `yaml:"remapped_id"`
	RemappedPrefix map[string]string     `yaml:"remapped_prefix"`
	RewriteTarget  []ManuelRewriteTarget `yaml:"rewrite_target"`
}

type ManuelRewriteTarget struct {
	FromSuffix string `yaml:"from_suffix"`
	ToSuffix   string `yaml:"to_suffix"`
	Group      string `yaml:"group"` // Only one rewrite target per group will be applied
}

func (mf *ManualFile) GetIfRemapped(key TypeKey, id int, title string) (string, bool) {
	checkIfRemappedId := func(remappedId map[int]string, remappedPrefix map[string]string, id int, title string) (string, bool) {
		if url, ok := remappedId[id]; ok {
			return url, true
		}
		for prefix, url := range remappedPrefix {
			if strings.HasPrefix(title, prefix) {
				return url, true
			}
		}
		return "", false
	}

	if key == QuestsKey {
		return checkIfRemappedId(mf.Quests.RemappedId, mf.Quests.RemappedPrefix, id, title)
	} else if key == DungeonsKey {
		return checkIfRemappedId(mf.Dungeons.RemappedId, mf.Dungeons.RemappedPrefix, id, title)
	}
	return "", false
}

func (mf *ManualFile) GetIfUnmapped(key TypeKey, id int) bool {
	checkIfUnmapped := func(unmapped []int, id int) bool {
		return slices.Contains(unmapped, id)
	}
	if key == QuestsKey {
		return checkIfUnmapped(mf.Quests.UnmappedId, id)
	} else if key == DungeonsKey {
		return checkIfUnmapped(mf.Dungeons.UnmappedId, id)
	}
	return false
}

func (mf *ManualFile) RewriteTarget(key TypeKey, target string) string {
	checkIfRewriteTarget := func(rewriteTargets []ManuelRewriteTarget, target string) string {
		var appliedGroups []string
		for _, rewriteTarget := range rewriteTargets {
			if slices.Contains(appliedGroups, rewriteTarget.Group) {
				continue
			}
			if len(rewriteTarget.FromSuffix) > 0 && strings.HasSuffix(target, rewriteTarget.FromSuffix) {
				target = strings.ReplaceAll(target, rewriteTarget.FromSuffix, rewriteTarget.ToSuffix)
				appliedGroups = append(appliedGroups, rewriteTarget.Group)
			}
		}
		return target
	}

	if key == QuestsKey {
		return checkIfRewriteTarget(mf.Quests.RewriteTarget, target)
	} else if key == DungeonsKey {
		return checkIfRewriteTarget(mf.Dungeons.RewriteTarget, target)
	}
	return target
}
