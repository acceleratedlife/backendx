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
	DeleteAuction(http.ResponseWriter, *http.Request)
	ExchangeRate(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	MakeAuction(http.ResponseWriter, *http.Request)
	PayTransaction(http.ResponseWriter, *http.Request)
	SearchAccount(http.ResponseWriter, *http.Request)
	SearchAllBucks(http.ResponseWriter, *http.Request)
	SearchClass(http.ResponseWriter, *http.Request)
	SearchClasses(http.ResponseWriter, *http.Request)
	SearchMarketItems(http.ResponseWriter, *http.Request)
	SearchSchool(http.ResponseWriter, *http.Request)
	SearchStudent(http.ResponseWriter, *http.Request)
	SearchStudentBucks(http.ResponseWriter, *http.Request)
	SearchStudents(http.ResponseWriter, *http.Request)
	SearchTeachers(http.ResponseWriter, *http.Request)
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
	GetStudentCount(http.ResponseWriter, *http.Request)
	SearchAdminTeacherClass(http.ResponseWriter, *http.Request)
}
// StaffApiRouter defines the required methods for binding the api requests to a responses for the StaffApi
// The StaffApiRouter implementation should parse necessary information from the http request,
// pass the data to a StaffApiServicer to perform the required actions, then write the service results to the http response.
type StaffApiRouter interface { 
	AuctionApprove(http.ResponseWriter, *http.Request)
	AuctionReject(http.ResponseWriter, *http.Request)
	AuctionsAll(http.ResponseWriter, *http.Request)
	DeleteMarketItem(http.ResponseWriter, *http.Request)
	DeleteStudent(http.ResponseWriter, *http.Request)
	Deleteclass(http.ResponseWriter, *http.Request)
	EditClass(http.ResponseWriter, *http.Request)
	GetSettings(http.ResponseWriter, *http.Request)
	KickClass(http.ResponseWriter, *http.Request)
	MakeClass(http.ResponseWriter, *http.Request)
	MakeMarketItem(http.ResponseWriter, *http.Request)
	MarketItemRefund(http.ResponseWriter, *http.Request)
	MarketItemResolve(http.ResponseWriter, *http.Request)
	PayTransactions(http.ResponseWriter, *http.Request)
	ResetPassword(http.ResponseWriter, *http.Request)
	SearchAuctionsTeacher(http.ResponseWriter, *http.Request)
	SearchEvents(http.ResponseWriter, *http.Request)
	SearchTransactions(http.ResponseWriter, *http.Request)
	SetSettings(http.ResponseWriter, *http.Request)
}
// StudentApiRouter defines the required methods for binding the api requests to a responses for the StudentApi
// The StudentApiRouter implementation should parse necessary information from the http request,
// pass the data to a StudentApiServicer to perform the required actions, then write the service results to the http response.
type StudentApiRouter interface { 
	AuctionBid(http.ResponseWriter, *http.Request)
	BuckConvert(http.ResponseWriter, *http.Request)
	BuyCD(http.ResponseWriter, *http.Request)
	CryptoConvert(http.ResponseWriter, *http.Request)
	LatestLotto(http.ResponseWriter, *http.Request)
	LottoPurchase(http.ResponseWriter, *http.Request)
	MarketItemBuy(http.ResponseWriter, *http.Request)
	PreviousLotto(http.ResponseWriter, *http.Request)
	RefundCD(http.ResponseWriter, *http.Request)
	SearchAuctionsStudent(http.ResponseWriter, *http.Request)
	SearchBuck(http.ResponseWriter, *http.Request)
	SearchBuckTransactions(http.ResponseWriter, *http.Request)
	SearchCDS(http.ResponseWriter, *http.Request)
	SearchCDTransactions(http.ResponseWriter, *http.Request)
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
	Deletetransaction(http.ResponseWriter, *http.Request)
	EditAccount(http.ResponseWriter, *http.Request)
	EditBuck(http.ResponseWriter, *http.Request)
	EditSchool(http.ResponseWriter, *http.Request)
	MakeAccount(http.ResponseWriter, *http.Request)
	MakeSchool(http.ResponseWriter, *http.Request)
	SearchSchools(http.ResponseWriter, *http.Request)
	SearchTransaction(http.ResponseWriter, *http.Request)
}
// UnregisteredApiRouter defines the required methods for binding the api requests to a responses for the UnregisteredApi
// The UnregisteredApiRouter implementation should parse necessary information from the http request,
// pass the data to a UnregisteredApiServicer to perform the required actions, then write the service results to the http response.
type UnregisteredApiRouter interface { 
	Register(http.ResponseWriter, *http.Request)
	ResetStaffPassword(http.ResponseWriter, *http.Request)
}


