package broadcastcontent

import (
	"database/sql"
	"log"

	"github.com/ifragasatt/goifut"
)

// GetBroadcastInfo fetches meta data about a broadcast
func GetBroadcastInfo(db *sql.DB, articleID string) (InfoData, error) {
	var i InfoData

	row := db.QueryRow(`SELECT b.id, subject, description, startTime, l.shortName as languageShortname, enableCarousel,
			autoCollapse, defaultOrder, endTime, b.deletedAt, showDescription, hasSportsPanel, autoscroll,
			hasComments, b.allowIfrComments, b.allowAnonComments, anonCommentRequireEmail, anonCommentAcceptTerms, userTermsVersion, b.customerId,
			b.autoArchive, c.shortName as customerShortName,
			cs.archiveAfterDays, b.embedJS, b.embedHTML, b.expandedMode, b.postsToShow, b.syndicate
		FROM direkt.broadcasts b
		JOIN ifragasatt.customerSignupSettings css ON css.customerId=b.customerId
		JOIN ifragasatt.languages l ON css.languageId=l.id
		JOIN ifragasatt.customers c ON c.id=b.customerId
		JOIN customerSettings cs ON cs.customerId=b.customerId
		WHERE articleId=?`, articleID)

	if err := row.Scan(&i.ID, &i.Subject, &i.Description, &i.BroadcastStartTime, &i.LanguageCode, &i.EnableCarousel,
		&i.AutoCollapse, &i.DefaultOrder, &i.BroadcastEndTime, &i.DeletedAtTime, &i.ShowDescription, &i.HasSportsPanel,
		&i.Autoscroll, &i.HasComments, &i.AllowIfrComments, &i.AllowAnonComments, &i.AnonCommentRequireEmail, &i.AnonCommentAcceptTerms,
		&i.UserTermsVersion, &i.CustomerID, &i.AutoArchive,
		&i.CustomerShortname, &i.ArchiveAfterDays, &i.EmbedJS, &i.EmbedHTML, &i.ExpandedMode, &i.PostsToShow, &i.Syndicate); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error scanning info rows: %s\n", err.Error())
			return i, err
		}
	}

	i.StartTime = goifut.ParseNullString(i.BroadcastStartTime)
	i.EndTime = goifut.ParseNullString(i.BroadcastEndTime)

	return i, nil
}

// GetInfoTexts returns an array of general, free format infotexts, which are displayed in the carousel
func GetInfoTexts(db *sql.DB, articleID string) ([]Infotext, error) {
	var texts []Infotext

	rows, err := db.Query(`SELECT i.text, i.createdAt, i.updatedAt, i.guid, i.userid, u.firstName, u.lastName, u.profilePicture, u.slug
		FROM broadcasts b
		JOIN infotexts i ON b.id = i.broadcastId
		JOIN ifragasatt.users u ON u.id=i.userId
		WHERE b.articleId=? AND i.deletedAt IS NULL`, articleID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error preparing get infotexts: %s\n", err.Error())
			return nil, err
		}
	}
	for rows.Next() {
		var i Infotext
		var u User
		if err := rows.Scan(&i.Text, &i.CreatedAt, &i.UpdatedAt, &i.GUID, &u.ID, &u.FirstName, &u.LastName, &u.ProfilePic, &u.Slug); err != nil {
			log.Printf("Error iterating infotexts %s\n", err.Error())
			return texts, err
		}
		i.User = u
		i.MessageType = MessageTypeInfo
		texts = append(texts, i)
	}

	return texts, nil
}

// GetSportResults returns the current match result for a broadcast
func GetSportResults(db *sql.DB, articleID string) ([]Sportresult, error) {
	var results []Sportresult

	rows, err := db.Query(`SELECT teamOneName, teamTwoName, teamOneLogo, teamTwoLogo, teamOneResult, teamTwoResult, guid, createdAt
		FROM sportResults
		WHERE articleId=?
		AND deletedAt IS NULL`, articleID)
	defer rows.Close()
	if err != nil {
		log.Printf("Error preparing get sports result %s\n", err.Error())
		return results, err
	}

	for rows.Next() {
		var r Sportresult
		if err := rows.Scan(&r.TeamOneName, &r.TeamTwoName, &r.TeamOneLogo, &r.TeamTwoLogo, &r.TeamOneResult, &r.TeamTwoResult, &r.GUID, &r.CreatedAt); err != nil {
			if err != sql.ErrNoRows {
				return results, nil
			}
			log.Printf("Error iterating sport results %s\n", err.Error())
			return results, err
		}
		results = append(results, r)
	}

	return results, nil
}

