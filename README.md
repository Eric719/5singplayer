学习了Go的前面一部分，做的小练习，就是想回家打开电脑能自动放音乐。

利用5sing的接口，获取每日推荐前10首，得到mp3链接后，写入文件，通过mpg123实现自动播放。
5sing接口地址：https://github.com/i5sing/5sing-mobile-api

~~本来要取的是5sing的古筝曲，奈何各种分类搜索都不能得到我想要的曲子，不明白为什么人气那么高的曲子通过搜索和分类查找就是找不到~~


------
如果在命令行单独使用mpg123， 方法：
mpg123 url --直接播放此url
mpg123 -@music.txt 依次播放文件内url
mpg123 url url url 依次播放
