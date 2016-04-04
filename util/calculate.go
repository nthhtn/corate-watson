package util

import (
	"time"
	"math"
	"strconv"
)

func TimeShow(createdAt float64) string{
	now:=time.Now().UTC().UnixNano()/1e6
	r:=float64((now-int64(createdAt))/1000)
	var timeShow string
	if (r<60){
		timeShow=strconv.Itoa(int(math.Floor(r)))+" s"
	} else{
		if (r/86400<1){
			if (r/3600<1){
				if (r/60<2){
					timeShow=strconv.Itoa(int(math.Floor(r/60)))+" min"
				} else{
					timeShow=strconv.Itoa(int(math.Floor(r/60)))+" mins"
				}
			} else{
				if (r/3600<2){
					timeShow=strconv.Itoa(int(math.Floor(r/3600)))+" hour"
				} else{
					timeShow=strconv.Itoa(int(math.Floor(r/3600)))+" hours"
				}
			}
		} else{
			if (r/86400<2){
				timeShow=strconv.Itoa(int(math.Floor(r/86400)))+" day"
			} else{
				timeShow=strconv.Itoa(int(math.Floor(r/86400)))+" days"
			}
		}
	}
	return timeShow
}