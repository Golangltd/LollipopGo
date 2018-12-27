ps -ef|grep LollipopGo|grep -v grep|awk '{print $2}'|xargs kill -9
sleep 3

