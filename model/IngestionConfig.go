package model

// Example:
// "ingestionConfig": {
// 	"transformConfigs": [
// 	  {
// 		"columnName": "ts",
// 		"transformFunction": "fromEpochDays(DaysSinceEpoch)"
// 	  },
// 	  {
// 		"columnName": "tsRaw",
// 		"transformFunction": "fromEpochDays(DaysSinceEpoch)"
// 	  }
// 	],
// 	"continueOnError": false,
// 	"rowTimeValueCheck": false,
// 	"segmentTimeValueCheck": true
//   }

type IngestionConfig struct {
}
