package bom

import "time"

// ForecastResponse represents the complete BOM forecast API response
type ForecastResponse struct {
	Metadata Metadata `json:"meta"`
	Forecast Forecast `json:"fcst"`
}

// Metadata contains forecast issue time and timezone information
type Metadata struct {
	IssueTimeUTC     time.Time `json:"issue_time_utc"`
	IssueTimeNextUTC time.Time `json:"issue_time_next_utc"`
	LocalTimezone    string    `json:"local_timezone"`
}

// Forecast contains summary and daily forecast data
type Forecast struct {
	Summary ForecastSummary `json:"summary"`
	Daily   []DailyForecast `json:"daily"`
}

// ForecastSummary contains regional forecast summaries
type ForecastSummary struct {
	RegionText           string  `json:"region_text"`
	RegionCoastalText    *string `json:"region_coastal_text"`
	SubRegionCoastalText *string `json:"sub_region_coastal_text"`
	SubRegionText        *string `json:"sub_region_text"`
	PublicDistrictText   *string `json:"public_district_text"`
	SeasText             *string `json:"seas_text"`
	CoastText            *string `json:"coast_text"`
	LocalityText         *string `json:"locality_text"`
}

// DailyForecast contains forecast data for a specific day
type DailyForecast struct {
	DateUTC     time.Time           `json:"date_utc"`
	Atmospheric AtmosphericForecast `json:"atm"`
	Ocean       OceanForecast       `json:"ocn"`
	Terrestrial TerrestrialForecast `json:"terr"`
}

// AtmosphericForecast contains atmospheric conditions forecast
type AtmosphericForecast struct {
	SurfaceAir SurfaceAirForecast `json:"surf_air"`
}

// SurfaceAirForecast contains surface air conditions
type SurfaceAirForecast struct {
	Radiation           RadiationForecast           `json:"radiation"`
	Wind                WindForecast                `json:"wind"`
	Weather             WeatherForecast             `json:"weather"`
	Heatwave            HeatwaveForecast            `json:"heatwave"`
	TropicalSystemState TropicalSystemStateForecast `json:"tropical_system_situation"`
}

// RadiationForecast contains UV and sun protection advice
type RadiationForecast struct {
	AdviceSummary RadiationAdviceSummary `json:"advice_summary"`
}

// RadiationAdviceSummary contains sun protection advice for different areas
type RadiationAdviceSummary struct {
	MetropolitanText   *string `json:"metropolitan_text"`
	PublicDistrictText *string `json:"public_district_text"`
	LocalityText       *string `json:"locality_text"`
}

// WindForecast contains wind conditions and warnings
type WindForecast struct {
	Coastal        WindCoastalForecast `json:"coastal"`
	WarningSummary WindWarningSummary  `json:"warning_summary"`
}

// WindCoastalForecast contains coastal wind conditions
type WindCoastalForecast struct {
	CoastText string `json:"coast_text"`
}

// WindWarningSummary contains wind warnings
type WindWarningSummary struct {
	CoastText *string `json:"coast_text"`
}

// WeatherForecast contains general weather conditions
type WeatherForecast struct {
	PrecisText         string  `json:"precis_text"`
	LocalityText       string  `json:"locality_text"`
	RegionText         *string `json:"region_text"`
	PublicDistrictText *string `json:"public_district_text"`
	MetropolitanText   *string `json:"metropolitan_text"`
	SeasText           *string `json:"seas_text"`
	CoastText          *string `json:"coast_text"`
}

// HeatwaveForecast contains heatwave conditions and warnings
type HeatwaveForecast struct {
	CountryText  *string `json:"country_text"`
	LinkMapImage *string `json:"link_map_image"`
}

// TropicalSystemStateForecast contains tropical system information
type TropicalSystemStateForecast struct {
	CoastText *string `json:"coast_text"`
}

// OceanForecast contains ocean and surf conditions
type OceanForecast struct {
	SurfaceWater SurfaceWaterForecast `json:"surf_water"`
}

// SurfaceWaterForecast contains detailed ocean conditions
type SurfaceWaterForecast struct {
	Caution          SurfCautionForecast `json:"caution"`
	SeaHeightSummary SeaHeightSummary    `json:"sea_height_summary"`
	Swell1stSummary  SwellSummary        `json:"swell_1st_summary"`
	Swell2ndSummary  SwellSummary        `json:"swell_2nd_summary"`
	WaveSummary      WaveSummary         `json:"wave_summary"`
	SurfDanger       SurfDangerForecast  `json:"surf_danger"`
}

// SurfCautionForecast contains surf safety cautions
type SurfCautionForecast struct {
	CoastText *string `json:"coast_text"`
}

// SeaHeightSummary contains sea height information
type SeaHeightSummary struct {
	CoastText *string `json:"coast_text"`
}

// SwellSummary contains swell information
type SwellSummary struct {
	CoastText *string `json:"coast_text"`
}

// WaveSummary contains wave information
type WaveSummary struct {
	CoastText *string `json:"coast_text"`
}

// SurfDangerForecast contains surf danger information
type SurfDangerForecast struct {
	MetropolitanText   *string `json:"metropolitan_text"`
	PublicDistrictText *string `json:"public_district_text"`
	LocalityText       *string `json:"locality_text"`
}

// TerrestrialForecast contains land-based conditions
type TerrestrialForecast struct {
	SurfaceLand SurfaceLandForecast `json:"surf_land"`
}

// SurfaceLandForecast contains surface land conditions
type SurfaceLandForecast struct {
	FireDanger FireDangerForecast `json:"fire_danger"`
}

// FireDangerForecast contains fire danger ratings and information
type FireDangerForecast struct {
	Rating             FireDangerRating `json:"rating"`
	MetropolitanText   *string          `json:"metropolitan_text"`
	PublicDistrictText *string          `json:"public_district_text"`
	LocalityText       *string          `json:"locality_text"`
	RegionText         *string          `json:"region_text"`
}

// FireDangerRating contains fire danger rating codes
type FireDangerRating struct {
	FireDistrictCode   *string `json:"fire_district_code"`
	PublicDistrictCode *string `json:"public_district_code"`
}
