
## In order to

- remove duplicate trainings
- improve data containing one training
  - missing HR
  - wrong temperature data, (is more accurate on handlebar than on a whatch

## Create a database

### Download all the fit files
  
- garmin connect
- strava
- zwift
- endomondo ??

### Import from arbitrary datasources

- garmin devices
- import from folders

### generate meta-data to manage

#### Meta-data

- single source of truth
- can store different filename for different medium
- only LocalDB source points to a real fit file

##### LocalDB

- what to do when local file is removed ? => nothing for now (maintenance should spot it and offer reimport)
- what to do with files placed next to localDB files, not in index ? => nothing for now (maintenance should spot it)

#### RemoteDB

Operations

- import -> upload exercise To 
- export <- download exercise From
- detect files not exported


### analyze

- automatic
- what data is needed to achieve goal ? -> generate meta-data until complete

## Make human decisions

- create a UI for it
- [ ] Mark files to upload and keep separate sources
- [ ] Generate files with extended data info ???

## Store changes

-upload
  - create activities
  - modify activities
  - remove activities

