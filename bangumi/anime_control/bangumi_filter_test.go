package anime_control

import (
	"log"
	"testing"
)

func TestSplit(t *testing.T) {
	unit := map[string][]string{
		"[梦蓝字幕组&VCB_S]Crayonshinchan 蜡笔小新[1105][2021.11.06][AVC][1080P][GB_JP][MP4]V2.mp4": {
			"梦蓝字幕组",
			"VCB_S",
			"Crayonshinchan",
			"蜡笔小新",
			"1105",
			"2021",
			"11",
			"06",
			"AVC",
			"1080P",
			"GB_JP",
			"MP4",
			"V2",
			"mp4",
		},
		"The Last Son 2021.mkv": {
			"The",
			"Last",
			"Son",
			"2021",
			"mkv",
		},
		"Midway 2019 2160p CAN UHD Blu-ray HEVC DTS-HD MA 5.1-THDBST@HDSky.nfo": {
			"Midway",
			"2019",
			"2160p",
			"CAN",
			"UHD",
			"Blu",
			"ray",
			"HEVC",
			"DTS",
			"HD",
			"MA",
			"5",
			"1",
			"THDBST",
			"HDSky",
			"nfo",
		},
	}
	for k, v := range unit {
		actual := SplitByDelimiter(k, " -.@&[]")
		for k2, v2 := range v {
			if actual[k2] != v2 {
				t.Errorf("Split(%s) = %s; expected %s", k, actual[k2], v2)
			}
		}
	}
}

var testCase []string

func TestGetEpisode(t *testing.T) {
	episodeFilter := NewBangumiFilter()
	for _, in := range testCase {
		res1 := episodeFilter.GetSingleEpisode(in)
		res2 := episodeFilter.GetMultiEpisode(in)
		res3 := episodeFilter.GetSeasonType(in)
		log.Printf("single: %s, multi: %s, movie: %s", res1, res2, res3)
	}
}

func TestGetTeam(t *testing.T) {
	episodeFilter := NewBangumiFilter()
	for _, in := range testCase {
		res1 := episodeFilter.GetTeam(in)
		log.Printf("team: %s", res1)
	}
}

func TestGetResolution(t *testing.T) {
	episodeFilter := NewBangumiFilter()
	for _, in := range testCase {
		res1 := episodeFilter.GetResolution(in)
		log.Printf("resolution: %s", res1)
	}
}

func TestGetMediaInfo(t *testing.T) {
	episodeFilter := NewBangumiFilter()
	for _, in := range testCase {
		res1 := episodeFilter.GetMediaInfo(in)
		log.Printf("meida info: %s", res1)
	}
}

func TestIntegrate(t *testing.T) {
	bgmFilter := NewBangumiFilter()
	for _, in := range testCase {
		log.Printf("GetResolution: %s", getString(bgmFilter.GetResolution(in)))
		log.Printf("GetMediaInfo: %s", getString(bgmFilter.GetMediaInfo(in)))
		log.Printf("GetTeam: %s", getString(bgmFilter.GetTeam(in)))
		log.Printf("GetSingleEpisode: %s", getString([]string{getString(bgmFilter.GetSeasonType(in)), bgmFilter.GetSingleEpisode(in), bgmFilter.GetMultiEpisode(in)}))
		log.Printf("GetMovie: %s", getString(bgmFilter.GetMovieType(in)))
	}
}

func getString(strList []string) string {
	res := ""
	for _, str := range strList {
		if str == "" {
			continue
		}
		if len(res) != 0 {
			res += " "
		}
		res += str
	}
	return res
}

