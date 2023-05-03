# 介绍
使用Golang生成getCPUID()和getDiskID()以及getMACAddr（）； 用于机器码生成。

# 说明
* 在 macOS 上，我们使用 system_profiler 命令来获取硬盘序列号，并使用 grep 和 awk 命令从输出中提取序列号。
* 在 Linux 上，我们仍然使用 lsblk 命令来获取硬盘序列号。


