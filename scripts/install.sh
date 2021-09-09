###
 # @Description: 
 # @Author: huyongchao
 # @Date: 2021-07-19 09:14:25
 # @LastEditTime: 2021-07-20 10:19:00
 # @LastEditors: huyongchao
 # @Usage: 
### 


cd ..
git pull
export GOPROXY=http://mirrors.aliyun.com/goproxy/
go install .
cd /data/apps/go/bin
killall -9 ant_push
nohup ant_push &
