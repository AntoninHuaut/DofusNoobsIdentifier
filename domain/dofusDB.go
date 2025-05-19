package domain

import (
	"regexp"
	"strconv"
)

var (
	regexAlignmentType  = regexp.MustCompile(`Ps=(\d+)`)
	regexAlignmentLevel = regexp.MustCompile(`Pa=(\d+)`)
)

type DofusDbSearchResource[T any] struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
	Data  []T `json:"data"`
}

type HasID interface {
	GetID() int
}

type DofusDbDungeonLight struct {
	ID   int               `json:"id"`
	Name map[string]string `json:"name"`
}

func (d DofusDbDungeonLight) GetID() int {
	return d.ID
}

type DofusDbQuestLight struct {
	ID             int               `json:"id"`
	Name           map[string]string `json:"name"`
	StartCriterion string            `json:"startCriterion"`
}

func (q DofusDbQuestLight) GetID() int {
	return q.ID
}

func (q DofusDbQuestLight) IsAlignment() bool {
	matches := regexAlignmentType.FindStringSubmatch(q.StartCriterion)
	return len(matches) > 1
}

func (q DofusDbQuestLight) GetAlignmentLevel() int {
	matches := regexAlignmentLevel.FindStringSubmatch(q.StartCriterion)
	if len(matches) > 1 {
		lvl, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0
		}
		return lvl
	}
	return 0
}
