## About

DIY command-line utility to schedule and publish tweets

## Building from source

### Build
    git clone github.com/selvan/turret.git
    export GOPATH=`pwd`/turret
    export PATH=$GOPATH/bin:$PATH
    cd turret
    go install github.com/selvan/turret

### Test
    go test github.com/selvan/turret


## Using

### Create a Twitter App
* Login to your twitter account
* Create an app at https://apps.twitter.com/
* Click on the "Keys and Access Tokens" tab in newly created app
* From "Application Settings" get "Consumer Key" and "Consumer Secret"
* From "Your Access Token" get "Access Token" and "Access Token Secret"

### Create CREDENTIALS file
    mkdir ~/.turret
    touch ~/.turret/CREDENTIALS

Contents of ~/.turret/CREDENTIALS

    <Twitter Consumer Key>
    <Twitter Consumer Secret>
    <Twitter Access Token>
    <Twitter Access Token Secret>

### Create tweets.csv with schedule date + time and status, in following format
    2017-Jan-02 10:30,"Test tweet"
    2017-Jan-02 10:33,"Another test tweet"
    2017-Jan-03 8:30,"Next day test tweet"

### Execution
Within the folder of tweets.csv and execute

    turret
