import urllib.request
import json

def rpc_call (url, method, args) :
    data = json.dumps({
        'jsonrpc': "2.0",
        'id': 1,
        'method': method,
        'params': [args]
    }).encode()

    req = urllib.request.Request(url, data, { 'Content-Type': 'application/json' })

    f = urllib.request.urlopen(req)
    response = f.read()
    return json.loads(response)

url = "http://localhost:4001/rpc"
args = {'Method': 'get_all', 'Dataset': "electricity", 'Collection': 'substation'}
jsondata = rpc_call(url, "TaskService.GetAll", args)

print(jsondata)
