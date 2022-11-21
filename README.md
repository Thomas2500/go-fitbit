# go-fitbit

Fitbit API for Go

The official docs do provide partially different data than provided by the API. This project uses the retuned data of the API as base instead of the official fields defined by the documentation. This fields do use `omitempty` to not break parsing.
Partially some fields were added in the documentation but are still missing in the Swagger file provided by Fitbit.

Please note that you need to register an app on https://dev.fitbit.com/apps/new to receive an API key to use any functionality provided by this project.
If you want to use the data only for yourself, you can use "Personal" as OAuth application type. This way the API allows access to intraday data like pulse data in second resolution. Otherwise, use "Server" as application type.
You can view your existing apps including your credentials at https://dev.fitbit.com/apps.

This project is provided as-is and should be tested before using in any productive environment.
Did I forgot something to implement, found a bug, something changed or recommendations? Please feel free to create an issue or a pull request!

## Installation

This project can be used as a dependency of your project.
```
go get github.com/Thomas2500/go-fitbit
```

## Example

You can find a working example how to use go-fitbit within the folder [`example/server/`](https://github.com/Thomas2500/go-fitbit/tree/master/example/server) which shows data of most API endpoints, shows how to use subscriptions and how to handle token updates.

Initialisizing a new API session can be done using the fitbit.Config struct. This will return a new API session.
```go
// Create a new fitbit session
fca = fitbit.New(fitbit.Config{
  ClientID:     clientID,
  ClientSecret: clientSecret,
  RedirectURL:  fmt.Sprintf("https://%s/callback", "localhost"),
  Scopes: []string{
    fitbit.ScopeActivity,
    fitbit.ScopeBreathingRate,
    fitbit.ScopeHeartrate,
    fitbit.ScopeLocation,
    fitbit.ScopeNutrition,
    fitbit.ScopeProfile,
    fitbit.ScopeSettings,
    fitbit.ScopeSleep,
    fitbit.ScopeSocial,
    fitbit.ScopeSpO2,
    fitbit.ScopeTemperature,
    fitbit.ScopeWeight,
  },
})
```

## Notes

As of https://dev.fitbit.com/build/reference/web-api/basics/#numerical-ids all IDs should be considered as unsigned int64.

## TODO

Some functions arn't tested because I do not have the hardware for it (I'm only using a Fitbit Versa 1 and MobileTrack of the iPhone app). If you have hardware which provides additional data (alarms, temperature, or Fitbit Aria) please test the functionality and let me know if everything works or something needs to be changed.

Functioons explicitly not tested (eventually broken, please test!) or not finished yet:
- activity favorite
- alarms (no hardware)
- meals (sounds very interesing, seems not to be implemented within the smartphone app and web version?)
- friends (only partially tested)
- foods
  - create custom food - https://dev.fitbit.com/build/reference/web-api/food-logging/#create-food
  - delete custom food - https://dev.fitbit.com/build/reference/web-api/food-logging/#delete-custom-food
- Temperature (no hardware)
- Breathing Rate (no hardware)

Further to do:
- combine similar structs and highlight differences

## Findings
- /1/foods/locales.json returns imageUpload true on en_US, but not with other languages like de_DE. No description how it can be used.
- Food search does only return PUBLIC records and no custom stored records. I found no way to find my own records or use them.

## What I use this API for

I do save the fetched data into MariaDB & InfluxDB databases for further processing and a simple overview in Grafana.
My plan for the future is to further process the data and combine it with other sources of data like weather, movie seen, music listened to, public travel vs driving with car, ...
The possibilities are endless and I do have currently enough free storage to store everything :-)

**Big thanks to Fitbit for providing such great devices and API access to the generated data!** Unfortunately, this is not a matter of course in the fitness and health sector, although it should be.
