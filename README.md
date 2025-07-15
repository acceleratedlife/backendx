
## Quick Start

Follow these steps to get the project up and running:

1.  **Fork and Clone:** Start by forking this repository and then cloning it to your local machine.
2.  **Install Go:** Ensure you have Go installed. You can find installation instructions here:
    [https://go.dev/doc/install](https://go.dev/doc/install)
3.  **Navigate to Project:** Change directory into the main project folder:
    ```bash
    cd your-project-folder
    ```
    (Replace `your-project-folder` with the actual name of your cloned directory)
4.  **Manage Dependencies:** Initialize and tidy up the Go modules:
    ```bash
    go mod tidy
    ```
5.  **Configuration File:** You may need an `alfcg.yml` file in the project's root directory with content similar to this:
    ```yaml
    adminpassword: qweasd
    serverport: 5000
    seedpassword: qweasd
    emailsmtp: qq@aa.com
    passwordsmtp: qweawd
    production: false
    ```
6.  **Build Executable:** Compile the Go application:
    ```bash
    go build
    ```
7.  **Run Application:** You should now see an executable file (e.g., `backend.exe` on Windows, or `backend` on Linux/macOS). Run it:
    ```bash
    .\backend.exe # On Windows
    ./backend     # On Linux/macOS
    ```

### To Test

To run the project's tests, use the following command:

```bash
go test

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

