
new Vue({
    el: "#app",
    data: {
        rows: []
    },
    created() {
        this.getData()
    },
    methods: {
        getData: async function () {
            const headers = new Headers({
                'Accept': '*/*',
                'Content-Type': 'application/json'
            })

            const args = {
                Method: 'get_all',
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
                    method: "TaskService.GetAll",
                    params: [args]
                })
            })

            if (res.ok) {
                const data = await res.json()
                this.rows = data.result.Data
                console.log(this.rows)
            }
            else {
                const err = await res.text()
                console.log(err)
            }
        }
    }
})