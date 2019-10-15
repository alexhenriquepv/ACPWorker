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

url = "http://10.84.125.24:4001/rpc"
args = {'Dataset': "electric", 'Collection': 'eo_cable'}
jsondata = rpc_call(url, "TaskService.Show", args)

print(jsondata)