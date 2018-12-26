ps -ef|grep Zi_ziRuiLi|grep -v grep|awk '{print $2}'|xargs kill -9
ps -ef|grep Zi_zzRuiLi|grep -v grep|awk '{print $2}'|xargs kill -9
sleep 3

