import configparser
import os
import threading

import capture





def open_proxy(port):
    os.system('chcp 65001')
    os.system('open_proxy.bat {}'.format(port))

def capture(port='8080'):
    open_proxy(port)
    # 开启监听
    # print('mitmdump -q -s {} -p  {}'.format('capture.py', port))
    os.system('mitmdump -q -s {} -p  {}'.format('capture.py', port))
if __name__ == '__main__':

#     # 抓包获取cookie
    t1 = threading.Thread(target=capture)
    t1.start()
    t1.join()
