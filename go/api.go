/*
 * AL
 *
 * This is a simple API
 *
 * API version: 1.0.1
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"
)

// AllApiRouter defines the required methods for binding the api requests to a responses for the AllApi
// The AllApiRouter implementation should parse necessary information from the http request,
// pass the data to a AllApiServicer to perform the required actions, then write the service results to the http response.
type AllApiRouter interface {
	AuthUser(http.ResponseWriter, *http.Request)
	ConfirmEmail(http.ResponseWriter, *http.Request)
	ExchangeRate(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	ResetPassword(http.ResponseWriter, *http.Request)
	SearchAccount(http.ResponseWriter, *http.Request)
	SearchBucks(http.ResponseWriter, *http.Request)
	SearchClass(http.ResponseWriter, *http.Request)
	SearchSchool(http.ResponseWriter, *http.Request)
	SearchStudent(http.ResponseWriter, *http.Request)
	SearchStudentBuck(http.ResponseWriter, *http.Request)
	SearchStudents(http.ResponseWriter, *http.Request)
	UserEdit(http.ResponseWriter, *http.Request)
}

// AllSchoolApiRouter defines the required methods for binding the api requests to a responses for the AllSchoolApi
// The AllSchoolApiRouter implementation should parse necessary information from the http request,
// pass the data to a AllSchoolApiServicer to perform the required actions, then write the service results to the http response.
type AllSchoolApiRouter interface {
	AddCodeClass(http.ResponseWriter, *http.Request)
	RemoveClass(http.ResponseWriter, *http.Request)
	SearchAuctions(http.ResponseWriter, *http.Request)
	SearchMyClasses(http.ResponseWriter, *http.Request)
}

// SchoolAdminApiRouter defines the required methods for binding the api requests to a responses for the SchoolAdminApi
// The SchoolAdminApiRouter implementation should parse necessary information from the http request,
// pass the data to a SchoolAdminApiServicer to perform the required actions, then write the service results to the http response.
type SchoolAdminApiRouter interface {
	SearchAdminTeacherClass(http.ResponseWriter, *http.Request)
}

// StaffApiRouter defines the required methods for binding the api requests to a responses for the StaffApi
// The StaffApiRouter implementation should parse necessary information from the http request,
// pass the data to a StaffApiServicer to perform the required actions, then write the service results to the http response.
type StaffApiRouter interface {
	DeleteAuction(http.ResponseWriter, *http.Request)
	Deleteclass(http.ResponseWriter, *http.Request)
	EditClass(http.ResponseWriter, *http.Request)
	KickClass(http.ResponseWriter, *http.Request)
	MakeAuction(http.ResponseWriter, *http.Request)
	MakeClass(http.ResponseWriter, *http.Request)
	PayTransaction(http.ResponseWriter, *http.Request)
	PayTransactions(http.ResponseWriter, *http.Request)
	SearchAllBucks(http.ResponseWriter, *http.Request)
	SearchAuctionsTeacher(http.ResponseWriter, *http.Request)
	SearchClasses(http.ResponseWriter, *http.Request)
	SearchEvents(http.ResponseWriter, *http.Request)
	SearchTransactions(http.ResponseWriter, *http.Request)
}

// StudentApiRouter defines the required methods for binding the api requests to a responses for the StudentApi
// The StudentApiRouter implementation should parse necessary information from the http request,
// pass the data to a StudentApiServicer to perform the required actions, then write the service results to the http response.
type StudentApiRouter interface {
	AuctionBid(http.ResponseWriter, *http.Request)
	BuckConvert(http.ResponseWriter, *http.Request)
	CryptoConvert(http.ResponseWriter, *http.Request)
	SearchAuctionsStudent(http.ResponseWriter, *http.Request)
	SearchBuckTransaction(http.ResponseWriter, *http.Request)
	SearchCrypto(http.ResponseWriter, *http.Request)
	SearchCryptoTransaction(http.ResponseWriter, *http.Request)
	SearchStudentCrypto(http.ResponseWriter, *http.Request)
	SearchStudentUbuck(http.ResponseWriter, *http.Request)
	StudentAddClass(http.ResponseWriter, *http.Request)
}

// SysAdminApiRouter defines the required methods for binding the api requests to a responses for the SysAdminApi
// The SysAdminApiRouter implementation should parse necessary information from the http request,
// pass the data to a SysAdminApiServicer to perform the required actions, then write the service results to the http response.
type SysAdminApiRouter interface {
	CreateBuck(http.ResponseWriter, *http.Request)
	DeleteAccount(http.ResponseWriter, *http.Request)
	DeleteBuck(http.ResponseWriter, *http.Request)
	DeleteSchool(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	Deletetransaction(http.ResponseWriter, *http.Request)
	EditAccount(http.ResponseWriter, *http.Request)
	EditBuck(http.ResponseWriter, *http.Request)
	EditSchool(http.ResponseWriter, *http.Request)
	MakeAccount(http.ResponseWriter, *http.Request)
	MakeSchool(http.ResponseWriter, *http.Request)
	SearchSchools(http.ResponseWriter, *http.Request)
	SearchTransaction(http.ResponseWriter, *http.Request)
}

// TeacherApiRouter defines the required methods for binding the api requests to a responses for the TeacherApi
// The TeacherApiRouter implementation should parse necessary information from the http request,
// pass the data to a TeacherApiServicer to perform the required actions, then write the service results to the http response.
type TeacherApiRouter interface {
	TeacherAddClass(http.ResponseWriter, *http.Request)
}

// UnregisteredApiRouter defines the required methods for binding the api requests to a responses for the UnregisteredApi
// The UnregisteredApiRouter implementation should parse necessary information from the http request,
// pass the data to a UnregisteredApiServicer to perform the required actions, then write the service results to the http response.
type UnregisteredApiRouter interface {
	Register(http.ResponseWriter, *http.Request)
}

// AllApiServicer defines the api actions for the AllApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AllApiServicer interface {
	AuthUser(context.Context) (ImplResponse, error)
	ConfirmEmail(context.Context, string) (ImplResponse, error)
	ExchangeRate(context.Context, string, string) (ImplResponse, error)
	Login(context.Context, RequestLogin) (ImplResponse, error)
	Logout(context.Context, string) (ImplResponse, error)
	ResetPassword(context.Context, UsersResetPasswordBody) (ImplResponse, error)
	SearchAccount(context.Context, string) (ImplResponse, error)
	SearchBucks(context.Context, string) (ImplResponse, error)
	SearchClass(context.Context, RequestUser) (ImplResponse, error)
	SearchSchool(context.Context, string) (ImplResponse, error)
	SearchStudent(context.Context, RequestUser) (ImplResponse, error)
	SearchStudentBuck(context.Context, string) (ImplResponse, error)
	SearchStudents(context.Context) (ImplResponse, error)
	UserEdit(context.Context, UsersUserBody) (ImplResponse, error)
}

// AllSchoolApiServicer defines the api actions for the AllSchoolApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AllSchoolApiServicer interface {
	AddCodeClass(context.Context, RequestUser) (ImplResponse, error)
	RemoveClass(context.Context, RequestKickClass) (ImplResponse, error)
	SearchAuctions(context.Context, string) (ImplResponse, error)
	SearchMyClasses(context.Context, RequestUser) (ImplResponse, error)
}

// SchoolAdminApiServicer defines the api actions for the SchoolAdminApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type SchoolAdminApiServicer interface {
	SearchAdminTeacherClass(context.Context, RequestUser) (ImplResponse, error)
}

// StaffApiServicer defines the api actions for the StaffApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StaffApiServicer interface {
	DeleteAuction(context.Context, string) (ImplResponse, error)
	Deleteclass(context.Context, RequestUser) (ImplResponse, error)
	EditClass(context.Context, RequestEditClass) (ImplResponse, error)
	KickClass(context.Context, RequestKickClass) (ImplResponse, error)
	MakeAuction(context.Context, string, AuctionsBody) (ImplResponse, error)
	MakeClass(context.Context, RequestMakeClass) (ImplResponse, error)
	PayTransaction(context.Context, TransactionsPayTransactionBody) (ImplResponse, error)
	PayTransactions(context.Context, TransactionsPayTransactionsBody) (ImplResponse, error)
	SearchAllBucks(context.Context, string) (ImplResponse, error)
	SearchAuctionsTeacher(context.Context, string) (ImplResponse, error)
	SearchClasses(context.Context, RequestUser) (ImplResponse, error)
	SearchEvents(context.Context, string) (ImplResponse, error)
	SearchTransactions(context.Context, string) (ImplResponse, error)
}

// StudentApiServicer defines the api actions for the StudentApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StudentApiServicer interface {
	AuctionBid(context.Context, AuctionsPlaceBidBody) (ImplResponse, error)
	BuckConvert(context.Context, string, TransactionsConversionTransactionBody) (ImplResponse, error)
	CryptoConvert(context.Context, string, TransactionCryptoTransactionBody) (ImplResponse, error)
	SearchAuctionsStudent(context.Context, string) (ImplResponse, error)
	SearchBuckTransaction(context.Context, string) (ImplResponse, error)
	SearchCrypto(context.Context, string, string) (ImplResponse, error)
	SearchCryptoTransaction(context.Context, string) (ImplResponse, error)
	SearchStudentCrypto(context.Context, string) (ImplResponse, error)
	SearchStudentUbuck(context.Context, string) (ImplResponse, error)
	StudentAddClass(context.Context, RequestAddClass) (ImplResponse, error)
}

// SysAdminApiServicer defines the api actions for the SysAdminApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type SysAdminApiServicer interface {
	CreateBuck(context.Context, BucksBuckBody1) (ImplResponse, error)
	DeleteAccount(context.Context, string) (ImplResponse, error)
	DeleteBuck(context.Context, string) (ImplResponse, error)
	DeleteSchool(context.Context, string) (ImplResponse, error)
	DeleteUser(context.Context, string) (ImplResponse, error)
	Deletetransaction(context.Context, string) (ImplResponse, error)
	EditAccount(context.Context, AccountsAccountBody) (ImplResponse, error)
	EditBuck(context.Context, BucksBuckBody) (ImplResponse, error)
	EditSchool(context.Context, SchoolsSchoolBody) (ImplResponse, error)
	MakeAccount(context.Context, AccountsAccountBody1) (ImplResponse, error)
	MakeSchool(context.Context, SchoolsSchoolBody1) (ImplResponse, error)
	SearchSchools(context.Context, int32) (ImplResponse, error)
	SearchTransaction(context.Context, string) (ImplResponse, error)
}

// TeacherApiServicer defines the api actions for the TeacherApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type TeacherApiServicer interface {
	TeacherAddClass(context.Context, RequestAddClass) (ImplResponse, error)
}

// UnregisteredApiServicer defines the api actions for the UnregisteredApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UnregisteredApiServicer interface {
	Register(context.Context, RequestRegister) (ImplResponse, error)
}
