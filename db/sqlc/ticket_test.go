package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTicket() (ticket Ticket, err error) {

	category, err := createRandomCategory()
	return testQueries.CreateTicket(context.Background(),
		CreateTicketParams{
			Title:        "RANDOM",
			Description:  "RANDOM TEXT",
			Status:       "CLOSED",
			UserAssigned: "RANDOM USER",
			CategoryID:   category.CategoryID,
		})

}
func TestCreateTicket(t *testing.T) {
	//TODO: create random category
	category, err := createRandomCategory()

	args := CreateTicketParams{
		Title:        "TestTicket",
		Description:  "Very bad issue",
		Status:       "Unassigned",
		UserAssigned: "my user",
		CategoryID:   category.CategoryID,
	}
	ticket, err := testQueries.CreateTicket(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, args.Title, ticket.Title)
	require.Equal(t, args.Description, ticket.Description)
	require.Equal(t, args.Status, ticket.Status)
	require.Equal(t, args.UserAssigned, ticket.UserAssigned)
	require.Equal(t, args.CategoryID, ticket.CategoryID)

	require.Empty(t, ticket.AssignedTo)
	require.WithinDuration(t, ticket.CreatedAt, time.Now(), time.Second)
	require.Empty(t, ticket.ClosedAt)
	require.NotEmpty(t, ticket.TicketID)
	require.WithinDuration(t, ticket.UpdatedAt, time.Now(), time.Second)

}

func TestDeleteTicket(t *testing.T) {
	ticket, err := createRandomTicket()
	require.NoError(t, err)

	err = testQueries.DeleteTicket(context.Background(), ticket.TicketID)
	require.NoError(t, err)

	ticket, err = testQueries.GetTicket(context.Background(), ticket.TicketID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, ticket)

}

func TestGetTicket(t *testing.T) {
	ticket1, err := createRandomTicket()
	require.NoError(t, err)

	ticket2, err := testQueries.GetTicket(context.Background(), ticket1.TicketID)
	require.NoError(t, err)
	require.Equal(t, ticket1, ticket2)
}

func TestGetTicketForUpdate(t *testing.T) {
	ticket1, err := createRandomTicket()
	require.NoError(t, err)

	ticket2, err := testQueries.GetTicketForUpdate(context.Background(), ticket1.TicketID)
	require.NoError(t, err)
	require.Equal(t, ticket1, ticket2)
}

func TestUpdateTicket(t *testing.T) {
	ticket1, err := createRandomTicket()

	args := UpdateTicketParams{
		TicketID:  ticket1.TicketID,
		UpdatedAt: time.Now(),
	}

	ticket, err := testQueries.UpdateTicket(context.Background(), args)
	require.NoError(t, err)
	require.WithinDuration(t, ticket.UpdatedAt, time.Now(), time.Second)
}

func TestListTickets(t *testing.T) {
	var lastTicket Ticket
	for i := 0; i < 10; i++ {
		lastTicket, _ = createRandomTicket()
	}

	args := ListTicketsParams{
		UserAssigned: lastTicket.UserAssigned,
		Limit:        5,
		Offset:       0,
	}

	tickets, err := testQueries.ListTickets(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, tickets)
}
