# to run spec 
java -jar .\Downloads\openapi-generator-cli-5.3.0.jar generate -i .\Downloads\al.json -g go-server -o .\Documents\backendx
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
                [account-id]  (ubuck/teacher-id/crypto)
                    -MMA: decimal // valid  MMAs
                    -value: decimal  // history of value. Value is how much uBuck you can buy for 1 buck - more number means more valuable currency
                       "date": decimal  
                    -name
                    "transactions"
                        [transaction-id]
                            -dateTime
                            -account-id
                            -xrate
                            -amount
                "CertificateOfDeposit"
                    [cd-id]
                        -principal investmant
                        -mature date
                        -current value
                        -refund value
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
                        [marketData]
                            - title
                            - cost
                            - count
                            - active
                        [buyers]
                            [studentPurchase_id]
                                - active
                                - student_id


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
        "lottery"
            [lotto-id]: date-time
                - odds: int32
                - amount: int32
                - number: int32
                - winner: user-id
                
"users"
    [userName]: UserInfo

"collegeJobs"
  datetime: incomplete

"jobs"
  datetime: incomplete

"negativeEvents"
  datetime: incomplete

"posotiveEvents"
  datetime: incomplete

"cryptos"
    [crypto] (crypto name)
        -UpdateAt: dateTime
        -usd: float32

