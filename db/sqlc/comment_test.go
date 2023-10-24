package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRadomComment(ticker_number int64) (Comment, error) {
	return testQueries.CreateComment(context.Background(), CreateCommentParams{
		CommentText:   "very positive yet random comment",
		TicketID:      ticker_number,
		UserCommented: "random user commenting",
	})
}

func TestCreateComment(t *testing.T) {
	ticket, _, _ := createRandomTicket()

	args := CreateCommentParams{
		CommentText:   "Cypod is the best",
		TicketID:      ticket.TicketID,
		UserCommented: "random user",
	}

	comment, err := testQueries.CreateComment(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, args.CommentText, comment.CommentText)
	require.Equal(t, args.TicketID, comment.TicketID)
	require.Equal(t, args.UserCommented, comment.UserCommented)
}

func TestDeleteComment(t *testing.T) {
	ticket, _, _ := createRandomTicket()
	comment, err := createRadomComment(ticket.TicketID)

	err = testQueries.DeleteComment(context.Background(), comment.CommentID)
	require.NoError(t, err)

	comment1, err := testQueries.GetCommentForUpdate(context.Background(), comment.CommentID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, comment1)

}

func TestUpdateComment(t *testing.T) {
	ticket, _, _ := createRandomTicket()
	comment, err := createRadomComment(ticket.TicketID)

	args := UpdateCommentParams{
		CommentText: "very good comment",
		CommentID:   comment.CommentID,
	}

	comment, err = testQueries.UpdateComment(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, comment.CommentText, args.CommentText)

}

func TestListComments(t *testing.T) {

	ticket, _, _ := createRandomTicket()

	var lastComment Comment
	var err error
	for i := 0; i < 10; i++ {
		lastComment, err = createRadomComment(ticket.TicketID)
	}

	args := ListCommentsParams{
		TicketID: lastComment.TicketID,
		Limit:    11,
		Offset:   0,
	}

	comments, err := testQueries.ListComments(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Len(t, comments, 10)
}
