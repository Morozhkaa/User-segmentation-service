package tests

import (
	"net/url"
	"segmentation-service/internal/domain/models"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:3000"
)

var u = url.URL{
	Scheme: "http",
	Host:   host,
}

type user struct {
	UserID string `path:"userID"`
}

type month struct {
	Period string `path:"period"`
}

func TestCreateSegmnet(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	// OK - added segments TEST1, TEST2
	e.POST("/api/v1/createSegment").
		WithJSON(models.Segment{
			Slug: "TEST1",
		}).Expect().Status(201)
	e.POST("/api/v1/createSegment").
		WithJSON(models.Segment{
			Slug: "TEST2",
		}).Expect().Status(201)

	// Error - adding a second segment with the same name
	e.POST("/api/v1/createSegment").
		WithJSON(models.Segment{
			Slug: "TEST1",
		}).Expect().Status(400)

	// Error - invalid slug format
	e.POST("/api/v1/createSegment").
		WithJSON(models.Segment{
			Slug: "32 *4",
		}).Expect().Status(400)

	// Error - missing required slug parameter
	e.POST("/api/v1/createSegment").
		Expect().Status(400)
}

func TestUpdateUserSegments(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	// Error - invalid UserID format
	u := user{UserID: "123"}
	e.POST("/api/v1/updateUserSegments/{userID}").WithPathObject(u).
		WithJSON(models.UpdateRequest{
			SegmentsToAdd:    []string{"TEST1", "TEST2"},
			SegmentsToRemove: []string{},
		}).Expect().Status(400)

	// OK - added the user to the TEST1, TEST2 segments
	u = user{UserID: "550e8400-e29b-41d4-a716-446655440000"}
	e.POST("/api/v1/updateUserSegments/{userID}").WithPathObject(u).
		WithJSON(models.UpdateRequest{
			SegmentsToAdd:    []string{"TEST1", "TEST2"},
			SegmentsToRemove: []string{},
		}).Expect().Status(200)

	// Error - attempt to add or remove a user from a non-existent segment
	e.POST("/api/v1/updateUserSegments/{userID}").WithPathObject(u).
		WithJSON(models.UpdateRequest{
			SegmentsToAdd:    []string{"NON-EXISTING-SEGMENT"},
			SegmentsToRemove: []string{},
		}).Expect().Status(404)

	e.POST("/api/v1/updateUserSegments/{userID}").WithPathObject(u).
		WithJSON(models.UpdateRequest{
			SegmentsToAdd:    []string{},
			SegmentsToRemove: []string{"NON-EXISTING-SEGMENT"},
		}).Expect().Status(404)

	// OK - removed the user from the TEST1, TEST2 segments
	e.POST("/api/v1/updateUserSegments/{userID}").WithPathObject(u).
		WithJSON(models.UpdateRequest{
			SegmentsToAdd:    []string{},
			SegmentsToRemove: []string{"TEST1", "TEST2"},
		}).Expect().Status(200)
}

func TestGetUserSegments(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	// OK - successfully got an empty list
	u := user{UserID: "550e8400-e29b-41d4-a716-446655440000"}
	e.GET("/api/v1/getUserSegments/{userID}").WithPathObject(u).
		Expect().Status(200)

	// Error - invalid UserID format
	u = user{UserID: "123"}
	e.GET("/api/v1/getUserSegments/{userID}").WithPathObject(u).
		Expect().Status(400)
}

func TestDeleteSegmnet(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	// OK - removed segments TEST1, TEST2
	e.DELETE("/api/v1/deleteSegment").
		WithJSON(models.Segment{
			Slug: "TEST1",
		}).Expect().Status(200)
	e.DELETE("/api/v1/deleteSegment").
		WithJSON(models.Segment{
			Slug: "TEST2",
		}).Expect().Status(200)

	// Error - attempt to delete non-existent segment
	e.DELETE("/api/v1/deleteSegment").
		WithJSON(models.Segment{
			Slug: "324",
		}).Expect().Status(404)

	// Error - invalid slug format
	e.DELETE("/api/v1/deleteSegment").
		WithJSON(models.Segment{
			Slug: "32 *4",
		}).Expect().Status(400)
}

func TestGetReport(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	// Error - invalid period format
	m := month{Period: "123-123"}
	e.GET("/api/v1/getReport/{period}").WithPathObject(m).
		Expect().Status(400)

	// OK
	m = month{Period: "2023-08"}
	e.GET("/api/v1/getReport/{period}").WithPathObject(m).
		Expect().Status(200)
}
