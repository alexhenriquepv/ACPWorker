define([
    "esri/Map",
    "esri/views/MapView",
    "esri/layers/FeatureLayer",
    "esri/widgets/Expand",
    "esri/widgets/LayerList",
    "esri/widgets/BasemapGallery",
    "local/featuresManager",
    "local/sketchManager",
    "dojo/text!./config/emt.json",
], (Map, MapView, FeatureLayer, Expand, LayerList, BasemapGallery, fm, sm, config) => {

    const jsonData = JSON.parse(config)
    const appURL = jsonData.appURL
    const map = new Map(jsonData.map)
    
    const view = new MapView({
        container: "app",
        map: map,
        center: jsonData.mapView.center,
        zoom: jsonData.mapView.zoom
    })

    const layerList = new LayerList({
        view: view,
        listItemCreatedFunction: (e) => {
            e.item.actionsSections = [
                [
                    {
                        id: "load",
                        title: "Carregar",
                        className: "esri-icon-refresh"
                    }
                ]
            ]
        }
    })

    const basemapGallery = new BasemapGallery({ view: view })

    const layerListExpand = new Expand({ view: view, content: layerList })
    const basemapExpand = new Expand({ view: view, content: basemapGallery })

    layerList.on("trigger-action", e => {
        console.log(e.action.id)
        getData(e.item.layer.id)
    })
    
    map.addMany(jsonData.layers.map(l => new FeatureLayer(l)))
    view.ui.add("btn-select-by-polygon", "top-left")
    view.ui.add(layerListExpand, { position: "top-left" })
    view.ui.add(basemapExpand, { position: "top-left" })
    sm.init(view)

    function getData (collection_name) {
        Pace.track(() => {
            $.ajax({
                method: 'post',
                url: jsonData.appURL,
                contentType: "application/json",
                dataType: "json",
                data: JSON.stringify({
                    Method: 'get_all',
                    Dataset: jsonData.dataset,
                    Collection: `${collection_name}`
                }),
                success: async (data) => {
                    let features = fm.preparedFeatures(data)
                    layer = map.findLayerById(collection_name)
                    if (layer) {
                        try {
                            await layer.applyEdits({ addFeatures: features })
                            // toastr.success("Features adicionadas", promise.addFeatureResults.length)
                        }
                        catch (err) {
                            console.log(err)
                        }
                    }
                },
                error: (error) => {
                    console.log(error)
                }
            })
        })
    }

    return {
        getURL: () => { return appURL },
        getMap: () => { return map },
        getView: () => { return view },
        getData: (collection_name) => { return getData(collection_name) }
    }
})