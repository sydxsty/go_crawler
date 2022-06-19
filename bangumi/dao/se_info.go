package dao

import bgm "crawler/bangumi/anime_control"

type SEInfo struct {
	Finished      bool
	SingleEpisode string
	MultiEpisode  string
	Seasons       []string
	Movie         []string
}

func NewSEInfoFromTitle(title string, bgmFilter *bgm.BangumiFilter) *SEInfo {
	s := &SEInfo{}
	// if contains movie, append it
	s.SingleEpisode = bgmFilter.GetSingleEpisode(title)
	s.MultiEpisode, s.Finished = bgmFilter.GetMultiEpisode(title)
	s.Seasons = bgmFilter.GetSeasonType(title)
	s.Movie = bgmFilter.GetMovieType(title)
	return s
}

// GetEpisodeStringList return an ordered string list
func (s *SEInfo) GetEpisodeStringList() []string {
	var sl []string
	// multi episodes can override single episode
	if len(s.MultiEpisode) != 0 {
		sl = append(sl, s.MultiEpisode)
		if s.Finished {
			sl = append(sl, "Fin")
		}
	} else if len(s.SingleEpisode) != 0 {
		sl = append(sl, s.SingleEpisode)
	}
	addFin := false
	// no valid episode found, append season
	if len(sl) == 0 && len(s.Seasons) != 0 {
		sl = s.Seasons
		addFin = true
	}
	// if contains movie, append it
	sl = append(sl, s.Movie...)
	// only contains season, append fin
	if addFin {
		sl = append(sl, "Fin")
	}
	return sl
}
