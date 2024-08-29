package setting

import (
	"fmt"
	"os"

	// 	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 调用: setting.Conf.Spide.Cron
var Conf = new(AppConfig)

type AppConfig struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
	Open bool   `mapstructure:"open"` // 是否启动打开应用

	*Spide `mapstructure:"spide"`
	Apps   []*App `mapstructure:"apps"`
}

type Spide struct {
	DetailUrl string `mapstructure:"detail_url"`
	ListUrl   string `mapstructure:"list_url"`
	Cron      string `mapstructure:"cron"`
}

type App struct {
	Name string `mapstructure:"name" json:"name"`
	Path string `mapstructure:"path" json:"path"`
}

func Init(filePath string) {
	// 方式1：直接指定配置文件路径（相对路径或者绝对路径）
	// 相对路径：相对执行的可执行文件的相对路径
	//viper.SetConfigFile("./conf/config.yaml")
	// 绝对路径：系统中实际的文件路径
	//viper.SetConfigFile("/Users/liwenzhou/Desktop/bluebell/conf/config.yaml")

	// 方式2：指定配置文件名和配置文件的位置，viper自行查找可用的配置文件
	// 配置文件名不需要带后缀
	// 配置文件位置可配置多个
	//viper.SetConfigName("config") // 指定配置文件名（不带后缀）
	//viper.AddConfigPath(".") // 指定查找配置文件的路径（这里使用相对路径）
	//viper.AddConfigPath("./conf")      // 指定查找配置文件的路径（这里使用相对路径）

	// 基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	//viper.SetConfigType("json")

	viper.SetConfigFile(filePath)

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		// panic(err)

		// 创建默认配置并写入文件
		defaultConfig := `name: "远程控制"
port: 666
open: false
apps:
  - name: 微信
    path: E:\Program Files (x86)\Tencent\WeChat\WeChat.exe
  - name: 网易云
    path: E:\Program Files (x86)\NetEase\CloudMusic\cloudmusic.exe
`
		err = os.WriteFile(filePath, []byte(defaultConfig), 0644)
		if err != nil {
			fmt.Printf("Failed to create config file, err: %v\n", err)
			panic(err)
		}

		fmt.Printf("Created default config file at %s\n", filePath)
		// 重新读取配置文件
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Failed to read the newly created config file, err: %v\n", err)
			panic(err)
		}
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		panic(err)
	}

	// viper.WatchConfig()
	// viper.OnConfigChange(func(in fsnotify.Event) {
	// 	fmt.Println("配置文件修改了...")
	// 	if err := viper.Unmarshal(Conf); err != nil {
	// 		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	// 		panic(err)
	// 	}
	// })
	return
}