// GetHeaderSortorder returns order items are sorted in the thread header
func GetHeaderSortorder(db *sql.DB, articleID string) ([]HeaderSortorder, error) {
	var order []HeaderSortorder

	rows, err := db.Query(`SELECT hso.guid, hso.sortorder 
		FROM headerSortorder hso
		JOIN broadcasts b on b.id=hso.broadcastId
		WHERE articleId=?`, articleID)
	if err != nil {
		log.Printf("Error querying header sort order %s\n", err.Error())
		return order, err
	}

	for rows.Next() {
		var o HeaderSortorder

		if err := rows.Scan(&o.GUID, &o.Index); err != nil {
			if err == sql.ErrNoRows {
				return order, nil
			}
			log.Printf("Error iterating header sort order: %s\n", err.Error())
			return order, err
		}
		order = append(order, o)
	}
	return order, nil
}

// GetPublishedThreadComments returns all thread comments published in article to client
func GetPublishedThreadComments(db *sql.DB, articleID string) ([]Comment, error) {
	var comments []Comment
	var rows *sql.Rows
	var err error

	type tempComment struct {
		commentGUID       string
		reportGUID        sql.NullString
		broadcastID       int64
		text              string
		status            string
		isPinned          bool
		parentCommentGUID sql.NullString
		writtenAt         sql.NullString
		createdAt         sql.NullString
		userName          sql.NullString
		email             sql.NullString
		profilePic        sql.NullString
		imgRotation       sql.NullInt32
		alias             sql.NullString
		aliasProfilePic   sql.NullString
		aliasImgRotation  sql.NullInt32
		guestUserName     sql.NullString
		guestEmail        sql.NullString
	}

	rows, err = db.Query(`SELECT c.commentGuid, r.guid, c.broadcastId, c.text, c.status, c.parentCommentGuid, c.isPinned,
			c.createdAt, writtenAt,
			CONCAT(u.firstName, ' ', u.lastName) as userName, u.email, u.profilePicture, u.imageRotation, a.name, a.profilePicture,
			a.imageRotation, cu.userName, cu.email
			FROM comments c
			JOIN broadcasts b ON c.broadcastId=b.id
			LEFT JOIN commentUsers cu ON cu.id=c.commentUserId
			LEFT JOIN ifragasatt.users u ON cu.ifrUserId = u.id
			LEFT JOIN ifragasatt.aliases a ON u.aliasId=a.id
			LEFT JOIN reportsComments rc ON c.id=rc.commentId
			LEFT JOIN reports r ON r.id = rc.reportId
			WHERE b.articleId=?
			AND c.status = "published"
			AND c.deletedAt IS NULL
		`, articleID)

	defer rows.Close()

	if err != nil {
		log.Printf("Error querying comments in broadcast: %s\n", err.Error())
		return comments, err
	}
	for rows.Next() {
		var tc tempComment
		if err := rows.Scan(&tc.commentGUID, &tc.reportGUID, &tc.broadcastID, &tc.text, &tc.status, &tc.parentCommentGUID, &tc.isPinned, &tc.createdAt, &tc.writtenAt,
			&tc.userName, &tc.email,
			&tc.profilePic, &tc.imgRotation, &tc.alias, &tc.aliasProfilePic, &tc.aliasImgRotation, &tc.guestUserName, &tc.guestEmail); err != nil {
			log.Printf("Error iterating comments: %s\n", err.Error())
			return comments, err
		}

		p := Comment{
			CommentGUID:       tc.commentGUID,
			ReportGUID:        goifut.ParseNullString(tc.reportGUID),
			BroadcastID:       tc.broadcastID,
			ParentCommentGUID: goifut.ParseNullString(tc.parentCommentGUID),
			Text:              tc.text,
			Status:            tc.status,
			IsPinned:          tc.isPinned,
			WrittenAt:         goifut.ParseNullString(tc.writtenAt),
			CreatedAt:         goifut.ParseNullString(tc.createdAt),
			User: CommentUser{
				UserName:         goifut.ParseNullString(tc.userName),
				Email:            goifut.ParseNullString(tc.email),
				ProfilePic:       goifut.ParseNullString(tc.profilePic),
				ImgRotation:      goifut.ParseNullInt32(tc.imgRotation),
				Alias:            goifut.ParseNullString(tc.alias),
				AliasProfilePic:  goifut.ParseNullString(tc.aliasProfilePic),
				AliasImgRotation: goifut.ParseNullInt32(tc.aliasImgRotation),
				GuestUserName:    goifut.ParseNullString(tc.guestUserName),
				GuestEmail:       goifut.ParseNullString(tc.guestEmail),
			},
		}

		p.MessageType = "threadComment"

		comments = append(comments, p)
	}

	return comments, nil
}
