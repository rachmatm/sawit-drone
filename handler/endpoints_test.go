package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetHello(t *testing.T) {
	t.Run("success: returns greeting with user id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := NewServer(repository.NewMockEstateRepository(ctrl))

		ctx, rec := newEchoContext(t, http.MethodGet, "/hello?id=1", nil)

		err := server.GetHello(ctx, generated.GetHelloParams{Id: 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.HelloResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, "Hello User 1", resp.Message)
	})

	t.Run("success: different user id returns correct greeting", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := NewServer(repository.NewMockEstateRepository(ctrl))

		ctx, rec := newEchoContext(t, http.MethodGet, "/hello?id=42", nil)

		err := server.GetHello(ctx, generated.GetHelloParams{Id: 42})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.HelloResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, "Hello User 42", resp.Message)
	})

	t.Run("success: zero id returns correct greeting", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := NewServer(repository.NewMockEstateRepository(ctrl))

		ctx, rec := newEchoContext(t, http.MethodGet, "/hello?id=0", nil)

		err := server.GetHello(ctx, generated.GetHelloParams{Id: 0})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.HelloResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, "Hello User 0", resp.Message)
	})
}
