# backend

## to build locally
go build

## test & coverage
go test -coverprofile cover.out
go tool cover -html cover.out

## local run
- default admin test@admin.com

# database

"schools"
    [school-id]
    - name
    - city
    - zip
    - addCode
    "uBuck"
        -totalCurrency
        -freeCurrency
        -history
    "teachers"
        [teacher-id]
            "buck"
                -name
                -totalCurrency
                -freeCurrency
                -history
            "auctions"
                -bid
                -maxBid
                -description
                -endDate
                -startDate
                -owner_id
                -visibility
            "classes"
                [class-id]
                - name
                - period
                - addCode
                "students"
                    user-id: ''
    "admins"
        user-id: ''
    "students"
        [userName]
            dayPayment: datetime
            "accounts"
                [account-id]
                - balance: decimal
                    "transactions"
                        datetime: {date-time, account-id, xrate, amount}
    "classes"
        [class-id]
        - name: string
        - period: int32
        - addCode: string
        "students"
            user-id: ''
                
"users"
    [userName]: UserInfo

"cb"
    "accounts"
        [account-id]  (ubuck/teacher-id)
        - balance: decimal
        "transactions"
            datetime: {date-time, account-id, xrate, amount}
   

"orders"
  datetime: OrderInfo
