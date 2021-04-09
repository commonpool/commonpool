package handler

import (
	"cp/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"time"
)

/**
requests	offers	hours in bank	acknowledgements received	acknowledgements sent	note
*/
type GroupRow struct {
	AllRequestCount          int
	AllOfferCount            int
	RequestCount             int
	OfferCount               int
	Credits                  time.Duration
	Note                     string
	AcknowledgementsReceived string
	AcknowledgementsSent     string
}

type GroupRowComparer func(a, b *GroupRow) bool

var DefaultGroupRowComparer GroupRowComparer = func(a, b *GroupRow) bool {
	if a.AllRequestCount != b.AllRequestCount {
		return false
	}
	if a.AllOfferCount != b.AllOfferCount {
		return false
	}
	if a.RequestCount != b.RequestCount {
		return false
	}
	if a.OfferCount != b.OfferCount {
		return false
	}
	if a.Credits != b.Credits {
		return false
	}
	if a.Note != b.Note {
		return false
	}
	if a.AcknowledgementsReceived != b.AcknowledgementsReceived {
		return false
	}
	if a.AcknowledgementsSent != b.AcknowledgementsSent {
		return false
	}
	return true
}

var IgnoreGroupRowComparer GroupRowComparer = func(a, b *GroupRow) bool {
	return true
}

func (g *GroupRow) Clone() *GroupRow {
	return &GroupRow{
		AllRequestCount:          g.AllRequestCount,
		AllOfferCount:            g.AllOfferCount,
		RequestCount:             g.RequestCount,
		OfferCount:               g.OfferCount,
		Credits:                  g.Credits,
		Note:                     g.Note,
		AcknowledgementsReceived: g.AcknowledgementsReceived,
		AcknowledgementsSent:     g.AcknowledgementsSent,
	}
}

type UserRow struct {
	UserID                   string
	RequestCount             int
	OfferCount               int
	Credits                  time.Duration
	AcknowledgementsReceived string
	AcknowledgementsSent     string
	Note                     string
}

type UserRowComparer func(a, b *UserRow) bool

var DefaultUserRowComparer UserRowComparer = func(a, b *UserRow) bool {
	if a.UserID != b.UserID {
		return false
	}
	if a.RequestCount != b.RequestCount {
		return false
	}
	if a.OfferCount != b.OfferCount {
		return false
	}
	if a.Credits != b.Credits {
		return false
	}
	if a.AcknowledgementsReceived != b.AcknowledgementsReceived {
		return false
	}
	if a.AcknowledgementsSent != b.AcknowledgementsSent {
		return false
	}
	if a.Note != b.Note {
		return false
	}
	return true
}

var IgnoreUserRowComparer UserRowComparer = func(a, b *UserRow) bool {
	return true
}

func (r *UserRow) Clone() *UserRow {
	return &UserRow{
		UserID:                   r.UserID,
		RequestCount:             r.RequestCount,
		OfferCount:               r.OfferCount,
		Credits:                  r.Credits,
		AcknowledgementsReceived: r.AcknowledgementsReceived,
		AcknowledgementsSent:     r.AcknowledgementsSent,
		Note:                     r.Note,
	}
}

type Row struct {
	Time          time.Time
	Description   string
	GroupRow      *GroupRow
	UserRows      []*UserRow
	userIndexMap  map[string]int
	UserIDs       []string
	concernsGroup bool
}

func (r *Row) ConcernsUser(userID string) bool {
	for _, rowUserID := range r.UserIDs {
		if userID == rowUserID {
			return true
		}
	}
	return false
}

func (r *Row) ConcernsOneOfUsers(userIDs []string) bool {
	for _, userID := range userIDs {
		if r.ConcernsUser(userID) {
			return true
		}
	}
	return false
}

func (r *Row) ConcernsGroup() bool {
	return r.concernsGroup
}

func (r *Row) FilterUsers(userIDs []string) *Row {
	var newRow = r.Clone()

	var userRows []*UserRow
	var userIndexMap = map[string]int{}

	for _, userID := range userIDs {
		userRow := newRow.getUserRow(userID)
		userRows = append(userRows, userRow)
		userIndexMap[userID] = len(userRows) - 1
	}
	newRow.UserIDs = userIDs
	newRow.userIndexMap = userIndexMap
	newRow.UserRows = userRows
	return newRow
}

func (r *Row) Clone() *Row {

	var clonedUserRows []*UserRow
	for _, userRow := range r.UserRows {
		clonedUserRows = append(clonedUserRows, userRow)
	}

	var clonedUserIndexMap = map[string]int{}
	for key, value := range r.userIndexMap {
		clonedUserIndexMap[key] = value
	}

	var clonedUserIds = make([]string, len(r.UserIDs))
	for i, userID := range r.UserIDs {
		clonedUserIds[i] = userID
	}

	return &Row{
		Time:          r.Time,
		Description:   r.Description,
		GroupRow:      r.GroupRow.Clone(),
		UserRows:      clonedUserRows,
		userIndexMap:  clonedUserIndexMap,
		UserIDs:       clonedUserIds,
		concernsGroup: r.concernsGroup,
	}

}

