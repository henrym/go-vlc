// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

type LogPriority uint8

const (
	Info LogPriority = iota
	Error
	Warning
	Debug
)

type EventType int

const (
	MediaMetaChanged EventType = iota
	MediaSubItemAdded
	MediaDurationChanged
	MediaParsedChanged
	MediaFreed
	MediaStateChanged
)

const (
	MediaPlayerMediaChanged EventType = 0x100 + iota
	MediaPlayerNothingSpecial
	MediaPlayerOpening
	MediaPlayerBuffering
	MediaPlayerPlaying
	MediaPlayerPaused
	MediaPlayerStopped
	MediaPlayerForward
	MediaPlayerBackward
	MediaPlayerEndReached
	MediaPlayerEncounteredError
	MediaPlayerTimeChanged
	MediaPlayerPositionChanged
	MediaPlayerSeekableChanged
	MediaPlayerPausableChanged
	MediaPlayerTitleChanged
	MediaPlayerSnapshotTaken
	MediaPlayerLengthChanged
)

const (
	MediaListItemAdded EventType = 0x200 + iota
	MediaListWillAddItem
	MediaListItemDeleted
	MediaListWillDeleteItem
)

const (
	MediaListViewItemAdded EventType = 0x300 + iota
	MediaListViewWillAddItem
	MediaListViewItemDeleted
	MediaListViewWillDeleteItem
)

const (
	MediaListPlayerPlayed EventType = 0x400 + iota
	MediaListPlayerNextItemSet
	MediaListPlayerStopped
)

const (
	MediaDiscovererStarted EventType = 0x500 + iota
	MediaDiscovererEnded
)

const (
	VlmMediaAdded EventType = 0x600 + iota
	VlmMediaRemoved
	VlmMediaChanged
	VlmMediaInstanceStarted
	VlmMediaInstanceStopped
	VlmMediaInstanceStatusInit
	VlmMediaInstanceStatusOpening
	VlmMediaInstanceStatusPlaying
	VlmMediaInstanceStatusPause
	VlmMediaInstanceStatusEnd
	VlmMediaInstanceStatusError
)

type PlaybackMode uint8

const (
	PMDefault PlaybackMode = iota
	PMLoop
	PMRepeat
)

type MarqueeOption uint8

const (
	MOEnable MarqueeOption = iota
	MOText
	MOColor
	MOOpacity
	MOPosition
	MORefresh
	MOSize
	MOTimeout
	MOX
	MOY
)

type AdjustOption uint8

const (
	AOEnable AdjustOption = iota
	AOContrast
	AOBrightness
	AOHue
	AOSaturation
	AOGamma
)

type LogoOption uint8

const (
	LOEnable LogoOption = iota
	LOFile
	LOX
	LOY
	LODelay
	LORepeat
	LOOpacity
	LOPosition
)

type MetaProperty uint8

const (
	MPTitle MetaProperty = iota
	MPArtist
	MPGenre
	MPCopyright
	MPAlbum
	MPTrackNumber
	MPDescription
	MPRating
	MPDate
	MPSetting
	MPURL
	MPLanguage
	MPNowPlaying
	MPPublisher
	MPEncodedBy
	MPArtworkURL
	MPTrackID
)

type MediaState uint8

const (
	MSNothingSpecial MediaState = iota
	MSOpening
	MSBuffering
	MSPlaying
	MSPaused
	MSStopped
	MSEnded
	MSError
)

type MediaOption uint16

const (
	MOTrusted MediaOption = 0x2
	MOUnique  MediaOption = 0x100
)

type TrackType int

const (
	TTUnknown TrackType = -1
	TTAudio   TrackType = 0
	TTVideo   TrackType = 1
	TTText    TrackType = 2
)

type AudioDevice int8

const (
	ADError  AudioDevice = -1
	ADMono   AudioDevice = 1
	ADStereo AudioDevice = 2
	AD2F2R   AudioDevice = 4
	AD3F2R   AudioDevice = 5
	AD5_1    AudioDevice = 6
	AD6_1    AudioDevice = 7
	AD7_1    AudioDevice = 8
	ADSPDIF  AudioDevice = 10
)

type AudioChannel int8

const (
	ACError   AudioChannel = -1
	ACStereo  AudioChannel = 1
	ACRStereo AudioChannel = 2
	ACLeft    AudioChannel = 3
	ACRight   AudioChannel = 4
	ACDolbys  AudioChannel = 5
)
