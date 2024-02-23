package model

// Example:
// "fieldConfigList": [
//       {
//         "name": "ts",
//         "encodingType": "DICTIONARY",
//         "indexType": "TIMESTAMP",
//         "indexTypes": [
//           "TIMESTAMP"
//         ],
//         "timestampConfig": {
//           "granularities": [
//             "DAY",
//             "WEEK",
//             "MONTH"
//           ]
//         },
//         "indexes": null,
//         "tierOverwrites": null
//       },
//       {
//         "name": "ArrTimeBlk",
//         "encodingType": "DICTIONARY",
//         "indexTypes": [],
//         "indexes": {
//           "inverted": {
//             "enabled": "true"
//           }
//         },
//         "tierOverwrites": {
//           "hotTier": {
//             "encodingType": "DICTIONARY",
//             "indexes": {
//               "bloom": {
//                 "enabled": "true"
//               }
//             }
//           },
//           "coldTier": {
//             "encodingType": "RAW",
//             "indexes": {
//               "text": {
//                 "enabled": "true"
//               }
//             }
//           }
//         }
//       }
//     ]

type FieldConfigList struct {
}