func (r *Row) AsCombined() *Row {
	var newRow = r.Clone()
	newRow.Description = ""
	newRow.Time = time.Time{}
	return newRow
}

type GroupHistoryQuery struct {
	Users     []string `query:"users"`
	ShowGroup bool     `query:"showGroup"`
}

func NewRow() *Row {
	return &Row{
		GroupRow:     &GroupRow{},
		UserRows:     []*UserRow{},
		userIndexMap: map[string]int{},
	}
}

func (r *Row) GetRow(idx int) *UserRow {
	return r.UserRows[idx]
}

func (r *Row) getUserRowIndex(userID string) int {
	if idx, ok := r.userIndexMap[userID]; ok {
		return idx
	}
	return -1
}

func (r *Row) getUserRow(userID string) *UserRow {
	idx := r.getUserRowIndex(userID)
	if idx == -1 {
		row := &UserRow{
			UserID: userID,
		}
		r.UserRows = append(r.UserRows, row)
		r.userIndexMap[userID] = len(r.UserRows) - 1
		return row
	}
	return r.UserRows[idx]
}

func (r *Row) processAcknowledgement(acknowledgement *api.Acknowledgement) {

	var obj string

	if acknowledgement.Type != api.Other {
		obj = "an acknowledgement"
	} else {
		obj = "a note"
	}

	var ackType string
	switch acknowledgement.Type {
	case api.ThanksObjectLent:
		ackType = "<p>Thanks for lending me an object</p>"
		break
	case api.ThanksObjectGift:
		ackType = "<p>Thanks for the gift (object)</p>"
		break
	case api.ThanksServiceGift:
		ackType = "<p>Thanks for the gift (service)</p>"
		break
	}

	notes := acknowledgement.Notes
	if notes != "" {
		notes = "<p>" + notes + "</p>"
	}

	r.Description = fmt.Sprintf("%s %s sent %s to %s %s",
		acknowledgement.SentBy.Type,
		acknowledgement.SentBy.HTMLLink(),
		obj,
		acknowledgement.SentTo.Type,
		acknowledgement.SentTo.HTMLLink())

	if acknowledgement.SentBy.IsGroup() {
		if acknowledgement.Type != api.Other {
			r.GroupRow.AcknowledgementsSent = fmt.Sprintf("Sent to %s %s:%s",
				acknowledgement.SentTo.Type,
				acknowledgement.SentTo.HTMLLink(),
				ackType)
		}
		if notes != "" {
			r.GroupRow.Note = fmt.Sprintf("Sent to %s %s:<br>%s%s",
				acknowledgement.SentTo.Type,
				acknowledgement.SentTo.HTMLLink(),
				ackType,
				notes)
		}
		r.concernsGroup = true
	}
	if acknowledgement.SentBy.IsUser() {
		userID := acknowledgement.SentBy.GetUserID()
		userRow := r.getUserRow(userID)
		if acknowledgement.Type != api.Other {
			userRow.AcknowledgementsSent = fmt.Sprintf("Sent to %s %s:%s",
				acknowledgement.SentTo.Type,
				acknowledgement.SentTo.HTMLLink(),
				ackType)
		}
		if notes != "" {
			userRow.Note = fmt.Sprintf("Sent to %s %s:%s",
				acknowledgement.SentTo.Type,
				acknowledgement.SentTo.HTMLLink(),
				notes)
		}
		r.UserIDs = append(r.UserIDs, userID)
	}
	if acknowledgement.SentTo.IsGroup() {
		if acknowledgement.Type != api.Other {
			r.GroupRow.AcknowledgementsReceived = fmt.Sprintf("Received from %s %s:%s",
				acknowledgement.SentBy.Type,
				acknowledgement.SentBy.HTMLLink(),
				ackType)
		}
		if notes != "" {
			r.GroupRow.Note = fmt.Sprintf("Received from %s %s:%s",
				acknowledgement.SentBy.Type,
				acknowledgement.SentBy.HTMLLink(),
				notes)
		}
		r.concernsGroup = true
	}
	if acknowledgement.SentTo.IsUser() {
		userID := acknowledgement.SentTo.GetUserID()
		userRow := r.getUserRow(userID)
		if acknowledgement.Type != api.Other {
			userRow.AcknowledgementsReceived = fmt.Sprintf("Received from %s %s:%s",
				acknowledgement.SentBy.Type,
				acknowledgement.SentBy.HTMLLink(),
				ackType,
			)
		}
		if notes != "" {
			userRow.Note = fmt.Sprintf("Received from %s %s:%s",
				acknowledgement.SentBy.Type,
				acknowledgement.SentBy.HTMLLink(),
				notes,
			)
		}
		r.UserIDs = append(r.UserIDs, userID)
	}

}

