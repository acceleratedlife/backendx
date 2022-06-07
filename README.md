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
        "cb"
            "accounts"
                    
                [account-id]  (ubuck/teacher-id)
                    "value"
                       "date": decimal  
                    -name
                    -totalCurrency :decimal
                    -freeCurrency :decimal
                    "history"
                        -trade
                        -date
                    "transactions"
                        [transaction-id]
                            -dateTime
                            -account-id
                            -xrate
                            -amount
        "auctions"
            [auction-id]
                - bid
                - maxBid
                - description
                - endDate
                - startDate
                - owner_id
                - winner_id
                "visibility"
                    class-id: ''
        "teachers"
            [teacher-id]
                "classes"
                    [class-id]
                        - name
                        - period
                        - addCode
                        "students"
                            [user-id]
                                - firstName
                                - lastName
                                - email
                "market"
                    [item_id]
                        - description
                        - cost


        "admins"
            user-id: ''
        "students"
            [userName]
                dayPayment: datetime
                event: datetime
                "bAccounts"
                    [account-id]
                        - balance: decimal
                        "transactions"
                            -dateTime
                            -account-id
                            -xrate
                            -amount
                        "history"
                            -trade
                            -date
                "accounts"
                    [account-id]
                        - balance: decimal
                        "transactions"
                            -dateTime
                            -account-id
                            -xrate
                            -amount
                        "history"
                            -trade
                            -date
                "cAccounts"
                    [account-id]
                        - balance: decimal
                        "transactions"
                            -dateTime
                            -account-id
                            -xrate
                            -amount
                        "history"
                            -trade
                            -date
        "classes"
            [class-id]
                - name: string
                - period: int32
                - addCode: string
                "students"
                    user-id: ''
                
"users"
    [userName]: UserInfo

"orders"
  datetime: OrderInfo
