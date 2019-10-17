import urllib.request
import json

def rpc_call (url, method, args) :
    data = json.dumps({
        'id': 1,
        'method': method,
        'params': [args]
    }).encode()

    req = urllib.request.Request(url, data, { 'Content-Type': 'application/json' })

    f = urllib.request.urlopen(req)
    response = f.read()
    return json.loads(response)

url = "http://localhost:4001/rpc"
args = {'Method': 'server_stop', 'Dataset': "electricity", 'Collection': 'cable'}
jsondata = rpc_call(url, "TaskService.Show", args)

print(jsondata)
