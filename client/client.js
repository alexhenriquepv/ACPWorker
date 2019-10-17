async function getdata() {

    const headers = new Headers({
        'Accept': '*/*',
        'Content-Type': 'application/json'
    })

    const args = { 
        Method: 'server_stop', 
        Dataset: "electricity", 
        Collection: 'cable' 
    }

    const res = await fetch("http://localhost:4001/rpc", {
        method: "POST",
        mode: "cors",
        headers: headers,
        body: JSON.stringify({
            jsonrpc: "2.0",
            id: 1,
            method: "TaskService.Show",
            params: [args]
        })
    })
    
    try {
        const data = await res.json()
        console.log(data)
    }
    catch (err) {
        console.log(err)
    }
}