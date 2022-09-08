package anime_control

import (
	"regexp"
	"strings"
)

type BangumiFilter struct {
	multiEpisode     *regexp.Regexp
	singleEpisode    *regexp.Regexp
	defaultDelimiter string
	coarseDelimiter  string
	seasonRegexp     *regexp.Regexp
	movieRegexp      *regexp.Regexp
	teamRegexp       *regexp.Regexp
	resolutionRegexp *regexp.Regexp
	mediaInfoRegexp  *regexp.Regexp
}

func NewBangumiFilter() *BangumiFilter {
	bf := &BangumiFilter{
		multiEpisode:     regexp.MustCompile(`[ 【\[第]([0-9]{1,2}-[0-9]{1,2}) ?(?i)(END|Fin|合集)?[】\]+ 话話集]`),
		singleEpisode:    regexp.MustCompile(`[ 【\[第]([0-9]{1,4}([Vv][2-9])?)[】\[\]+ 话話集]`),
		defaultDelimiter: " []&/+【】()（）",
		coarseDelimiter:  "[]/()【】",
		seasonRegexp:     regexp.MustCompile(`([sS](0|)[0-9]+)|第.季|第.期`),
		// strict matching, also match the start and end of a substr
		movieRegexp: regexp.MustCompile(`^(?i)(剧场版|OVA|SP|OAD|Movie)$`),
		// currently, |字幕社|工作室 are not included in teams
		teamRegexp:       regexp.MustCompile(`喵萌|LoliHouse|字幕组|悠哈璃羽字幕社`),
		resolutionRegexp: regexp.MustCompile("[0-9]{3,}[pPiI]|[24][kK]|[0-9]{3,4}[xX][0-9]{3,4}"),
		mediaInfoRegexp:  regexp.MustCompile("(?i)(AVC|HEVC|AAC|WebRip|TVrip|MP4|MKV|WEB-DL|BDRip|[0-9]+-?bit|Ma10|Hi10|FLAC|BDMV|M2TS|x264|x265)"),
	}
	return bf
}

// CoarseSplit is used to find the team and resolution info
func (bf *BangumiFilter) CoarseSplit(title string) []string {
	return SplitByDelimiter(title, bf.coarseDelimiter)
}

// GetMultiEpisode return true if the episode is finished
func (bf *BangumiFilter) GetMultiEpisode(episode string) (string, bool) {
	str := bf.multiEpisode.FindStringSubmatch(episode)
	if str == nil {
		return "", false
	}
	// check is season finished
	finStr := str[2]
	if len(finStr) != 0 {
		return str[1], true
	}
	// Additional check if it is a complete set (01-12, 01-13, 01-24)
	finLikeMap := map[string]interface{}{"01-12": nil, "01-13": nil, "01-24": nil}
	if _, ok := finLikeMap[str[1]]; ok {
		return str[1], true
	}
	return str[1], false
}

func (bf *BangumiFilter) GetSingleEpisode(episode string) string {
	strList := bf.singleEpisode.FindAllStringSubmatch(episode, -1)
	if strList == nil {
		return ""
	}
	// in a reverse order
	for i := len(strList) - 1; i >= 0; i-- {
		str := strList[i]
		if len(str) > 0 && str[1] != "" {
			return str[1]
		}
	}
	return strList[0][1]
}

func (bf *BangumiFilter) GetSeasonType(title string) []string {
	return GetResultWithSplit(title, bf.defaultDelimiter, bf.seasonRegexp)
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
	for k := range resMap {
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
	if start < len(name) && !resetPos {
		subStr := name[start:]
		if len(subStr) != 0 {
			split = append(split, subStr)
		}
	}
	return split
}
