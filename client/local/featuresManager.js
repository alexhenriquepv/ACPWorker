define([], () => {

    function fieldInfos (attributes) {
        let fieldInfos = []
        Object.keys(attributes).forEach(k => fieldInfos.push({ fieldName: k }))
        return fieldInfos
    }

    function chainPaths (coordinates) {
        let paths = []
        for (let i = 0; i < coordinates.length; i++) {
            let path = []
            path.push([coordinates[i][0], coordinates[i][1]])
            path.push([coordinates[i + 1][0], coordinates[i + 1][1]])
            paths.push(path)
            i++
        }

        return paths
    }

    function polygonRings(coordinates) {
        let rings = []
        coordinates.forEach(c => rings.push([c[0], c[1]]))
        return rings
    }

    return {
        preparedFeatures: (data) => {
            let features = []
            data.features.forEach((f) => {
                if (f.Geometry) {
                    if (f.Geometry.type == "chain") {
                        features.push({
                            geometry: {
                                type: "polyline",
                                paths: chainPaths(f.Geometry.coordinates)
                            },
                            attributes: f.attributes,
                            popupTemplate: {
                                title: "{id}",
                                content: [
                                    {
                                        type: "fields",
                                        fieldInfos: fieldInfos(f.attributes)
                                    }
                                ]
                            }
                        })
                    }
                    else if (f.Geometry.type == "area") {
                        features.push({
                            geometry: {
                                type: "polygon",
                                rings: polygonRings(f.Geometry.coordinates)
                            },
                            attributes: f.attributes,
                            popupTemplate: {
                                title: "{id}",
                                content: [
                                    {
                                        type: "fields",
                                        fieldInfos: fieldInfos(f.attributes)
                                    }
                                ]
                            }
                        })
                    }
                    else if (f.Geometry.type == "point") {
                        features.push({
                            geometry: {
                                type: "point",
                                longitude: f.Geometry.coordinates[0][0],
                                latitude: f.Geometry.coordinates[0][1]
                            },
                            attributes: f.attributes,
                            popupTemplate: {
                                title: "{id}",
                                content: [
                                    {
                                        type: "fields",
                                        fieldInfos: fieldInfos(f.attributes)
                                    }
                                ]
                            }
                        })
                    }
                }
            })

            return features
        }
    }
})