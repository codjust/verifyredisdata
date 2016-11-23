# VerifyRedisData

language：golang<br>
author：Spotless<br>

主要使用go完成了一下工作：<br>
（1）使用os/exec包下的Command对象调用两个exe可执行程序，并获取其输出信息<br>
（2）exe程序完成了对redis的数据插入，数据结构为hash和string<br>
（3）分别取出redis的所有数据，并写入到txt文件<br>
（4）实现对哈希key的排序<br>
（5）格式化文件的数据<br>
最后使用beyond compare 工具对比两个文件内容，即可得出两个文件的差异。
