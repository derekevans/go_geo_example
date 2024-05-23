package main

import (
    "encoding/json"
    "fmt"
    "github.com/go-pg/pg/v10"
    "os"
)

type Metadata struct {
    FarmId      string `json:"FarmID"`
    FarmName    string `json:"FarmName"`
    FieldId     string `json:"FieldID"`
    FieldName   string `json:"FieldName"`
    CropName    string `json:"CropToken"`
}

type Farm struct {
    Id     int64
    Name   string
    Fields []*Field `pg:"rel:has-many"`
}

type Field struct {
    Id     int64
    Name   string
    FarmId int64
    Farm   *Farm `pg:"rel:has-one"`
}

type Crop struct {
    Id int64
    Name string
}

func parseMetadata(path string) Metadata {
    dat, _ := os.ReadFile(path)
    var meta Metadata
    json.Unmarshal(dat, &meta)
    return meta
}

func selectOrInsertFarm(db *pg.DB, meta Metadata) Farm {
    farm := Farm{
	    Name: meta.FarmName,
	}
    _, err := db.Model(&farm).SelectOrInsert(&farm)
    if err != nil {
        panic(err)
    }
    return farm
}

func selectOrInsertField(db *pg.DB, meta Metadata, farm Farm) Field {
    field := Field{
        Name: meta.FieldName,
        FarmId: farm.Id,
    }
    _, err := db.Model(&field).SelectOrInsert(&field)
    if err != nil {
        panic(err)
    }
    return field
}

func selectCrop(db *pg.DB, meta Metadata) *Crop {
    crop := new(Crop)
    err := db.Model(crop).
        Where("name ILIKE ?", meta.CropName).
        First()
    if err != nil {
        panic(err)
    }
    return crop
}

func main() {

    db := pg.Connect(&pg.Options{
        Addr: "postgis:5432",
        User: "postgres",
        Password: "password",
        Database: "go_geo",
    })
    defer db.Close()

    plant_meta := parseMetadata("../data/Harvest/Merriweather Farms-JT-01-Soybeans-Deere-Metadata.json")
    
    farm := selectOrInsertFarm(db, plant_meta)
    field := selectOrInsertField(db, plant_meta, farm)
    crop := selectCrop(db, plant_meta)

    fmt.Println(field.Name)
    fmt.Println(crop.Name)
}