func (r *Row) processCredits(credits *api.Credits) {
	r.Description = fmt.Sprintf("%s %s sent %s credits to %s %s",
		credits.SentBy.Type,
		credits.SentBy.HTMLLink(),
		credits.Amount.String(),
		credits.SentTo.Type,
		credits.SentTo.HTMLLink(),
	)

	if credits.SentBy.IsGroup() {
		r.GroupRow.Credits = r.GroupRow.Credits - credits.Amount
		r.concernsGroup = true
	} else if credits.SentBy.IsUser() {
		userID := credits.SentBy.GetUserID()
		userRow := r.getUserRow(userID)
		userRow.Credits = userRow.Credits - credits.Amount
		r.UserIDs = append(r.UserIDs, userID)
	}

	if credits.SentTo.IsGroup() {
		r.GroupRow.Credits = r.GroupRow.Credits + credits.Amount
		r.concernsGroup = true
	} else if credits.SentTo.IsUser() {
		userID := credits.SentTo.GetUserID()
		userRow := r.getUserRow(userID)
		userRow.Credits = userRow.Credits + credits.Amount
		r.UserIDs = append(r.UserIDs, userID)
	}

}

func (r *Row) processPost(post *api.Post) {

	r.Description = fmt.Sprintf("%s added post %s to group %s",
		post.Author.HTMLLink(),
		post.HTMLLink(),
		post.Group.HTMLLink(),
	)

	userID := post.AuthorID
	userRow := r.getUserRow(userID)
	r.UserIDs = append(r.UserIDs, userID)
	r.concernsGroup = true

	if post.DeletedAt == nil {
		if post.Type == api.RequestPost {
			r.GroupRow.AllRequestCount++
			userRow.RequestCount++
		} else if post.Type == api.OfferPost {
			r.GroupRow.AllOfferCount++
			userRow.OfferCount++
		}
	} else {
		if post.Type == api.RequestPost {
			r.GroupRow.AllRequestCount--
			userRow.RequestCount--
		} else if post.Type == api.OfferPost {
			r.GroupRow.AllOfferCount--
			userRow.OfferCount--
		}
	}

}

func (r *Row) processPreviousRow(previousRow *Row) {
	r.GroupRow = &GroupRow{
		AllRequestCount: previousRow.GroupRow.AllRequestCount,
		AllOfferCount:   previousRow.GroupRow.AllOfferCount,
		RequestCount:    previousRow.GroupRow.RequestCount,
		OfferCount:      previousRow.GroupRow.OfferCount,
		Credits:         previousRow.GroupRow.Credits,
	}
	r.UserRows = make([]*UserRow, len(previousRow.UserRows))
	for i, previous := range previousRow.UserRows {
		newRow := &UserRow{
			UserID:       previous.UserID,
			RequestCount: previous.RequestCount,
			OfferCount:   previous.OfferCount,
			Credits:      previous.Credits,
		}
		r.UserRows[i] = newRow
		r.userIndexMap[previous.UserID] = i
	}
}

func (r *Row) processEntry(entry *Entry) {
	r.Time = entry.time
	if entry.acknowledgement != nil {
		r.processAcknowledgement(entry.acknowledgement)
	} else if entry.credit != nil {
		r.processCredits(entry.credit)
	} else if entry.post != nil {
		r.processPost(entry.post)
	}
}

type Entry struct {
	time            time.Time
	acknowledgement *api.Acknowledgement
	credit          *api.Credits
	post            *api.Post
}

type RowComparer func(a, b *Row) bool

var groupHistoryRowComparer = func(comparers []RowComparer) RowComparer {
	return func(a, b *Row) bool {
		for _, comparer := range comparers {
			if !comparer(a, b) {
				return false
			}
		}
		return true
	}
}

