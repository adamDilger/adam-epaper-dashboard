package bomsummary

type IconDefinition struct {
	Code                       int
	Precis                     string
	AlternatePrecis            string
	DayIconName                string
	AlternateDayIconName       string
	NightIconName              string
	AlternateNightIconName     string
	DayMood                    IconImageDefinition
	NightMood                  IconImageDefinition
	DayIconMoodName            string
	AlternativeDayIconMoodName string
	NightIconMoodName          string
}

type IconImageDefinition struct {
	ImageName     string
	StartGradient string
	EndGradient   string
}

var IconDefinitions []IconDefinition = []IconDefinition{
	{
		Code:          1,
		Precis:        "Sunny.",
		DayIconName:   "sunny",
		NightIconName: "clear-night",
		DayMood: IconImageDefinition{
			ImageName:     "sunny-day",
			StartGradient: "#4ea5fb",
			EndGradient:   "#aad5ff",
		},
		NightMood: IconImageDefinition{
			ImageName:     "clear-night",
			StartGradient: "#001e3c",
			EndGradient:   "#003b71",
		},
		NightIconMoodName: "clear-night-mood",
	},
	{
		Code:          2,
		Precis:        "Clear.",
		DayIconName:   "sunny",
		NightIconName: "clear-night",
		DayMood: IconImageDefinition{
			ImageName:     "sunny-day",
			StartGradient: "#4ea5fb",
			EndGradient:   "#aad5ff",
		},
		NightMood: IconImageDefinition{
			ImageName:     "clear-night",
			StartGradient: "#001e3c",
			EndGradient:   "#003b71",
		},
		NightIconMoodName: "clear-night-mood",
	},
	{
		Code:                   3,
		Precis:                 "Mostly sunny.",
		AlternatePrecis:        "Mostly clear.",
		DayIconName:            "mostly-sunny",
		AlternateDayIconName:   "partly-cloudy",
		NightIconName:          "mostly-clear-night",
		AlternateNightIconName: "partly-cloudy-night",
		DayMood: IconImageDefinition{
			ImageName:     "mostly-sunny-day",
			StartGradient: "#61b0ff",
			EndGradient:   "#d4ebff",
		},
		NightMood: IconImageDefinition{
			ImageName:     "mostly-clear-night",
			StartGradient: "#001e3c",
			EndGradient:   "#003b71",
		},
		DayIconMoodName:            "mostly-sunny-mood",
		AlternativeDayIconMoodName: "partly-cloudy-mood",
	},
	{
		Code:          4,
		Precis:        "Cloudy.",
		DayIconName:   "cloudy",
		NightIconName: "cloudy-night",
		DayMood: IconImageDefinition{
			ImageName:     "cloudy-day",
			StartGradient: "#d5dadd",
			EndGradient:   "#eeeeee",
		},
		NightMood: IconImageDefinition{
			ImageName:     "cloudy-night",
			StartGradient: "#081623",
			EndGradient:   "#384e60",
		},
		DayIconMoodName: "cloudy-mood",
	},
	{
		Code:          5,
		DayIconName:   "missing",
		NightIconName: "missing-night",
		DayMood: IconImageDefinition{
			ImageName:     "",
			StartGradient: "#DEDEDE",
			EndGradient:   "#DEDEDE",
		},
		NightMood: IconImageDefinition{
			ImageName:     "",
			StartGradient: "#4A4A4A",
			EndGradient:   "#4A4A4A",
		},
	},
	{
		Code:          6,
		Precis:        "Hazy.",
		DayIconName:   "haze",
		NightIconName: "haze-night",
		DayMood: IconImageDefinition{
			ImageName:     "haze-day",
			StartGradient: "#dce2e3",
			EndGradient:   "#ede5d1",
		},
		NightMood: IconImageDefinition{
			ImageName:     "haze-night",
			StartGradient: "#3a3322",
			EndGradient:   "#3c3c3c",
		},
		NightIconMoodName: "haze-night-mood",
	},
	{
		Code:          7,
		DayIconName:   "missing",
		NightIconName: "missing-night",
		DayMood: IconImageDefinition{
			ImageName:     "",
			StartGradient: "#949494",
			EndGradient:   "#949494",
		},
		NightMood: IconImageDefinition{
			ImageName:     "",
			StartGradient: "#a4a4a4",
			EndGradient:   "#a4a4a4",
		},
	},
	{
		Code:            8,
		Precis:          "Light rain.",
		AlternatePrecis: "Possible rain.",
		DayIconName:     "light-rain",
		NightIconName:   "light-rain-night",
		DayMood: IconImageDefinition{
			ImageName:     "light-rain-day",
			StartGradient: "#cccccc",
			EndGradient:   "#eaeaea",
		},
		NightMood: IconImageDefinition{
			ImageName:     "light-rain-night",
			StartGradient: "#081623",
			EndGradient:   "#414e5a",
		},
		DayIconMoodName: "light-rain-mood",
	},
	{
		Code:          9,
		Precis:        "Windy.",
		DayIconName:   "wind",
		NightIconName: "wind-night",
		DayMood: IconImageDefinition{
			ImageName:     "wind-day",
			StartGradient: "#c5e2fd",
			EndGradient:   "#d4dee6",
		},
		NightMood: IconImageDefinition{
			ImageName:     "wind-night",
			StartGradient: "#081623",
			EndGradient:   "#384e60",
		},
		DayIconMoodName: "wind-mood",
	},
	{
		Code:          10,
		Precis:        "Fog.",
		DayIconName:   "fog",
		NightIconName: "fog-night",
		DayMood: IconImageDefinition{
			ImageName:     "fog-day",
			StartGradient: "#cfd8df",
			EndGradient:   "#f2f5f7",
		},
		NightMood: IconImageDefinition{
			ImageName:     "fog-night",
			StartGradient: "#081623",
			EndGradient:   "#384e60",
		},
		NightIconMoodName: "fog-night-mood",
	},
	{
		Code:            11,
		Precis:          "Shower or two.",
		AlternatePrecis: "Showers.",
		DayIconName:     "showers",
		NightIconName:   "showers-night",
		DayMood: IconImageDefinition{
			ImageName:     "showers-day",
			StartGradient: "#d9eafa",
			EndGradient:   "#d0dde8",
		},
		NightMood: IconImageDefinition{
			ImageName:     "showers-night",
			StartGradient: "#081623",
			EndGradient:   "#414e5a",
		},
	},
	{
		Code:            12,
		Precis:          "Rain at times.",
		AlternatePrecis: "Heavy rain.",
		DayIconName:     "rain",
		NightIconName:   "rain-night",
		DayMood: IconImageDefinition{
			ImageName:     "rain-day",
			StartGradient: "#b2b4b5",
			EndGradient:   "#c6c8c9",
		},
		NightMood: IconImageDefinition{
			ImageName:     "rain-night",
			StartGradient: "#081623",
			EndGradient:   "#414e5a",
		},
		DayIconMoodName: "rain-mood",
	},
	{
		Code:          13,
		Precis:        "Dusty.",
		DayIconName:   "dust",
		NightIconName: "dust-night",
		DayMood: IconImageDefinition{
			ImageName:     "dust-day",
			StartGradient: "#fae9c1",
			EndGradient:   "#b9d8ea",
		},
		NightMood: IconImageDefinition{
			ImageName:     "dust-night",
			StartGradient: "#2e2614",
			EndGradient:   "#37393b",
		},
		DayIconMoodName: "dust-mood",
	},
	{
		Code:          14,
		Precis:        "Frost.",
		DayIconName:   "frost",
		NightIconName: "frost-night",
		DayMood: IconImageDefinition{
			ImageName:     "frost-day",
			StartGradient: "#ccddec",
			EndGradient:   "#eef7fa",
		},
		NightMood: IconImageDefinition{
			ImageName:     "frost-night",
			StartGradient: "#081623",
			EndGradient:   "#384e60",
		},
		DayIconMoodName:   "frost-mood",
		NightIconMoodName: "frost-night-mood",
	},
	{
		Code:            15,
		Precis:          "Possible snow.",
		AlternatePrecis: "Snow.",
		DayIconName:     "snow",
		NightIconName:   "snow-night",
		DayMood: IconImageDefinition{
			ImageName:     "snow-day",
			StartGradient: "#d0dbdf",
			EndGradient:   "#f0f0f0",
		},
		NightMood: IconImageDefinition{
			ImageName:     "snow-night",
			StartGradient: "#081623",
			EndGradient:   "#384e60",
		},
		DayIconMoodName:   "snow-mood",
		NightIconMoodName: "snow-night-mood",
	},
	{
		Code:          16,
		Precis:        "Possible storm.",
		DayIconName:   "storms",
		NightIconName: "storms-night",
		DayMood: IconImageDefinition{
			ImageName:     "storms-day",
			StartGradient: "#878d95",
			EndGradient:   "#adb6bc",
		},
		NightMood: IconImageDefinition{
			ImageName:     "storms-night",
			StartGradient: "#525252",
			EndGradient:   "#333333",
		},
		DayIconMoodName: "storms-mood",
	},
	{
		Code:          17,
		Precis:        "Possible shower.",
		DayIconName:   "light-showers",
		NightIconName: "light-showers-night",
		DayMood: IconImageDefinition{
			ImageName:     "light-showers-day",
			StartGradient: "#d9eafa",
			EndGradient:   "#d0dde8",
		},
		NightMood: IconImageDefinition{
			ImageName:     "light-showers-night",
			StartGradient: "#081623",
			EndGradient:   "#414e5a",
		},
	},
	{
		Code:            18,
		Precis:          "Heavy shower or two.",
		AlternatePrecis: "Heavy showers.",
		DayIconName:     "heavy-showers",
		NightIconName:   "heavy-showers-night",
		DayMood: IconImageDefinition{
			ImageName:     "heavy-showers-day",
			StartGradient: "#bdd2e6",
			EndGradient:   "#d4dee6",
		},
		NightMood: IconImageDefinition{
			ImageName:     "heavy-showers-night",
			StartGradient: "#081623",
			EndGradient:   "#414e5a",
		},
	},
	{
		Code:          19,
		Precis:        "Cyclone.",
		DayIconName:   "cyclone",
		NightIconName: "cyclone-night",
		DayMood: IconImageDefinition{
			ImageName:     "cyclone-day",
			StartGradient: "#878d95",
			EndGradient:   "#adb6bc",
		},
		NightMood: IconImageDefinition{
			ImageName:     "cyclone-night",
			StartGradient: "#525252",
			EndGradient:   "#333333",
		},
	},
}

var IconDefinitionMap map[int]IconDefinition

func init() {
	IconDefinitionMap = make(map[int]IconDefinition)
	for _, def := range IconDefinitions {
		IconDefinitionMap[def.Code] = def
	}
}
