# 本项目目前包含两部分系统
1. neubt爬虫（自动RSS指定板块，根据过滤条件下载种子）
2. 自动发布bangumi.moe种子到neubt（目前不稳定，仍在开发）

# NEUBT爬虫
编译入口地址
https://github.com/sydxsty/go_crawler/blob/master/neubt/rss/main.go
## 使用方法
1. 编译（或者直接下载release界面的）
2. 在运行目录上新建data文件夹
3. 在data文件夹内新建download文件夹（必须）
4. 在data文件夹内新建config.yaml里面加上如下内容：
```
username: 你的用户名
password: 你的密码
use_cookie: false    #第二次就可以填true了
cookie_path: cookie
qb_addr: "http://127.0.0.1:8080/"
qb_username: 你的qbittorrentWEBUI的用户名
qb_password: 你的qbittorrentWEBUI的密码
torrent_path: "./data/download/"
leveldb_path: "./data/db/"
discount_water_mark: 100     # 100为只下载免费种子，50为下载50折扣的种子，0为下载全部种子
thread_water_mark: 1695272    #第一次启动时候从哪个ID开始爬
```
程序设置为600秒爬取一次，已经爬取的界面不会再次爬取（会保存在leveldb中），资源索引更新的本来就很慢，不建议改

# 自动发种
编译入口地址
https://github.com/sydxsty/go_crawler/blob/master/neubt/bgm_auto_poster/main.go
## 工作流程
1. 爬取bangumi上更新的种子
2. 查询ptgen上是否有对应动漫信息
3. 如果有，通过qbittorrent先下载种子
4. 当下载完成时，使用mediainfo获取视频参数
5. 在neubt上发布种子
6. 重新下载发布过的种子，通过qbittorrent做种

## 使用方法
目录结构同NEUBT爬虫，运行即可
### 强烈不建议直接用，目前过滤结构仍不完善，会发布大量重复种子
