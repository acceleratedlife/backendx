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
                    -MMA: decimal // valid  MMAs
                    "value"  // history of value. Value is how mucj uBuck you can buy for 1 buck - more number means more valuable currency
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
                "accounts"
                    [account-id]
                        - balance: decimal
                        "transactions"
                            -dateTime
                            -account-id
                            -xrate
                            -amount
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
