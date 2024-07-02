package models

import (
	"time"
)

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

type Variety struct {
    Id int64
    Name string
}

type PlantingPt struct {
    Id int64
    FieldId int64
    Field *Field `pg:"rel:has-one"`
    CropId int64
    Crop *Crop `pg:"rel:has-one"`
    VarietyId int64
    Variety   *Variety `pg:"rel:has-one"`
    Time time.Time
    Section int64
    SwathWidthFt float64
    DistanceFt float64
    HeadingDeg float64
    ElevationFt float64
    TargetRate float64
    AppliedRate float64
    Geog string
}

type HarvestPt struct {
    Id int64
    FieldId int64
    Field *Field `pg:"rel:has-one"`
    CropId int64
    Crop *Crop `pg:"rel:has-one"`
    Time time.Time
    Section int64
    SwathWidthFt float64
    DistanceFt float64
    HeadingDeg float64
    ElevationFt float64
    MoisturePer float64
    YieldBuAc float64
    Geog string
}