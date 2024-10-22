package fbServices

import "runtime"

//根据操作系统判断服务器配置
func InitOSConfig() {
	//log.Println("============ 正在使用"+runtime.GOOS+"系统 ============")
	switch runtime.GOOS {
	case "windows":
		//ccpPath = windowCcpPath
		//credPath = windowCredPath
		chainBrowserConfigPath = windowConfigPath
		break
	case "linux":
		//ccpPath = linuxCcpPath
		//credPath = linuxCredPath
		chainBrowserConfigPath = linuxConfigPath
		break
	default:
		//ccpPath = linuxCcpPath
		//credPath = linuxCredPath
		chainBrowserConfigPath = linuxConfigPath
		break
	}
}
