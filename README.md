# QNQ
神奇圈圈 V0.0.4.8

## 产品特性
### 文件同步
#### 本地同步（支持周期和定时同步策略以及进度显示）
- 支持批量同步
- 支持单文件同步
- 支持分区同步
#### 远程同步
- 支持远程单文件同步（仅支持小文件）
### 系统信息
- 支持磁盘分区基础信息显示
- 支持磁盘分区测速（读取速度偏大，写入速度较为准确）
### 系统设置
- 操作日志记录

## 更新说明
1. 解决了不打印日志文件的问题。
2. 解决远程同步中可能出现的panic和内存、CPU占用过大以及文件句柄不释放的问题。
3. 现在一个QNQ可以连接多个远程QNQ。
4. 对QNQ间网络通讯做了进一步规范。
5. 去除了冗余的代码以及不安全的rest接口（未来将重新引入rest接口）。
6. 这次一个月才更新版本，是因为开发文件快照失败了，低估了快照的实现难度。计划在第四季度实现该功能。

## 0.0.5开发计划
1. 预计6月23号发布。
2. 增加远程批量同步，远程同步支持同步策略。
3. 增加同步结果的GUI。
4. 尽量优化下界面吧，有点丑。


## 使用说明
### 本地同步
#### 以分区同步为例
1. 选择源路径（Source Path）。
2. 选择目标路径（Target Path）。
3. 点击开始按钮（Start按钮）。
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/part_sync.PNG)
#### 同步策略
首先需要打开全局开关（Global switch），再根据自己的需求选择周期策略（periodic sync）或者是定期策略（timing sync），两种策略可以叠加。
##### 周期策略（periodic sync）
根据选择的时间单元和时间单位进行同步。
如选择1 Hour，则每个小时同步一次。
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/sync_policy.PNG)
##### 定期策略（timing sync）
根据选择的日期及时间进行同步。
如选择了每周三、周五的15点32进行同步，则当系统日期到达这几个时间点时，开始同步。
### 远程同步
#### 连接远端QNQ
限制：两端的QNQ版本号必须大于等于0.0.4.8，这个版本规范了远程QNQ间的通讯。
1. 点击➕，新建一个远程QNQ页签。
2. 输入目标IP，点击认证（Auth）按钮。
3. 目标端QNQ点击同意按钮，方可完成验证。
4. 切换至其他页面再回到该页面，可以看到连接状态（这是将来的优化点，将数据与GUI做双向数据绑定，避免不动态刷新的问题）。
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/remote_auth_dia.PNG)
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/remote_auth_success.PNG)
5. 查看连接了本端的IP。
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/remote_qnq_list.PNG)
#### 远程单文件同步
在有可使用的远程QNQ时，选择目标QNQ的IP以及本地路径，点击开始按钮后，将会把本地文件同步至对端的相同路径下。（0.0.5版本将实现选择对端路径的功能，目前这里体验很差）
如选择了E:/test.txt，请确保对端存在E这个盘符。
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/remote_auth_success.PNG)
### 系统信息
#### 磁盘信息
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/disk_info.PNG)
#### 磁盘测速
这里的GUI布局很奇怪，因为左下方计划做一个实时性能统计图表，但是框架目前不支持。
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/disk_speed_test.PNG)

### 操作日志
![image](https://github.com/wangshenghao1/QNQ/blob/main/instructed/img/QLog.PNG)