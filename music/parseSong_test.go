package music

import (
	"reflect"
	"testing"
)

func TestParseJson(t *testing.T) {
	youtubeURL := "https://www.youtube.com/watch?v=QXV8iBwvZH4&ab_channel=Kozko"
	expected := Song{
		URL:       "https://rr1---sn-f5f7lnl7.googlevideo.com/videoplayback?expire=1708920254&ei=XrnbZaK3FM7F6dsPyYCykA0&ip=194.29.137.25&id=o-ACOQL7uTYhnEaAIBwoMkl5vszzgXqE1FoaZKREdl8kzH&itag=251&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&mh=oH&mm=31%2C26&mn=sn-f5f7lnl7%2Csn-4g5e6nss&ms=au%2Conr&mv=m&mvi=1&pl=19&initcwndbps=1456250&spc=UWF9f8CeR7Wcg2vW9pJZszO8EG9p9nbXsYo8-BJIcW_MsAw&vprv=1&svpuc=1&mime=audio%2Fwebm&gir=yes&clen=2406114&dur=142.061&lmt=1671315105425811&mt=1708898211&fvip=3&keepalive=yes&fexp=24007246&c=ANDROID&txp=5318224&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cspc%2Cvprv%2Csvpuc%2Cmime%2Cgir%2Cclen%2Cdur%2Clmt&sig=AJfQdSswRQIhAJvJTbDXvKdA7zU1zcfnj475JN8JFPUVUTtugQPGxvtqAiB0muCezMBx_D64KIUmQ9IQjSfAJqqCEOatk9DLxWC6Kg%3D%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=APTiJQcwRQIhAKYTGn8NWpH7bU3OjxRoAuanq8k1pAznsyPzJ_Yyq5nCAiAE-kn0By0pzFrj_mMySwX5UNgkAMCQHkRHIvBlEqKiig%3D%3D",
		Title:     "Palion - Zielone 2 Nightcore",
		Thumbnail: "https://i.ytimg.com/vi/QXV8iBwvZH4/maxresdefault.jpg",
		Duration:  "2:22",
	}

	output, _ := GetSongInfo(youtubeURL)
	result := NewSong(output)

	reflect.DeepEqual(expected, result)
}
