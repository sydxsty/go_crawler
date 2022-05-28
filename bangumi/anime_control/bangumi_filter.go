package anime_control

import (
	"regexp"
	"strings"
)

type BangumiFilter struct {
	multiEpisode     *regexp.Regexp
	singleEpisode    *regexp.Regexp
	finEpisode       *regexp.Regexp
	defaultDelimiter string
	coarseDelimiter  string
	movieRegexp      *regexp.Regexp
	teamRegexp       *regexp.Regexp
	resolutionRegexp *regexp.Regexp
	mediaInfoRegexp  *regexp.Regexp
}

func NewBangumiFilter() *BangumiFilter {
	bf := &BangumiFilter{
		multiEpisode:     regexp.MustCompile(`[ 【\[]([0-9]{1,2}-[0-9]{1,2})[】\] ]`),
		singleEpisode:    regexp.MustCompile(`[ 【\[第]([0-9]{1,3})[】\] 话話]`),
		finEpisode:       regexp.MustCompile(`[ 【\[](?i)( ?fin)[】\] ]`),
		defaultDelimiter: " []&/【】()（）",
		coarseDelimiter:  "[]/()【】",
		movieRegexp:      regexp.MustCompile(`剧场版|OVA|OAD|(?i)Movie|([sS](0|)[0-9]+)|第.季`),
		// currently, |字幕社|工作室 are not included in teams
		teamRegexp:       regexp.MustCompile(`喵萌|LoliHouse|字幕组`),
		resolutionRegexp: regexp.MustCompile("[0-9]{3,}[pPiI]|[24][kK]|[0-9]{3,4}x[0-9]{3,4}"),
		mediaInfoRegexp:  regexp.MustCompile("(?i)(AVC|HEVC|AAC|WebRip|TVrip|MP4|MKV|WEB-DL|BDRip|[0-9]+-?bit|Ma10|Hi10|FLAC|BDMV|M2TS|x264|x265)"),
	}
	return bf
}

// CoarseSplit is used to find the team and resolution info
func (bf *BangumiFilter) CoarseSplit(title string) []string {
	return SplitByDelimiter(title, bf.coarseDelimiter)
}

func (bf *BangumiFilter) GetMultiEpisode(episode string) string {
	str := bf.multiEpisode.FindStringSubmatch(episode)
	if str == nil {
		return ""
	}
	// check is season finished
	finStr := bf.finEpisode.FindAllString(episode, -1)
	if len(finStr) != 1 {
		return str[1]
	}
	return str[1] + " Fin"
}

func (bf *BangumiFilter) GetSingleEpisode(episode string) string {
	strList := bf.singleEpisode.FindAllStringSubmatch(episode, -1)
	if strList == nil {
		return ""
	}
	for _, str := range strList {
		if len(str) > 0 && str[1][0] == '0' {
			return str[1]
		}
	}
	return strList[0][1]
}

func (bf *BangumiFilter) GetMovieType(title string) []string {
	return GetResultWithSplit(title, bf.defaultDelimiter, bf.movieRegexp)
}

func (bf *BangumiFilter) GetResolution(title string) []string {
	return GetResultWithSplit(title, bf.defaultDelimiter, bf.resolutionRegexp)
}

func (bf *BangumiFilter) GetMediaInfo(title string) []string {
	return GetResultWithSplit(title, bf.defaultDelimiter, bf.mediaInfoRegexp)
}

func GetResultWithSplit(title, delim string, regexp *regexp.Regexp) []string {
	strList := SplitByDelimiter(title, delim)
	resMap := make(map[string]bool)
	for _, str := range strList {
		for _, res := range regexp.FindAllString(str, -1) {
			resMap[res] = true
		}
	}
	var resList []string
	for k, _ := range resMap {
		resList = append(resList, k)
	}
	return resList
}

// GetTeam the team name usually in piece 1
func (bf *BangumiFilter) GetTeam(title string) []string {
	coarse := SplitByDelimiter(title, bf.coarseDelimiter)
	// nil or not cut yet
	if len(coarse) <= 1 {
		return nil
	}
	teams := SplitByDelimiter(coarse[0], "&")
	valid := false
	for _, team := range teams {
		if len(bf.teamRegexp.FindString(team)) != 0 {
			valid = true
			break
		}
	}
	if valid {
		return teams
	}
	return nil
}

func SplitByDelimiter(name, delimiter string) []string {
	var split []string
	start := 0
	resetPos := true
	for k, v := range name {
		if resetPos {
			start = k
			resetPos = false
		}
		if strings.Contains(delimiter, string(v)) {
			subStr := name[start:k]
			resetPos = true
			if len(subStr) != 0 {
				split = append(split, subStr)
			}
		}
	}
	if start < len(name) {
		subStr := name[start:]
		if len(subStr) != 0 {
			split = append(split, subStr)
		}
	}
	return split
}