func (h *Handler) handleGetGroupHistory(c echo.Context) error {

	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	acknowledgements, err := h.acknowledgementStore.GetAllInGroup(group.ID)
	if err != nil {
		return err
	}

	credits, err := h.creditsStore.GetAllForGroup(group.ID)
	if err != nil {
		return err
	}

	posts, err := h.postStore.GetByGroup(group.ID)
	if err != nil {
		return err
	}

	var query GroupHistoryQuery
	if err := c.Bind(&query); err != nil {
		return err
	}

	var entries []*Entry
	for _, acknowledgement := range acknowledgements {
		entries = append(entries, &Entry{
			time:            acknowledgement.CreatedAt,
			acknowledgement: acknowledgement,
		})
	}
	for _, credit := range credits {
		entries = append(entries, &Entry{
			time:   credit.CreatedAt,
			credit: credit,
		})
	}
	for _, post := range posts {
		var createdPost api.Post
		createdPost = *post
		createdPost.DeletedAt = nil
		entries = append(entries, &Entry{
			time: post.CreatedAt,
			post: &createdPost,
		})

		if post.DeletedAt != nil {
			entries = append(entries, &Entry{
				time: *post.DeletedAt,
				post: post,
			})
			entries = append(entries, &Entry{
				time: post.CreatedAt,
				post: post,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].time.Before(entries[j].time)
	})

	var userIdMap = map[string]bool{}
	var userMap = map[string]*api.User{}
	for _, acknowledgement := range acknowledgements {
		if acknowledgement.SentBy.IsUser() {
			userID := acknowledgement.SentBy.GetUserID()
			userIdMap[userID] = true
			userMap[userID] = acknowledgement.SentBy.GetUser()
		}
		if acknowledgement.SentTo.IsUser() {
			userID := acknowledgement.SentTo.GetUserID()
			userIdMap[userID] = true
			userMap[userID] = acknowledgement.SentTo.User
		}
	}
	for _, credit := range credits {
		if credit.SentBy.IsUser() {
			userID := credit.SentBy.GetUserID()
			userIdMap[userID] = true
			userMap[userID] = credit.SentBy.User
		}
		if credit.SentTo.IsUser() {
			userID := credit.SentTo.GetUserID()
			userIdMap[userID] = true
			userMap[userID] = credit.SentTo.User
		}
	}
	for _, post := range posts {
		userIdMap[post.AuthorID] = true
		userMap[post.AuthorID] = post.Author
	}
	var userIds []string
	var allUsers []*api.User

	var nullRow = NewRow()
	for userID, _ := range userIdMap {
		// populate user columns by calling getUserRow
		nullRow.getUserRow(userID)
		userIds = append(userIds, userID)
		allUsers = append(allUsers, userMap[userID])
	}

	var previousRow *Row
	var rows []*Row
	for _, entry := range entries {
		if entry.post != nil && entry.post.Type == api.CommentPost {
			continue
		}
		row := NewRow()
		if previousRow != nil {
			row.processPreviousRow(previousRow)
		} else {
			row.processPreviousRow(nullRow)
		}
		row.processEntry(entry)
		rows = append(rows, row)
		previousRow = row
	}

	var rowUsers = []*api.User{}
	var comparers []RowComparer

	queryUserIdMap := map[string]bool{}
	for _, queryUserID := range query.Users {
		rowUsers = append(rowUsers, userMap[queryUserID])
		queryUserIdMap[queryUserID] = true
		comparers = append(comparers, func(a, b *Row) bool {
			userA := a.getUserRow(queryUserID)
			userB := b.getUserRow(queryUserID)
			return DefaultUserRowComparer(userA, userB)
		})
	}

	if query.ShowGroup {
		comparers = append(comparers, func(a, b *Row) bool {
			return DefaultGroupRowComparer(a.GroupRow, b.GroupRow)
		})
	}

	var rowComparer = groupHistoryRowComparer(comparers)

	var filteredRows []*Row
	var skippedRows []*Row
	var isFirstRow = true
	for _, row := range rows {
		if (query.ShowGroup && row.concernsGroup) || row.ConcernsOneOfUsers(query.Users) {
			if len(skippedRows) > 0 {
				if !rowComparer(skippedRows[len(skippedRows)-1], row) {
					filteredRows = append(filteredRows, skippedRows[len(skippedRows)-1].AsCombined())
				}
			}
			isFirstRow = false
			skippedRows = []*Row{}
			filteredRows = append(filteredRows, row.FilterUsers(query.Users))
		} else {
			if !isFirstRow {
				skippedRows = append(skippedRows, row)
			}
		}
	}

	for i, j := 0, len(filteredRows)-1; i < j; i, j = i+1, j-1 {
		filteredRows[i], filteredRows[j] = filteredRows[j], filteredRows[i]
	}

	type UserOption struct {
		*api.User
		Selected bool
	}

	var options []*UserOption
	for _, user := range allUsers {
		var selected = false
		for _, s := range query.Users {
			if s == user.ID {
				selected = true
			}
		}
		options = append(options, &UserOption{
			User:     user,
			Selected: selected,
		})
	}

	return c.Render(http.StatusOK, "group_history_view", map[string]interface{}{
		"Title":     "Hello",
		"Rows":      filteredRows,
		"Users":     options,
		"RowUsers":  rowUsers,
		"ShowGroup": query.ShowGroup,
	})

}
