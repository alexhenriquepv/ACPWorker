
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
                console.log(data.result.Message, this.rows)
                this.initMap(this.rows)
            }
            else {
                const err = await res.text()
                console.log(err)
            }
        },
        initMap (data) {
            const streets = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                id: 'mapbox.streets',
                attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            })

            const map = L.map('app', {
                center: [52.221841, 0.109577],
                zoom: 13,
                layers: [streets]
            })

            const cable = L.geoJSON(data, {
                onEachFeature: (feature, layer) => {
                    if (feature.geometry.type == "LineString") {
                        const latLngs = feature.geometry.coordinates.map(function (latLng) { return L.latLng(latLng[1], latLng[0]) })
                        geom = L.polyline(latLngs, { color: 'red' })
                        geom.addTo(map)
                    }
                }
            })

            cable.bindPopup(layer => {
                const prop = layer.feature.properties
                return `Cable id: ${prop.id}; length: ${prop.length}; status: ${prop.status}; voltage: ${prop.voltage}`
            })

            cable.addTo(map)
        }
    }
})