package main

type AppState struct{ Bag }

const scratchpadKey = "scratchpad"
const draftnoteKey = "draftnote"
const taglineKey = "tagline"
const recentSearchKey = "recentSearchTerms"
const appNameKey = "appName"

func (b AppState) ScratchPad() string {
	return b.GetOr(scratchpadKey, "")
}

func (b AppState) DraftNote() string {
	return b.GetOr(draftnoteKey, "")
}

func (b AppState) Tagline() string {
	return b.GetOr(taglineKey, "")
}

func (b AppState) RecentSearchTerms() string {
	return b.GetOr(recentSearchKey, "")
}

func (b AppState) AppName() string {
	return b.GetOr(appNameKey, "Gills")
}
