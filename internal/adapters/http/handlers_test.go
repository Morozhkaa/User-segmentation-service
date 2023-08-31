package http

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"segmentation-service/internal/domain/models"
	"segmentation-service/internal/ports/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

var (
	r   *gin.Engine
	svc *mocks.MockSegmentService
)

func TestInit(t *testing.T) {
	// create mock service and router
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc = mocks.NewMockSegmentService(ctrl)
	_, err := New(svc, AdapterOptions{HTTP_port: 3030, Timeout: 10 * time.Second, IdleTimeout: 60 * time.Second})
	require.NoError(t, err)
	r = GetRouter()
}

func TestCreateSegment(t *testing.T) {
	// prepare test data
	testCases := []struct {
		name            string
		inputBody       string
		useMock         bool
		mockBehaviour   func(m *mocks.MockSegmentService)
		expStatusCode   int
		expResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug":"TEST"}`,
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService) {
				m.EXPECT().CreateSegment("TEST").Return(nil)
			},
			expStatusCode:   201,
			expResponseBody: `{"success":"segment with slug 'TEST' created"}`,
		},
		{
			name:            "Incorrect name",
			inputBody:       `{"slug":"# %TEST"}`,
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"invalid format of parameter 'slug'"}`,
		},
		{
			name:            "No slug param",
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"missing required parameters"}`,
		},
		{
			name:      "Slug already exists",
			inputBody: `{"slug":"TEST"}`,
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService) {
				m.EXPECT().CreateSegment("TEST").Return(models.ErrSegmentAlreadyExists)
			},
			expStatusCode:   400,
			expResponseBody: `{"error":"segment with this slug already exists"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"slug":"TEST"}`,
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService) {
				m.EXPECT().CreateSegment("TEST").Return(errors.New("some error"))
			},
			expStatusCode:   500,
			expResponseBody: `{"error":"some error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// change the behavior of the mock service if necessary
			if tc.useMock {
				tc.mockBehaviour(svc)
			}

			// create and execute request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/createSegment", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			r.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.Equal(t, tc.expResponseBody, w.Body.String())
		})
	}
}

func TestDeleteSegment(t *testing.T) {
	// prepare test data
	testCases := []struct {
		name            string
		inputBody       string
		useMock         bool
		mockBehaviour   func(m *mocks.MockSegmentService)
		expStatusCode   int
		expResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug":"TEST"}`,
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService) {
				m.EXPECT().DeleteSegment("TEST").Return(nil)
			},
			expStatusCode:   200,
			expResponseBody: `{"success":"segment with slug 'TEST' deleted"}`,
		},
		{
			name:            "Incorrect name",
			inputBody:       `{"slug":"# %TEST"}`,
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"invalid format of parameter 'slug'"}`,
		},
		{
			name:            "No slug param",
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"missing required parameters"}`,
		},
		{
			name:      "Segment not found",
			inputBody: `{"slug":"TEST"}`,
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService) {
				m.EXPECT().DeleteSegment("TEST").Return(models.ErrSegmentNotFound)
			},
			expStatusCode:   404,
			expResponseBody: `{"error":"segment not found"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"slug":"TEST"}`,
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService) {
				m.EXPECT().DeleteSegment("TEST").Return(errors.New("some error"))
			},
			expStatusCode:   500,
			expResponseBody: `{"error":"some error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// change the behavior of the mock service if necessary
			if tc.useMock {
				tc.mockBehaviour(svc)
			}

			// create and execute request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/deleteSegment", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			r.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.Equal(t, tc.expResponseBody, w.Body.String())
		})
	}
}

func TestUpdateUserSegments(t *testing.T) {
	// prepare test data
	testCases := []struct {
		name            string
		inputBody       string
		userID          string
		data            models.UpdateRequest
		useMock         bool
		mockBehaviour   func(m *mocks.MockSegmentService, data models.UpdateRequest, userID string)
		expStatusCode   int
		expResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"segments-to-add":["TEST1", "TEST2"],"segments-to-remove":[]}`,
			userID:    "550e8400-e29b-41d4-a716-446655440000",
			data:      models.UpdateRequest{SegmentsToAdd: []string{"TEST1", "TEST2"}, SegmentsToRemove: []string{}},
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService, data models.UpdateRequest, userID string) {
				m.EXPECT().UpdateUserSegments(data, uuid.MustParse(userID)).Return(nil)
			},
			expStatusCode:   200,
			expResponseBody: `{"success":"segment information for user with userID = 550e8400-e29b-41d4-a716-446655440000 updated"}`,
		},
		{
			name:            "Invalid uuid format",
			inputBody:       `{"segments-to-add":["TEST1", "TEST2"],"segments-to-remove":[]}`,
			userID:          "123",
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"invalid format of parameter 'userID'"}`,
		},
		{
			name:            "Missing required parameters",
			userID:          "550e8400-e29b-41d4-a716-446655440000",
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"missing required parameters"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"segments-to-add":["TEST1", "TEST2"],"segments-to-remove":[]}`,
			userID:    "550e8400-e29b-41d4-a716-446655440000",
			data:      models.UpdateRequest{SegmentsToAdd: []string{"TEST1", "TEST2"}, SegmentsToRemove: []string{}},
			useMock:   true,
			mockBehaviour: func(m *mocks.MockSegmentService, data models.UpdateRequest, userID string) {
				m.EXPECT().UpdateUserSegments(data, uuid.MustParse(userID)).Return(errors.New("some error"))
			},
			expStatusCode:   500,
			expResponseBody: `{"error":"some error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// change the behavior of the mock service if necessary
			if tc.useMock {
				tc.mockBehaviour(svc, tc.data, tc.userID)
			}

			// create and execute request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/updateUserSegments/%s", tc.userID), bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			r.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.Equal(t, tc.expResponseBody, w.Body.String())
		})
	}
}

func TestGetUserSegments(t *testing.T) {
	// prepare test data
	testCases := []struct {
		name            string
		userID          string
		segments        models.SegmentsList
		useMock         bool
		mockBehaviour   func(m *mocks.MockSegmentService, userID string, segments models.SegmentsList)
		expStatusCode   int
		expResponseBody string
	}{
		{
			name:     "OK",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			segments: models.SegmentsList{S: []string{"TEST1", "TEST2"}},
			useMock:  true,
			mockBehaviour: func(m *mocks.MockSegmentService, userID string, segments models.SegmentsList) {
				m.EXPECT().GetUserSegments(uuid.MustParse(userID)).Return(segments, nil)
			},
			expStatusCode:   200,
			expResponseBody: `{"segments":["TEST1","TEST2"]}`,
		},
		{
			name:            "Invalid uuid format",
			userID:          "123",
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"invalid format of parameter 'userID'"}`,
		},
		{
			name:            "Invalid uuid format",
			userID:          "123",
			useMock:         false,
			expStatusCode:   400,
			expResponseBody: `{"error":"invalid format of parameter 'userID'"}`,
		},
		{
			name:     "Internal server error",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			segments: models.SegmentsList{S: []string{"TEST1", "TEST2"}},
			useMock:  true,
			mockBehaviour: func(m *mocks.MockSegmentService, userID string, segments models.SegmentsList) {
				m.EXPECT().GetUserSegments(uuid.MustParse(userID)).Return(segments, errors.New("some error"))
			},
			expStatusCode:   500,
			expResponseBody: `{"error":"database error: some error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// change the behavior of the mock service if necessary
			if tc.useMock {
				tc.mockBehaviour(svc, tc.userID, tc.segments)
			}

			// create and execute request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/getUserSegments/%s", tc.userID), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			r.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.Equal(t, tc.expResponseBody, w.Body.String())
		})
	}
}

func TestGetPeriodFromPath_Bad(t *testing.T) {
	period := "123-123"

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/getReport/%s", period), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	r.ServeHTTP(w, req)

	// check response
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid format of parameter 'period'"}`, w.Body.String())
}
