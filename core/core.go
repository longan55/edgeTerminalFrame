package core

//1.连接数据服务器，mqtt broker

//1. 核心启动，查询数据库，获取设备信息
//2. 判断设备状态，是否要连接
//3. 连接设备
//4. 启动查询任务

var EdgeCore edge

type edge struct {
}

func (edge *edge) Preload() error {
	////### 1.加载全部插件
	////获取插件存放路径
	//pluginpath := path.Join(viper.GetString(global.CORE_WORKPATH), viper.GetString(global.CORE_PLUGINPATH))
	////查询插件文件
	//pfs := os.DirFS(pluginpath)
	// matches, err := fs.Glob(pfs, "*.plugin")
	// if err != nil {
	// 	return err
	// }
	////加载插件
	// for _, match := range matches {
	// 	if err0 := LoadPlugin(path.Join(pluginpath, match)); err != nil {
	// 		err = errors.Join(err0)
	// 	}
	// }

	////查询数据库，加载上次运行的数据
	//### 2.加载设备信息
	//LoadDevice()

	return nil
}

//1. 读取插件信息，有哪些插件？
//2. 读取设备信息，设备需要哪些插件？
//3. 加载插件
//4. 开启任务。
