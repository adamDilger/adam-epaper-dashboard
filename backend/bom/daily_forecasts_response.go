package bom

import "time"

// DailyForecastsResponse represents the complete BOM daily forecasts API response
type DailyForecastsResponse struct {
	Metadata       DailyForecastMetadata `json:"meta"`
	ForecastsDaily DailyForecastsData    `json:"fcst"`
}

// DailyForecastMetadata contains forecast issue time and timezone information
type DailyForecastMetadata struct {
	IssueTimeUTC     time.Time `json:"issue_time_utc"`
	IssueTimeNextUTC time.Time `json:"issue_time_next_utc"`
	LocalTimezone    string    `json:"local_timezone"`
}

// DailyForecastsData contains the array of daily forecasts
type DailyForecastsData struct {
	Daily []DailyForecastEntry `json:"daily"`
}

// DailyForecastEntry represents a single day's forecast data
type DailyForecastEntry struct {
	DateUTC     time.Time                 `json:"date_utc"`
	Atmospheric DailyAtmosphericForecast  `json:"atm"`
	Terrestrial DailyTerrestrialForecast  `json:"terr"`
	Ocean       DailyOceanForecast        `json:"ocn"`
	Astronomy   DailyAstronomicalForecast `json:"astro"`
}

// DailyAtmosphericForecast contains atmospheric forecast data
type DailyAtmosphericForecast struct {
	SurfaceAir DailySurfaceAirForecast `json:"surf_air"`
}

// DailySurfaceAirForecast contains surface air conditions for the day
type DailySurfaceAirForecast struct {
	TempMaxCel    *float64                   `json:"temp_max_cel"`
	TempMinCel    *float64                   `json:"temp_min_cel"`
	Precipitation DailyPrecipitationForecast `json:"precip"`
	Weather       DailyWeatherForecast       `json:"weather"`
	Radiation     DailyRadiationForecast     `json:"radiation"`
}

// DailyPrecipitationForecast contains detailed precipitation probability data
type DailyPrecipitationForecast struct {
	Exceeding10PercentChanceTotalMm *float64 `json:"exceeding_10percentchance_total_mm"`
	Exceeding25PercentChanceTotalMm *float64 `json:"exceeding_25percentchance_total_mm"`
	Exceeding50PercentChanceTotalMm *float64 `json:"exceeding_50percentchance_total_mm"`
	Exceeding75PercentChanceTotalMm *float64 `json:"exceeding_75percentchance_total_mm"`
	AnyProbabilityPercent           *float64 `json:"any_probability_percent"`
	AnyRestOfDayProbabilityPercent  *float64 `json:"any_restofday_probability_percent"`
	Rain10mmProbabilityPercent      *float64 `json:"10mm_probability_percent"`
	Rain25mmProbabilityPercent      *float64 `json:"25mm_probability_percent"`
}

// DailyWeatherForecast contains weather icon information
type DailyWeatherForecast struct {
	IconCode int `json:"icon_code"`
}

// DailyRadiationForecast contains UV radiation forecast data
type DailyRadiationForecast struct {
	UVClearSkyMaxCode *float64   `json:"uv_clear_sky_max_code"`
	UVPeriodStart     *time.Time `json:"uv_period_start"`
	UVPeriodEnd       *time.Time `json:"uv_period_end"`
}

// DailyTerrestrialForecast contains land-based forecast data
type DailyTerrestrialForecast struct {
	SurfaceLand DailySurfaceLandForecast `json:"surf_land"`
}

// DailySurfaceLandForecast contains surface land conditions
type DailySurfaceLandForecast struct {
	Snow DailySnowForecast `json:"snow"`
}

// DailySnowForecast contains snow-related forecast data (currently empty structure)
type DailySnowForecast struct {
	// Empty struct - no snow data fields in the provided JSON
}

// DailyOceanForecast contains ocean forecast data
type DailyOceanForecast struct {
	SurfaceWater DailySurfaceWaterForecast `json:"surf_water"`
}

// DailySurfaceWaterForecast contains surface water conditions
type DailySurfaceWaterForecast struct {
	Sea DailySeaForecast `json:"sea"`
}

// DailySeaForecast contains sea-related forecast data (currently empty structure)
type DailySeaForecast struct {
	// Empty struct - no sea data fields in the provided JSON
}

// DailyAstronomicalForecast contains astronomical information
type DailyAstronomicalForecast struct {
	SunriseUTC *time.Time `json:"sunrise_utc"`
	SunsetUTC  *time.Time `json:"sunset_utc"`
}
