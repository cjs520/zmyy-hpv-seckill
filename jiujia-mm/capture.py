import json
import os

from mitmproxy import ctx
from utils import *

flag = False
isBind = False


def request(flow):
    global flag
    global isBind
    url = flow.request.url
    print(url)
    if url == 'https://wx.scmttec.com/base/region/childRegions.do?parentCode=' and isBind:
        log('保存headers')
#         print(flow.request.headers)
        with open('header.txt', 'w', encoding='utf-8') as f:
            for k, v in dict(flow.request.headers).items():
                f.write('{}: {}'.format(k, v))
                f.write('\n')
        flag = True
def response(flow):
    global flag
    global isBind
    if 'https://wx.scmttec.com/member/wx/memberCenter/getMemberInfo.do' in flow.request.url:
        res = json.loads(str(flow.response.text))
        if res['code'] == "1001":
            error('非法请求,请登录授权后访问!')
        else:
            isBind = True
    if flag:
        os.system('close_proxy.bat')
        ctx.master.shutdown()

