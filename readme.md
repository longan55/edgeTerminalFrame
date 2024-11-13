smt 2.0 设计

1. 对象管理容器
	1.1 管理连接
	1.2 管理状态 连接状态，启停状态
	1.3 管理重新连接
2. 对象数据热更新
3. 重新设计断点续传
4. 添加数据状态
5. 重新设计数据源之类的

软件功能:

### 第一阶段: 终端级
1. 网关自身相关
   1.1 自身信息, 状态.
   1.2 系统级操作, ip查看、设置ip.
   1.3 和平台通信、同步信息
2. 外界设备相关
   2.1 自身信息、状态信息、实时同步更新。 平台修改实时更新。
   2.1 物联设备协议对接 -》 指令
   2.3 通过cifs/smb、sftp、等获取外部电脑的数据。
   2.4 自定义功能
3. 平台相关

### 第二阶段: 平台级
概述: 边缘端设备框架, 可以对接平台. 用MQTT通信.
设备状态存在服务器内存中, 不存在数据库. 