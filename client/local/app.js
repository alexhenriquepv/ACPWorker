define([
    "esri/Map",
    "esri/views/MapView",
    "esri/layers/GraphicsLayer",
    "esri/widgets/LayerList",
    "dojo/text!./config/emt.json"
], (Map, MapView, GraphicsLayer, LayerList, config) => {

    jsonData = JSON.parse(config)

    const appURL = jsonData.appURL
    const map = new Map(jsonData.map)

    const view = new MapView({
        container: "app",
        map: map,
        center: jsonData.mapView.center,
        zoom: jsonData.mapView.zoom
    })

    const layerList = new LayerList({ view: view })
    
    map.addMany(jsonData.features.map(f => new GraphicsLayer(f)))
    view.ui.add(layerList, { position: "top-right" })

    return {
        getURL: () => { return appURL },
        getMap: () => { return map },
        getView: () => { return view }
    }
})