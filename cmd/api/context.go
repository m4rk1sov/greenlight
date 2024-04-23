package main

import (
	"context"
	"greenlight.m4rk1sov.github.com/internal/data"
	"net/http"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextGetUser() retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

// For assignment 2 |
// 					|
//					V

// Define a custom contextKey2 type, with the underlying type string.
type contextKey2 string

// Convert the string "user" to a contextKey2 type and assign it to the userContextKey2
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey2 = contextKey2("user_info")

// The contextSetUser2() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func (app *application) contextSetUser2(r *http.Request, user2 *data.UserInfo) *http.Request {
	ctx2 := context.WithValue(r.Context(), userContextKey2, user2)
	return r.WithContext(ctx2)
}

// The contextGetUser2() retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
func (app *application) contextGetUser2(r *http.Request) *data.UserInfo {
	user2, ok := r.Context().Value(userContextKey2).(*data.UserInfo)
	if !ok {
		panic("missing user value in request context")
	}
	return user2
}
