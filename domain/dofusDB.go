package domain

type DofusDbSearchQuest struct {
	Total int                 `json:"total"`
	Limit int                 `json:"limit"`
	Skip  int                 `json:"skip"`
	Data  []DofusDbQuestLight `json:"data"`
}

type DofusDbQuestLight struct {
	ID   int               `json:"id"`
	Name map[string]string `json:"name"`
}
