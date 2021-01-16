# In progress

# TODO

* support manage exercises : delete by filter
* User resolution on overlapping exercises? => enhance files
- a) MANUAL in order to make decisions user should be able to filter the correct exercises to enhance with some others!
* list command: filter for remote+band existance
* support strava http API !!!
* export to garmin
* export to strava
* there should be a function to fix database.
  - compact or upgrade when necessary!
  - keep migration safe upon metaDataExtractor next version!!
* import assets
* console output feedback of
  * import process to be able to print, log, mock interface

# Done

* count in fetchDiff result on export. (when export removed a file, put it back - only the missing ones)
* overlapping ids must BE removed when an excercise is overwritten! or completely recalculated
 should not point to a non existant excercise ever!!!
* build an internal map id-exercise to be operate on overlaps
* automatic file enhancement
* list command: store exercises ordered - database is sorted on every save operation - this might not be efficient when dataset gets large
* list command filter for remoteName
* split metadataextractor from garmin connector !!!
* register overlaps on Save
* summary display about the db
* multiple connectors - to be able to import from my watch or a directory
* garmin-connect authenticate on connect
* console input parser
  - to be able to pass os.Args and call connected callbacks
* do not duplicate stored on remotes
* Recover from any exception of underlying interface implementations of
  Export, Import, MetaDataExtractor
