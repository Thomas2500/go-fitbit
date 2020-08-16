# go-fitbit

Fitbit API for Go

The official docs do provide partially different data than provided by the API. This project uses the retuned data of the API as base instead of the official fields defined by the documentation. This fields do use `omitempty` to not break parsing.

Please not that you need to register an app on https://dev.fitbit.com/apps/new to receive an API key.
If you want to use the data only for yourself, you can use "Personal" as OAuth Application Type. This way the API allows access to intraday data like pulse data in second resolution. Otherwise, use "Server" as application type.
You can view your existing apps at https://dev.fitbit.com/apps.

Did I forgot something to implement, found a bug, something changed or recommendations? Please feel free to create an issue or a pull request!

## Installation

```
go get github.com/Thomas2500/go-fitbit
```

## Example

You can find a working example how to use go-fitbit within the folder `example/server/` which shows data of all API endpoints and shows how to use subscriptions where fitbit informs you about user data changes.

## Notes

As of https://dev.fitbit.com/build/reference/web-api/basics/#numerical-ids all IDs should be considered as unsigned int64.

## TODO

Some functions arn't tested because I do not have the hardware for it (I'm only using a Fitbit Versa and MobileTrack of the iPhone app). If you have hardware which provides additional data (alarms or Fitbit Aria) please test the functions and let me know if every works or something needs to be changed.

Functioons explicitly not tested (eventually broken, please test!) or not finished yet:
- activity favorite
- activity TCX (activity with GPS datapoints, hat not the time yet to implement this)
- alarms (no hardware)
- meals (sounds very interesing, seems not to be implemented within the smartphone app and web version?)
- friends (pnly partially tested)
- foods
  - create custom food - https://dev.fitbit.com/build/reference/web-api/food-logging/#create-food
  - delete custom food - https://dev.fitbit.com/build/reference/web-api/food-logging/#delete-custom-food

Further to do:
- combine similar structs and highlight differences

## Findings
- /1/foods/locales.json returns imageUpload true on en_US, but not with other languages like de_DE. No description how it can be used.
- Food search does only return PUBLIC records and no custom stored records. I found no way to find my own records or use them.
- SpO2 data isn't provided by the API

## What I use this API for

I do save the fetched data into MariaDB & InfluxDB databases for further processing and a simple overview in Grafana.
My plan for the future is to further process the data and combine it with other sources of data like weather, movie seen, music listened to, public travel vs driving with car, ...
The possibilities are endless and I do have currently enough free storage to store everything :-)

Big thanks to Fitbit for providing such great devices and API access to the generated data. Unfortunately, this is not a matter of course in the fitness and health sector, although it should be.