// AllApiServicer defines the api actions for the AllApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AllApiServicer interface { 
	AuthUser(context.Context) (ImplResponse, error)
	ConfirmEmail(context.Context, string) (ImplResponse, error)
	DeleteAuction(context.Context, string) (ImplResponse, error)
	ExchangeRate(context.Context, string, string) (ImplResponse, error)
	Login(context.Context, RequestLogin) (ImplResponse, error)
	Logout(context.Context, string) (ImplResponse, error)
	MakeAuction(context.Context, RequestMakeAuction) (ImplResponse, error)
	PayTransaction(context.Context, RequestPayTransaction) (ImplResponse, error)
	SearchAccount(context.Context, string) (ImplResponse, error)
	SearchAllBucks(context.Context) (ImplResponse, error)
	SearchClass(context.Context, string) (ImplResponse, error)
	SearchClasses(context.Context) (ImplResponse, error)
	SearchMarketItems(context.Context, string) (ImplResponse, error)
	SearchSchool(context.Context, string) (ImplResponse, error)
	SearchStudent(context.Context, string) (ImplResponse, error)
	SearchStudentBucks(context.Context) (ImplResponse, error)
	SearchStudents(context.Context) (ImplResponse, error)
	SearchTeachers(context.Context) (ImplResponse, error)
	UserEdit(context.Context, RequestUserEdit) (ImplResponse, error)
}


// AllSchoolApiServicer defines the api actions for the AllSchoolApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AllSchoolApiServicer interface { 
	AddCodeClass(context.Context, RequestUser) (ImplResponse, error)
	RemoveClass(context.Context, RequestKickClass) (ImplResponse, error)
	SearchAuctions(context.Context, string) (ImplResponse, error)
	SearchMyClasses(context.Context, string) (ImplResponse, error)
}


// SchoolAdminApiServicer defines the api actions for the SchoolAdminApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type SchoolAdminApiServicer interface { 
	GetStudentCount(context.Context, string) (ImplResponse, error)
	SearchAdminTeacherClass(context.Context, string) (ImplResponse, error)
}


// StaffApiServicer defines the api actions for the StaffApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StaffApiServicer interface { 
	AuctionApprove(context.Context, RequestAuctionAction) (ImplResponse, error)
	AuctionReject(context.Context, string) (ImplResponse, error)
	AuctionsAll(context.Context) (ImplResponse, error)
	DeleteMarketItem(context.Context, string) (ImplResponse, error)
	DeleteStudent(context.Context, string) (ImplResponse, error)
	Deleteclass(context.Context, string) (ImplResponse, error)
	EditClass(context.Context, RequestEditClass) (ImplResponse, error)
	GetSettings(context.Context) (ImplResponse, error)
	KickClass(context.Context, RequestKickClass) (ImplResponse, error)
	MakeClass(context.Context, RequestMakeClass) (ImplResponse, error)
	MakeMarketItem(context.Context, RequestMakeMarketItem) (ImplResponse, error)
	MarketItemRefund(context.Context, RequestMarketRefund) (ImplResponse, error)
	MarketItemResolve(context.Context, RequestMarketRefund) (ImplResponse, error)
	PayTransactions(context.Context, RequestPayTransactions) (ImplResponse, error)
	ResetPassword(context.Context, RequestUser) (ImplResponse, error)
	SearchAuctionsTeacher(context.Context) (ImplResponse, error)
	SearchEvents(context.Context) (ImplResponse, error)
	SearchTransactions(context.Context, string) (ImplResponse, error)
	SetSettings(context.Context, Settings) (ImplResponse, error)
}


// StudentApiServicer defines the api actions for the StudentApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StudentApiServicer interface { 
	AuctionBid(context.Context, RequestAuctionBid) (ImplResponse, error)
	BuckConvert(context.Context, RequestBuckConvert) (ImplResponse, error)
	BuyCD(context.Context, RequestBuyCd) (ImplResponse, error)
	CryptoConvert(context.Context, RequestCryptoConvert) (ImplResponse, error)
	LatestLotto(context.Context) (ImplResponse, error)
	LottoPurchase(context.Context, int32) (ImplResponse, error)
	MarketItemBuy(context.Context, RequestMarketRefund) (ImplResponse, error)
	PreviousLotto(context.Context) (ImplResponse, error)
	RefundCD(context.Context, RequestUser) (ImplResponse, error)
	SearchAuctionsStudent(context.Context) (ImplResponse, error)
	SearchBuck(context.Context, string) (ImplResponse, error)
	SearchBuckTransactions(context.Context) (ImplResponse, error)
	SearchCDS(context.Context) (ImplResponse, error)
	SearchCDTransactions(context.Context) (ImplResponse, error)
	SearchCrypto(context.Context, string) (ImplResponse, error)
	SearchCryptoTransaction(context.Context) (ImplResponse, error)
	SearchStudentCrypto(context.Context) (ImplResponse, error)
	SearchStudentUbuck(context.Context) (ImplResponse, error)
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
	Deletetransaction(context.Context, string) (ImplResponse, error)
	EditAccount(context.Context, AccountsAccountBody) (ImplResponse, error)
	EditBuck(context.Context, BucksBuckBody) (ImplResponse, error)
	EditSchool(context.Context, SchoolsSchoolBody) (ImplResponse, error)
	MakeAccount(context.Context, AccountsAccountBody1) (ImplResponse, error)
	MakeSchool(context.Context, SchoolsSchoolBody1) (ImplResponse, error)
	SearchSchools(context.Context, int32) (ImplResponse, error)
	SearchTransaction(context.Context, string) (ImplResponse, error)
}


// UnregisteredApiServicer defines the api actions for the UnregisteredApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UnregisteredApiServicer interface { 
	Register(context.Context, RequestRegister) (ImplResponse, error)
	ResetStaffPassword(context.Context, RequestUser) (ImplResponse, error)
}
