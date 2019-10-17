var client = new RestClient("http://localhost:4001/rpc");
var request = new RestRequest(Method.POST);
request.AddHeader("cache-control", "no-cache");
request.AddHeader("Accept-Encoding", "gzip, deflate");
request.AddHeader("Cache-Control", "no-cache");
request.AddHeader("Accept", "*/*");
request.AddHeader("Content-Type", "application/json");
request.AddParameter("undefined", "{\r\n\"jsonrpc\": \"2.0\",\r\n\"id\": 1,\r\n\"method\": \"TaskService.Show\",\r\n\"params\": []\r\n}", ParameterType.RequestBody);
IRestResponse response = client.Execute(request);