# Android QQ Pic Cleaner
该工具用于清理安卓QQ上指定好友或群聊的聊天图片文件，只适用于非NT版本（旧版），且需要root权限。

建议运行之前做好备份工作，以免出现意外情况导致数据丢失。

## 使用方法
将该程序放到安卓设备中，通过终端以root权限运行。

## 运行参数
* -h, --help: 显示说明
* -q: 用户的QQ号，用于查找数据库文件
* -f, --friends: 指定的好友QQ号，如有多个使用冒号 **:** 分隔
* -g, --groups: 指定的群号，如有多个使用冒号 **:** 分隔
* -l, --list: 列出用户所有的好友和群聊
* --dry-run: 模拟运行，但不进行删除

## 示例

列出用户1234的所有好友和群聊：

`android-qq-pic-cleaner -q 1234 -l`

删除用户1234的好友2345、3456及群聊4567、5678、6789的聊天图片：

`android-qq-pic-cleaner -q 1234 -f 2345:3456 -g 4567:5678:6789`

模拟删除用户1234的好友2345的聊天图片（但实际并没有删除）：

`android-qq-pic-cleaner -q 1234 -f 2345 --dry-run`

## 参考
1. <https://github.com/QQBackup/QQ-History-Backup>