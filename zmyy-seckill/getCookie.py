import configparser
import os
import threading
import yaml
config = configparser.ConfigParser()


def getCookie():
    headers = {}
    f = open('headers.txt', 'r', encoding='utf-8')
    for i in f.readlines():
        if i.rstrip():
            k, v = i.rstrip().split(': ', 2)
            if k == 'Content-Length':
                continue
            headers[k.lower()] = v
    cookie=headers['cookie']
    print(cookie)



def open_proxy(port):
    os.system('chcp 65001')
    os.system('open_proxy.bat {}'.format(port))


def capture(port='8080'):
    open_proxy(port)
    # 开启监听
    # print('mitmdump -q -s {} -p  {}'.format('capture.py', port))
    os.system('mitmdump -q -s {} -p  {}'.format('capture.py', port))


if __name__ == '__main__':
    t1 = threading.Thread(target=capture)
    t1.start()
    t1.join()
    getCookie()