func init() {
	testCase = []string{
		"[NC-Raws] 杜鵑婚約 [特別篇] / Kakkou no Iinazuke (A Couple of Cuckoos) - 05 (Baha 1920x1080 AVC AAC MP4)",
		"【喵萌Production】★04月新番★[歌愈少女/Healer Girl][08][1080p][繁日雙語][招募翻譯]",
		"【喵萌奶茶屋】★04月新番★[夏日重现/Summer Time Rendering][07][1080p][简日双语][招募翻译]",
		"[ANi] 杜鵑婚約 [特別篇] - 05 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4]",
		"[jibaketa合成&音頻壓制][TVB粵語]食戟之靈 / Shokugeki no Souma - 21 [粵日雙語+內封繁體中文字幕][BD 1920x1080 x264 AACx2 SRT TVB CHT]",
		"[NC-Raws] 少年歌行 风花雪月篇 / Youths and Golden Coffin S2 - 31 (B-Global Donghua 1920x1080 HEVC AAC MKV)",
		"[轻之国度字幕组][盾之勇者成名录 SEASON2][08][720P][MP4]",
		"[霜庭云花Sub][夏日重现 / サマータイムレンダ / Summer Time Rendering][07][720P][AVC][简日内嵌][TVRip先行][招募]",
		"[Lilith-Raws] 川尻小玉的懶散生活 / Atasha Kawajiri Kodama Da yo - 17 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4]",
		"[桜都字幕組] 3秒後、野獣。~在联谊会的角落里、他是肉食系 / 3 Byou Go, Yajuu. Goukon de Sumi ni Ita Kare wa Midara na Nikushoku Deshita [07][1080p][繁體內嵌]",
		"[DHR动研字幕组&茉语星梦&VCB-Studio] DanManchi 3/ 在地下城寻求邂逅是否搞错了什么 第三季 8-bit 720p AVC BDRip [S3 MP4 CHS Ver]",
		"[BDMV][220225][Words Worth Blu-ray Archive BOX SPECIAL EDITION][JP]",
		"[LoliHouse] 攻壳机动队 SAC_2045 S2 / Ghost in the Shell SAC_2045 Season2 [WebRip 1080p HEVC-10bit AAC EAC3][简繁英日字幕][Fin]",
		"[c.c動漫][4月新番][RPG不動產][03-08][BIG5][1080P][MP4]",
		"【喵萌Production】★01月新番★[CUE! 短篇动画/CUE! Short Anime][05-06][BDRip][1080p][简日双语][招募翻译]",
		"[jibaketa合成&壓制][代理商粵語]咒術迴戰 咒胎戴天篇 / Jujutsu Kaisen 01-08 Fin [粵日雙語+內封繁體中文字幕][BD 1920x1080 x264 AACx2 SRT Ani-One CHT]",
		"[GM-Team][国漫][凡人修仙传 再别天南][Fan Ren Xiu Xian Zhuan][2022][02-05][AVC][GB][1080P]",
		"[桜都字幕组] 白金终局 / Platinum End [01-24 Fin][1080p][简体内嵌]",
		"[VCB-Studio] Kono Subarashii Sekai ni Shukufuku wo! / 为美好的世界献上祝福！ 10-bit 1080p HEVC BDRip [S1-S2 Reseed + Movie Fin]",
		"[雪飘工作室][Tropical-Rouge！Precure/トロピカル～ジュ！プリキュア][剧场版_冰雪公主与奇迹的戒指！][720p][简体内嵌](检索:光之美少女/Q娃) 附外挂字幕",
		"[LoliHouse闹钟保护协会] 少女☆歌剧 剧场版 / Shoujo☆Kageki Revue Starlight Movie / 劇場版 少女☆歌劇 [BDRip 1920x804 HEVC-10bit FLAC]",
		"[风车字幕组][名侦探柯南剧场版第24部][绯色的子弹/绯色的弹丸][1080P][简体][MP4][BDRip]",
		"【千夏字幕组】【剧场版 紫罗兰永恒花园/薇尔莉特·伊芙嘉登_Violet Evergarden the Movie】[剧场版][BDRip_Full HD_AVC][简体]",
		"【极影字幕社】★10月新番 结城友奈是勇者 大满开之章 第07话 GB 1080P MP4（字幕社招人内详）",
		"[Lilith-Raws] 不起眼女主角培育法 / Saenai Heroine no Sodatekata Flat [00-11][Baha][WEB-DL][1080p][AVC AAC][CHT][MP4]",
		"[NC-Raws] 键等 / Kaginado - 06 [B-Global][WEB-DL][1080p][AVC AAC][Multiple Subtitle][MKV]",
		"【极影字幕社】LoveLive! 虹咲学园学园偶像同好会 第2期 第08集 加料剪辑版 GB_CN HEVC_opus 1080p",
		"【千夏字幕组】【来玩游戏吧 / 游戏3人娘_Asobi Asobase】[第01-12话][BDRip_1080p_AVC][简体][修正合集]",
	}
}
