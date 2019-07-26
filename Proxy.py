# -*- coding: utf-8 -*-
# Create by a2si

import pycurl

from Core import Config
from Core import Logs


class WebProxy(object):

    __slots__ = [
        'ProxyEnable',
        'Handle',
        '__ProxyIP', '__ProxyPort', '__ProxyUser', '__ProxyPwd'
    ]

    def __init__(self, Handle=None):
        self.Handle = Handle
        self.__ProxyIP = ""
        self.__ProxyPort = 0
        self.__ProxyUser = ""
        self.__ProxyPwd = ""
        Config.SetDefaults("MiniWeb", {
            "UsedProxy": False,
            "ProxyType": "",
            "ProxyIP": "",
            "Port": "",
            "ProxyUser": "",
            "ProxyPwd": ""
        })
        self.ProxyEnable = Config.Get("MiniWeb", "UsedProxy")

    def __del__(self):
        pass

    # 基础函数
    def SetOption(self, *args):
        if self.Handle is not None:
            self.Handle.setopt(*args)

    def GetOption(self, *args):
        if self.Handle is not None:
            return self.Handle.getinfo(*args)
        return None

    def UsedProxy(self):
        return self.ProxyEnable

    def SetProxyEnable(self, Enable):
        self.ProxyEnable = Enable
        Config.Set("MiniWeb", "UsedProxy", Enable)

    def CancelProxy(self):
        self.ProxyEnable = False
        Config.Set("MiniWeb", "UsedProxy", False)
        self.SetOption(pycurl.PROXY, "")
        self.SetOption(pycurl.PROXYUSERPWD, "")

    def SetProxyType(self, ProxyType):
        """ProxyType
            -- "HTTP"          PROXYTYPE_HTTP
            -- "HTTP10"        PROXYTYPE_HTTP_1_0
            -- "SOCKS4"        PROXYTYPE_SOCKS4
            -- "SOCKS4A"       PROXYTYPE_SOCKS4A
            -- "SOCKS5"        PROXYTYPE_SOCKS5
            -- "SOCKS5_HOSTNAME"      PROXYTYPE_SOCKS5_HOSTNAME
        """
        Config.Set("MiniWeb", "ProxyType", ProxyType)
        if ProxyType == "HTTP":
            self.SetOption(pycurl.PROXYTYPE, pycurl.PROXYTYPE_HTTP)
        elif ProxyType == "HTTP10":
            self.SetOption(pycurl.PROXYTYPE, pycurl.PROXYTYPE_HTTP_1_0)
        elif ProxyType == "SOCKS4":
            self.SetOption(pycurl.PROXYTYPE, pycurl.PROXYTYPE_SOCKS4)
        elif ProxyType == "SOCKS4A":
            self.SetOption(pycurl.PROXYTYPE, pycurl.PROXYTYPE_SOCKS4A)
        elif ProxyType == "SOCKS5":
            self.SetOption(pycurl.PROXYTYPE, pycurl.PROXYTYPE_SOCKS5)
        elif ProxyType == "SOCKS5_HOSTNAME":
            self.SetOption(pycurl.PROXYTYPE, pycurl.PROXYTYPE_SOCKS5_HOSTNAME)
        pass

    def SetProxyIP(self, IP, Port):
        Config.Set("MiniWeb", "ProxyIP", IP)
        Config.Set("MiniWeb", "Port", Port)
        if IP == "" and Port == "":
            return
        self.SetOption(pycurl.PROXY, "%s:%s" % (IP, Port))

    def SetProxyUserPwd(self, User, Pwd):
        Config.Set("MiniWeb", "ProxyUser", User)
        Config.Set("MiniWeb", "ProxyPwd", Pwd)
        if User == "" and Pwd == "":
            return
        self.SetOption(pycurl.PROXYUSERPWD, "%s:%s" % (User, Pwd))

    def CheckProxyError(self, ErrCode, ErrMsg):
        if self.ProxyEnable:
            '''
                56
                    # Received HTTP code 405 from proxy after CONNECT
                    # Received HTTP code 403 from proxy after CONNECT
                52
                    # Empty reply from server
                7
                    # 连接失败
                    # Connection refused
                28
                    # 连接超时
                35
                    # SSL 错误
            '''
            ProxyError = [56, 52, 7, 28, 35]
            if ErrCode in ProxyError:
                self.Proxy.CancelProxy()
                return True
            print("ProxyError[%s]: %s" % (ErrCode, ErrMsg))
            return True
        return False

    def gwcInitProxy(self):
        # 代理模式支持
        if self.ProxyEnable:
            ProxyType = Config.Get("MiniWeb", "ProxyType")
            self.SetProxyType(ProxyType)
            IP = Config.Get("MiniWeb", "ProxyIP")
            Port = int(Config.Get("MiniWeb", "Port"))
            self.SetProxyIP(IP, Port)
            ProxyUser = Config.Get("MiniWeb", "ProxyUser")
            ProxyPwd = Config.Get("MiniWeb", "ProxyPwd")
            self.SetProxyUserPwd(ProxyUser, ProxyPwd)
        else:
            self.CancelProxy()

    def gwcProxyCheck(self):
        if self.ProxyEnable:
            ErrorCode = self.GetOption(pycurl.OS_ERRNO)
            StatusCode = self.GetOption(pycurl.HTTP_CODE)
            self.__StatusCode = StatusCode
            self.__ErrorCode = ErrorCode
            if StatusCode == 596:
                ErrorCode = 28
            if ErrorCode == 28:
                print("代理无效, 重新访问...")
                Logs.Info("MiniWeb", "Proxy.Error.Refresh.Proxy")
        #         {
        #             global $Console;
        #             $Console->PrintError("代理无效, 重新访问...");
        #             $this->ProxySDK->SwitchProxy();
        #             $ProxyInfo = $this->ProxySDK->GetProxyInfo();
        #             $this->SetProxy( $ProxyInfo['Host'], $ProxyInfo['Port'], $ProxyInfo['User'], $ProxyInfo['Pwd'] );
        #         }
            if ErrorCode == 28:
                print("代理无效28, 重新访问...")
                return False
        return True
