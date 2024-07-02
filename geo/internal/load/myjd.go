
// Loads data exported from MyJohnDeere into a PostGIS database
package load

import (
    "github.com/jonas-p/go-shp"
    "encoding/json"
    "path/filepath"
    "fmt"
    geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
    "github.com/go-pg/pg/v10"
    "os"
    "strconv"
    "time"
    "geo/internal/models"
)

// Load data from directory into the PostGIS database.
// The following folder structure and files must be present:
// {inDir}/Planting/*.shp
// {inDir}/Planting/*.json
// {inDir}/Harvest/*.shp
func LoadMyJD(inDir string) {
    db := pg.Connect(&pg.Options{
        Addr: "postgis:5432",
        User: "postgres",
        Password: "password",
        Database: "go_geo",
    })
    defer db.Close()

    tx, _ := db.Begin()
    defer tx.Close()

    fmt.Println("Loading MyJohnDeere data from", inDir)

    plantMetaPath := getMetaPath(inDir)
    plantShpPath := getPlantingPath(inDir)
    harvestShpPath := getHarvestPath(inDir)

    fmt.Println("Parsing metadata from", plantMetaPath)
    plantMeta := parseMetadata(plantMetaPath)
    farm := selectOrInsertFarm(tx, plantMeta.FarmName)
    field := selectOrInsertField(tx, plantMeta.FieldName, farm)
    crop := selectCrop(tx, plantMeta.CropName)
    fmt.Println("Farm:", farm.Name)
    fmt.Println("Field:", field.Name)
    fmt.Println("Crop:", crop.Name)
    
    fmt.Println("Loading planting data from", plantShpPath)
    plantingPts := parsePlantingPts(tx, plantShpPath, field, crop)
    insertInBatches(tx, plantingPts, 5000)
    
    fmt.Println("Loading harvest data from", harvestShpPath)
    harvestPts := parseHarvestPts(harvestShpPath, field, crop)
    insertInBatches(tx, harvestPts, 5000)
    
    commit(tx)

    fmt.Println(len(plantingPts), "planting points successfully loaded.")
    fmt.Println(len(harvestPts), "harvest points successfully loaded.")
    fmt.Println("Load complete!")
}

func getMetaPath(inDir string) string {
    files, _ := filepath.Glob(inDir + "/Planting/*json")
    return files[0]
}

func getPlantingPath(inDir string) string {
    files, _ := filepath.Glob(inDir + "/Planting/*shp")
    return files[0]
}

func getHarvestPath(inDir string) string {
    files, _ := filepath.Glob(inDir + "/Harvest/*shp")
    return files[0]
}

type metadata struct {
    FarmId      string `json:"FarmID"`
    FarmName    string `json:"FarmName"`
    FieldId     string `json:"FieldID"`
    FieldName   string `json:"FieldName"`
    CropName    string `json:"CropToken"`
}

func parseMetadata(path string) metadata {
    dat, _ := os.ReadFile(path)
    var meta metadata
    json.Unmarshal(dat, &meta)
    return meta
}

func selectOrInsertFarm(tx *pg.Tx, name string) models.Farm {
    farm := models.Farm{
	    Name: name,
	}
    _, err := tx.Model(&farm).SelectOrInsert(&farm)
    if err != nil {
        panic(err)
    }
    return farm
}

func selectOrInsertField(tx *pg.Tx, name string, farm models.Farm) models.Field {
    field := models.Field{
        Name: name,
        FarmId: farm.Id,
    }
    _, err := tx.Model(&field).SelectOrInsert(&field)
    if err != nil {
        panic(err)
    }
    return field
}

func selectCrop(tx *pg.Tx, name string) models.Crop {
    crop := new(models.Crop)
    err := tx.Model(crop).
        Where("name ILIKE ?", name).
        First()
    if err != nil {
        panic(err)
    }
    return *crop
}

func parsePlantingPts(tx *pg.Tx, path string, field models.Field, crop models.Crop)  []models.PlantingPt {
    shape, err := shp.Open(path)
    if err != nil { 
        panic(err)
    } 
    defer shape.Close()
    
    varietyMap := make(map[string]models.Variety)
    fieldMap := getFieldMap(shape)
        
    plantingPts := make([]models.PlantingPt, int(shape.AttributeCount()))
    for shape.Next() {
        featureIdx, feature := shape.Shape()
        
        point := getPoint(feature)
        variety := getVariety(tx, shape, featureIdx, fieldMap, varietyMap)
        time, _ := time.Parse(time.RFC3339, shape.ReadAttribute(featureIdx, fieldMap["IsoTime"]))
        section, _ := strconv.ParseInt(shape.ReadAttribute(featureIdx, fieldMap["SECTIONID"]), 10, 64)
        swathWidthFt, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["SWATHWIDTH"]), 64)
        distanceFt, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["DISTANCE"]), 64)
        headingDeg, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["Heading"]), 64)
        elevationFt, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["Elevation"]), 64)
        targetRate, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["AppliedRate"]), 64)
        appliedRate, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["TargetRate"]), 64)

        plantingPt := models.PlantingPt{
            FieldId: field.Id,
            CropId: crop.Id,
            VarietyId: variety.Id,
            Time: time,
            Section: section,
            SwathWidthFt: swathWidthFt,
            DistanceFt: distanceFt,
            HeadingDeg: headingDeg,
            ElevationFt: elevationFt,
            TargetRate: targetRate,
            AppliedRate: appliedRate,
            Geog: point,
        }

        plantingPts[featureIdx] = plantingPt
    }
    return plantingPts
    
}

