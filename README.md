# backend

## to build locally
go build

## test & coverage
go test -coverprofile cover.out
go tool cover -html=cover.out

## local run
- default admin test@admin.com

# database

"schools"
    [school-id]
    - name
    - city
    - zip
    - addCode
    "teachers"
        [teacher-id]
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
            "accounts"
                [account-id]
                    - balance: decimal
                    "transactions"
                        datetime: {date-time, account-id, xrate, amount}
    "classes"
        [class-id]
        - name
        - period
        - addCode
        "students"
            user-id: ''
                
"users"
    [userName]: UserInfo
        "accounts"
            [account-id]
                - balance: number
                "transactions"
                    datetime: {date-time, account-id, xrate, amount}
"cb"
  - ubuck: decimal 
  - teacher-id: issued

"orders"
  datetime: OrderInfo
