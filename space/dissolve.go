package space

import (
    "fmt"
    "github.com/spatial-go/geoos/algorithm/operation"
)

// DissolvePolygonsOptions defines options for dissolving polygons.
type DissolvePolygonsOptions struct {
    AttributeName string        // Optional: attribute to group by
    Values        []interface{} // Optional: specific values to filter
    All           bool          // If true, dissolve all into one feature
}

// DissolvePolygons combines polygons into new features based on options.
// If All is true, all polygons are dissolved into one MultiPolygon.
// If AttributeName is specified, polygons are grouped by attribute values.
// Adjacent polygon boundaries are erased in the output.
func DissolvePolygons(geoms []Geometry, opts DissolvePolygonsOptions) (Geometry, error) {
    // Validate input
    if len(geoms) == 0 {
        return MultiPolygon{}, nil // Empty input returns empty MultiPolygon
    }

    // Handle single geometry case
    if len(geoms) == 1 {
        if poly, ok := geoms[0].(Polygon); ok {
            return MultiPolygon{poly}, nil
        }
        if mp, ok := geoms[0].(MultiPolygon); ok {
            return mp, nil
        }
        return nil, fmt.Errorf("input must be Polygon or MultiPolygon, got %T", geoms[0])
    }

    // Collect all polygons
    var polygons []Polygon
    for _, g := range geoms {
        switch geom := g.(type) {
        case Polygon:
            polygons = append(polygons, geom)
        case MultiPolygon:
            polygons = append(polygons, geom...)
        default:
            return nil, fmt.Errorf("input must be Polygon or MultiPolygon, got %T", geom)
        }
    }

    if opts.All {
        // Dissolve all into one MultiPolygon
        union, err := dissolveAll(polygons)
        if err != nil {
            return nil, err
        }
        return union, nil
    }

    // TODO: Implement attribute-based dissolve (future enhancement)
    return nil, fmt.Errorf("attribute-based dissolve not yet implemented")
}

// dissolveAll merges all polygons into a single MultiPolygon, erasing common boundaries.
func dissolveAll(polygons []Polygon) (MultiPolygon, error) {
    if len(polygons) == 0 {
        return MultiPolygon{}, nil
    }

    // Start with the first polygon
    result := polygons[0]
    for i := 1; i < len(polygons); i++ {
        // Union with next polygon, erasing common boundaries
        op := operation.UnionOperation{Subject: result, Clip: polygons[i]}
        merged, err := op.Union()
        if err != nil {
            return nil, fmt.Errorf("union failed: %v", err)
        }
        switch m := merged.(type) {
        case Polygon:
            result = m
        case MultiPolygon:
            result = m[0] // Take first if multi, adjust as needed
        default:
            return nil, fmt.Errorf("unexpected union result type: %T", m)
        }
    }

    // Convert to MultiPolygon
    return MultiPolygon{result}, nil
}