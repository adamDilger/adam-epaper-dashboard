package bom

import "time"

// WeatherResponse represents the complete BOM weather API response
type WeatherResponse struct {
	Station     Station     `json:"stn"`
	Observation Observation `json:"obs"`
}

// Station contains station information
type Station struct {
	Identity StationIdentity `json:"identity"`
	Location StationLocation `json:"location"`
}

// StationIdentity contains station identification details
type StationIdentity struct {
	BomStationNum   int     `json:"bom_stn_num"`
	RiverStationID  *string `json:"river_stn_id"`
	BomStationName  string  `json:"bom_stn_name"`
	WmoStationID    int     `json:"wmo_stn_id"`
	WigosStationID  *string `json:"wigos_stn_id"`
	HeightAboveMSL  float64 `json:"ht_above_msl"`
	HeightBarometer float64 `json:"ht_barometer"`
}

// StationLocation contains geographical location information
type StationLocation struct {
	LatitudeDegrees  float64 `json:"lat_dec_deg"`
	LongitudeDegrees float64 `json:"long_dec_deg"`
	Timezone         string  `json:"timezone"`
}

// Observation contains all weather observation data
type Observation struct {
	DatetimeUTC   time.Time     `json:"datetime_utc"`
	Temperature   Temperature   `json:"temp"`
	Pressure      Pressure      `json:"pres"`
	Wind          Wind          `json:"wind"`
	Precipitation Precipitation `json:"precip"`
	Visibility    Visibility    `json:"visibility"`
	Cloud         Cloud         `json:"cloud"`
}

// Temperature contains temperature-related measurements
type Temperature struct {
	DryBulb1MinCel          float64   `json:"dry_bulb_1min_cel"`
	Apparent1MinCel         float64   `json:"apparent_1min_cel"`
	DewPoint1MinCel         float64   `json:"dew_pnt_1min_cel"`
	WetBulb1MinAvgCel       float64   `json:"wet_bulb_1min_avg_cel"`
	WetBulbGlobeSunCel      *float64  `json:"wet_bulb_globe_sun_cel"`
	WetBulbGlobeShadeCel    float64   `json:"wet_bulb_globe_shade_cel"`
	WetBulbDepressionCel    float64   `json:"wet_bulb_depression_cel"`
	DryBulbMaxCel           float64   `json:"dry_bulb_max_cel"`
	DryBulbMaxTimeUTC       time.Time `json:"dry_bulb_max_time_utc"`
	DryBulbMinCel           float64   `json:"dry_bulb_min_cel"`
	DryBulbMinTimeUTC       time.Time `json:"dry_bulb_min_time_utc"`
	RelativeHumidityPercent float64   `json:"rel_hum_percent"`
}

// Pressure contains atmospheric pressure measurements
type Pressure struct {
	StationLevelHpa *float64 `json:"stn_lvl_hpa"`
	MeanSeaLevelHpa float64  `json:"msl_hpa"`
	QNHHpa          *float64 `json:"qnh_hpa"`
}

// Wind contains wind-related measurements
type Wind struct {
	Speed10mMps          float64   `json:"speed_10m_mps"`
	Direction10mOrd      string    `json:"dirn_10m_ord"`
	GustSpeed10mMps      float64   `json:"gust_speed_10m_mps"`
	GustDirection10mDegT float64   `json:"gust_dirn_10m_deg_t"`
	GustSpeed10mMaxMps   float64   `json:"gust_speed_10m_max_mps"`
	Gust10mMaxUTC        time.Time `json:"gust_10m_max_utc"`
	Run2mTotalM          *float64  `json:"run_2m_total_m"`
}

// Precipitation contains rainfall measurements
type Precipitation struct {
	Since0900LCTTotalMm    float64  `json:"since_0900lct_total_mm"`
	Since0000LCTTotalMm    float64  `json:"since_0000lct_total_mm"`
	Hours24_0900LCTTotalMm float64  `json:"24h_0900lct_total_mm"`
	Minutes10TotalMm       float64  `json:"10min_total_mm"`
	Hour1TotalMm           float64  `json:"1h_total_mm"`
	Hours24TotalMm         *float64 `json:"24h_total_mm"`
}

// Visibility contains visibility measurements
type Visibility struct {
	HorizontalM *float64 `json:"horiz_m"`
}

// Cloud contains cloud-related observations
type Cloud struct {
	BaseHeightS1M         *float64 `json:"base_ht_s1_m"`
	BaseHeightS2M         *float64 `json:"base_ht_s2_m"`
	BaseHeightS3M         *float64 `json:"base_ht_s3_m"`
	BaseHeightS4M         *float64 `json:"base_ht_s4_m"`
	BaseHeightS5M         *float64 `json:"base_ht_s5_m"`
	TotalCoverAmtText     *string  `json:"total_cover_amt_text"`
	LowLayerCoverAmtOkta  *float64 `json:"low_layer_cover_amt_okta"`
	LowLayerHeightM       *float64 `json:"low_layer_height_m"`
	MedLayerCoverAmtOkta  *float64 `json:"med_layer_cover_amt_okta"`
	MedLayerHeightM       *float64 `json:"med_layer_height_m"`
	HighLayerCoverAmtOkta *float64 `json:"high_layer_cover_amt_okta"`
	HighLayerHeightM      *float64 `json:"high_layer_height_m"`
}
