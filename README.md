# go相关的一个工具包

### loadSetting
 定时从数据库中读取数据，如果有数据更改，将更改的数据传到chan中去。
 
> 使用时首次初始化需要注意的内容：
>1. 接收消息需要新开一个routine。
>2. 首次初始化初始一个空值，然后运行一次reload，以后都在routine中运行即可。

### dealjson
 处理json文件，添加json scheme的默认值


### hashlist
 一个由map和list组成的hash list。
 目的：map便于查找数据，list便于按照时序排列。