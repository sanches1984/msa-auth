//+build !ci

package test

import (
	"context"
	auth "github.com/sanches1984/auth/proto/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestClientFlow(t *testing.T) {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	ctx := context.Background()
	authService := auth.NewAuthServiceClient(conn)
	manageService := auth.NewManageServiceClient(conn)

	// create user
	user, err := manageService.CreateUser(ctx, &auth.CreateUserRequest{
		Login:    "user123",
		Password: "passwd123",
	})
	require.NoError(t, err)
	require.NotEmpty(t, user.Id)

	// try create the same user
	_, err = manageService.CreateUser(ctx, &auth.CreateUserRequest{
		Login:    "user123",
		Password: "passwd123",
	})
	require.EqualError(t, err, "rpc error: code = AlreadyExists desc = duplicate key value violates unique constraint \"uindex_users_login\"") // todo

	// login user (fail)
	_, err = authService.Login(ctx, &auth.LoginRequest{
		Login:    "user123",
		Password: "passwd111",
		Data:     []byte("some user data"),
	})
	require.EqualError(t, err, "rpc error: code = PermissionDenied desc = incorrect password")

	// login user success
	loginResp, err := authService.Login(ctx, &auth.LoginRequest{
		Login:    "user123",
		Password: "passwd123",
		Data:     []byte("some user data"),
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Access.Token)
	require.NotEmpty(t, loginResp.Refresh.Token)

	// validate token (fail)
	_, err = authService.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: loginResp.Refresh.Token})
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid token")

	// validate token (success)
	validResp, err := authService.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: loginResp.Access.Token})
	require.NoError(t, err)
	require.NotEmpty(t, validResp.SessionId)
	require.NotEmpty(t, validResp.Data)

	// change password
	_, err = authService.ChangePassword(ctx, &auth.ChangePasswordRequest{
		Token:       loginResp.Access.Token,
		NewPassword: "newpwd123",
	})
	require.NoError(t, err)

	// update user data
	_, err = authService.UpdateSessionData(ctx, &auth.UpdateSessionDataRequest{
		Token: loginResp.Access.Token,
		Data:  []byte("new user data"),
	})
	require.NoError(t, err)

	// have a break :)
	time.Sleep(time.Second)

	// get new access token
	newLoginResp, err := authService.NewAccessTokenByRefreshToken(ctx, &auth.NewAccessTokenByRefreshTokenRequest{RefreshToken: loginResp.Refresh.Token})
	require.NoError(t, err)
	require.NotEmpty(t, newLoginResp.Access.Token)
	require.NotEmpty(t, newLoginResp.Refresh.Token)
	require.NotEqual(t, loginResp.Access.Token, newLoginResp.Access.Token)

	// validate old token (fail)
	_, err = authService.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: loginResp.Access.Token})
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid token")

	// validate new token (success)
	_, err = authService.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: newLoginResp.Access.Token})
	require.NoError(t, err)
	require.NotEmpty(t, validResp.SessionId)
	require.NotEmpty(t, []byte("new user data"))

	// logout
	_, err = authService.Logout(ctx, &auth.LogoutRequest{Token: newLoginResp.Access.Token})
	require.NoError(t, err)

	// validate new token (fail)
	_, err = authService.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: newLoginResp.Access.Token})
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid token")

	// get users list
	userList, err := manageService.GetUserList(ctx, &auth.GetUserListRequest{Login: "user123"})
	require.NoError(t, err)
	require.Len(t, userList.Users, 1)
	require.Equal(t, userList.Users[0].Login, "user123")

	// delete user
	_, err = manageService.DeleteUser(ctx, &auth.DeleteUserRequest{Id: user.Id})
	require.NoError(t, err)
}
