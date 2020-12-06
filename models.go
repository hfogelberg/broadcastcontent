package broadcastcontent

import (
	"database/sql"

	"github.com/ifragasatt/goifut"
)

// Infotext can be displayed in the header section and may contain anny type of information which is
// not available in the data flow
type Infotext struct {
	Text        string `json:"text,omitempty"`
	GUID        string `json:"guid,omitempty"`
	User        User   `json:"user,omitempty"`
	MessageType string `json:"messageType,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// User is a struct representing a user with a valid JWT token.
type User struct {
	ID                   int64          `json:"userId"`
	FirstName            string         `json:"userFirstName"`
	LastName             string         `json:"userLastName"`
	UserName             string         `json:"userName,omitempty"`
	BroadcastDisplayName sql.NullString `json:"-"`
	DisplayName          string         `json:"displayName"`
	ProfilePic           string         `json:"userProfilePic,omitempty"`
	ImageRotation        int32          `json:"imageRotation"`
	Slug                 string         `json:"userSlug,omitempty"`
	Email                string         `json:"email"`
	CustomerID           int64          `json:"userCustomerId"`
	RoleName             sql.NullString `json:"roleName,omitempty"`
	Role                 string         `json:"role"`
	RoleID               int64          `json:"roleId"`
	InviteDate           string         `json:"inviteDate"`
	Accepted             bool           `json:"accepted"`
	Language             string         `json:"language"`
	AliasID              int64          `json:"aliasId"`
	AliasName            string         `json:"aliasName"`
	AlisaPic             string         `json:"aliasPic"`
	AliasRotation        int64          `json:"aliasRotation"`
	UserIP               string         `json:"userIP"`
	GeoData              goifut.Geo     `json:"geoData"`
}

// InfoData is used internally, but mapped to Info for null handling
type InfoData struct {
	ID                      int64
	Subject                 string
	Description             string
	StartTime               string
	BroadcastStartTime      sql.NullString
	Reporters               []User
	LanguageCode            sql.NullString
	EnableCarousel          bool
	AutoCollapse            bool
	EndTime                 string
	BroadcastEndTime        sql.NullString
	DefaultOrder            string
	DeletedAtTime           sql.NullString
	ShowDescription         bool
	HasSportsPanel          bool
	Autoscroll              bool
	HasComments             bool
	AllowIfrComments        bool
	AllowAnonComments       bool
	AnonCommentRequireEmail bool
	AnonCommentAcceptTerms  bool
	UserTermsVersion        string
	CustomerID              int64
	AutoArchive             bool
	CustomerShortname       string
	ArchiveAfterDays        int
	EmbedJS                 string
	EmbedHTML               string
	ExpandedMode            bool
	PostsToShow             int
	Syndicate               bool
}

// Sportresult is the current score in a game between two competitors, eg football or ice hockey
type Sportresult struct {
	TeamOneName   string `json:"teamOneName"`
	TeamTwoName   string `json:"teamTwoName"`
	TeamOneLogo   string `json:"teamOneLogo"`
	TeamTwoLogo   string `json:"teamTwoLogo"`
	TeamOneResult int    `json:"teamOneResult"`
	TeamTwoResult int    `json:"teamTwoResult"`
	ArticleID     string `json:"articleId"`
	GUID          string `json:"guid"`
	CreatedAt     string `json:"createdAt"`
}

// Sportexternal is a copy of sports result, used by the client
type Sportexternal struct {
	ArticleID     string `json:"articleId,omitempty"`
	GUID          string `json:"guid,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	TeamOneName   string `json:"teamOneName,omitempty"`
	TeamTwoName   string `json:"teamTwoName,omitempty"`
	TeamOneLogo   string `json:"teamOneLogo,omitempty"`
	TeamTwoLogo   string `json:"teamTwoLogo,omitempty"`
	TeamOneResult string `json:"teamOneResult,omitempty"`
	TeamTwoResult string `json:"teamTwoResult,omitempty"`
}

type HeaderSortorder struct {
	GUID  string `json:"itemGuid"`
	Index int    `json:"index"`
}

type Report struct {
	ID          int64         `json:"id"`
	BroadcastID int64         `json:"broadcastId"`
	GUID        string        `json:"guid"`
	CommentGUID string        `json:"commentGuid,omitempty"`
	Text        string        `json:"text"`
	Comments    []Comment     `json:"comments"`
	IsNew       bool          `json:"isNew"`
	IsAppended  bool          `json:"isAppended"`
	IsPinned    bool          `json:"isPinned"`
	IsPrepost   bool          `json:"isPrepost"`
	IsChanged   bool          `json:"isChanged"`
	User        *User         `json:"user,omitempty"`
	Guest       *GuestProfile `json:"guest,omitempty"`
	CreatedAt   string        `json:"createdAt"`
	UpdatedAt   string        `json:"updatedAt"`
	MessageType string        `json:"messageType"`
}

// CommentUser has written a comment
type CommentUser struct {
	IfrUserID        int64  `json:"ifrUserId"`
	UserName         string `json:"userName"`
	Email            string `json:"email"`
	ProfilePic       string `json:"profilePic"`
	ImgRotation      int32  `json:"imgRotation"`
	Alias            string `json:"alias"`
	AliasProfilePic  string `json:"aliasProfilePic"`
	AliasImgRotation int32  `json:"aliasImgRotation"`
	Token            string `json:"token,omitempty"`
	IsIfragasattUser bool   `json:"isIfragasattUser"`
	GuestUserName    string `json:"guestUserName"`
	GuestEmail       string `json:"guestEmail"`
}

// Comment is a comment to a report or broadcast
type Comment struct {
	ID                int64       `json:"id"`
	ArticleID         string      `json:"articleId"`
	BroadcastID       int64       `json:"broadcastId"`
	CommentGUID       string      `json:"guid"`
	ReportGUID        string      `json:"reportGuid"`
	ParentCommentGUID string      `json:"parentCommentGuid"`
	Status            string      `json:"status"`
	IsPinned          bool        `json:"isPinned"`
	Text              string      `json:"text"`
	User              CommentUser `json:"commentUser"`
	WrittenAt         string      `json:"writtenAt"`
	UpdatedAt         string      `json:"updatedAt,omitempty"`
	CreatedAt         string      `json:"createdAt,omitempty"`
	DeletedAt         string      `json:"deletedAt,omitempty"`
	MessageType       string      `json:"messageType"`
	PublishedByUserID int64       `json:"publishedByUserId"`
}

// GuestProfile is used in the report struct when sending a report.
type GuestProfile struct {
	ID            int64  `json:"id,omitempty"`
	DisplayName   string `json:"displayName"`
	ProfilePic    string `json:"profilePicture"`
	ImageRotation int    `json:"imageRotation"`
}
