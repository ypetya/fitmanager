package metadataExtractor

type IMetadataExtractor interface {
	// Method to access
	Extract(file string) (
		// The type of activitiy recorded
		Activity string,
		// Device used for recording the exercise
		Device string,
		// begining time when the exercise started
		Start int64,
		// finish time when the exercise ended
		End int64,
		// count of samples stored
		Samples int64,
		// bands stored in the sample
		// TODO order
		Bands []string,
		// time-stamp when the meta-info parsed
		Created int64,
	)
}