func parseHarvestPts(path string, field models.Field, crop models.Crop)  []models.HarvestPt {
    shape, err := shp.Open(path)
    if err != nil { 
        panic(err)
    } 
    defer shape.Close()

    fieldMap := getFieldMap(shape)
        
    harvestPts := make([]models.HarvestPt, int(shape.AttributeCount()))
    for shape.Next() {
        featureIdx, feature := shape.Shape()
        
        point := getPoint(feature)
        time, _ := time.Parse(time.RFC3339, shape.ReadAttribute(featureIdx, fieldMap["IsoTime"]))
        section, _ := strconv.ParseInt(shape.ReadAttribute(featureIdx, fieldMap["SECTIONID"]), 10, 64)
        swathWidthFt, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["SWATHWIDTH"]), 64)
        distanceFt, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["DISTANCE"]), 64)
        headingDeg, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["Heading"]), 64)
        elevationFt, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["Elevation"]), 64)
        moisturePer, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["Moisture"]), 64)
        yieldBuAc, _ := strconv.ParseFloat(shape.ReadAttribute(featureIdx, fieldMap["VRYIELDVOL"]), 64)

        harvestPt := models.HarvestPt{
            FieldId: field.Id,
            CropId: crop.Id,
            Time: time,
            Section: section,
            SwathWidthFt: swathWidthFt,
            DistanceFt: distanceFt,
            HeadingDeg: headingDeg,
            ElevationFt: elevationFt,
            MoisturePer: moisturePer,
            YieldBuAc: yieldBuAc,
            Geog: point,
        }

        harvestPts[featureIdx] = harvestPt
    }
    
    return harvestPts
}

func getFieldMap(shape *shp.Reader) map[string]int {
    var fieldMap = make(map[string]int)
    for field_idx, field := range shape.Fields() {
        fieldMap[field.String()] = field_idx
    }
    return fieldMap
}

func getPoint(shape shp.Shape) string {
    var point *geom.Point

    switch shape.(type) {
    case *shp.Point:
        shpPoint := shape.(*shp.Point)
        point = geom.NewPoint(geom.XY).MustSetCoords([]float64{shpPoint.X, shpPoint.Y}).SetSRID(4326)
    case *shp.PointZ:
        shpPoint := shape.(*shp.PointZ)
        point = geom.NewPoint(geom.XY).MustSetCoords([]float64{shpPoint.X, shpPoint.Y}).SetSRID(4326)
    }
    ewkbhexGeom, err := ewkbhex.Encode(point, ewkbhex.NDR)
    if err != nil {
        panic(err)
    }
    return ewkbhexGeom
}

func getVariety(db *pg.Tx, shape *shp.Reader, featureIdx int, fieldMap map[string]int, varietyMap map[string]models.Variety) models.Variety {
    varietyName := shape.ReadAttribute(featureIdx, fieldMap["Variety"])
    variety, variety_exists := varietyMap[varietyName]
    if ! variety_exists {
        variety = selectOrInsertVariety(db, varietyName)
    }
    return variety
}

func selectOrInsertVariety(tx *pg.Tx, name string) models.Variety {
    variety := models.Variety{
        Name: name,
    }
    _, err := tx.Model(&variety).SelectOrInsert(&variety)
    if err != nil {
        panic(err)
    }
    return variety
}

func insertInBatches[T any](db *pg.Tx, instances []T, batchSize int) {
    batches := getBatches(instances, batchSize)
    for _, batch := range batches {
        _, err := db.Model(&batch).Insert()
        if err != nil {
            panic(err)
        }
    }
}

func getBatches[T any](items []T, batchSize int) [][]T {
    numBatches := len(items)/batchSize + 1
    batches := make([][]T, 0, numBatches)
    for batchSize < len(items) {
        batches = append(batches, items[0:batchSize])
        items = items[batchSize:]
    }
    return append(batches, items)
}

func commit(tx *pg.Tx) {
    if err := tx.Commit(); err != nil {
        tx.Rollback()
        panic(err)
    }
}
