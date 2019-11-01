define([
    "esri/layers/GraphicsLayer",
    "esri/widgets/Sketch/SketchViewModel",
    "dgrid/Grid"
], (GraphicsLayer, SketchViewModel, Grid) => {

    let layerViews = []
    let view, sketchViewModel, highlight, grid
    let polygonGraphicsLayer = new GraphicsLayer({ id: 'polygons', title: 'polygons', listMode: 'hide' })
    
    $("#btn-select-by-polygon").click(e => {
        view.popup.close()
        sketchViewModel.create('polygon')
    })

    function init (mapView) {
        view = mapView
        sketchViewModel = new SketchViewModel({
            view: view,
            layer: polygonGraphicsLayer,
            pointSymbol: {
                type: 'simple-marker',
                color: [255, 255, 255, 0],
                size: '1px',
                outline: { color: 'gray', width: 0 }
            }
        })

        sketchViewModel.on("create", e => {
            if (e.state === "complete") {
                polygonGraphicsLayer.remove(e.graphic)
                selectFeatures(e.graphic.geometry)
            }
        })

        view.map.add(polygonGraphicsLayer)
        view.ui.add('grid', 'bottom-right')

        view.map.layers.forEach(async (l) => {
            if (l.type == "feature") {
                const lv = await view.whenLayerView(l)
                if (lv) layerViews.push(lv)
            }
        })
    }

    async function selectFeatures (geometry) {
        const query = { geometry: geometry }
        view.graphics.removeAll()
        
        layerViews.forEach(l => filterLayer(l, query))
    }

    async function filterLayer (layerView, query) {
        if (layer.type == "feature") {
            const filter = await layerView.queryFeatures(query)

            if (filter.features.length) {
                if (highlight) highlight.remove()
                if (grid) grid.destroy()

                const gridContainer = document.createElement("div")
                const results = document.getElementById("queryResults")
                results.appendChild(gridContainer)

                view.goTo(query.geometry.extent.expand(2))
                highlight = layerView.highlight(filter.features)

                const rows = filter.features.map(f => f.attributes)
                grid = new Grid({
                    columns: Object.keys(rows[0]).map(fieldName => {
                        return {
                            label: fieldName,
                            field: fieldName,
                            sortable: true
                        }
                    })
                }, gridContainer)

                grid.renderArray(rows)
            }
        }
    }

    return {
        init: init
    }
})