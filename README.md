# antena

cd /Users/ohnishi/home/go/src/github.com/ohnishi/antena

go install ./backend/...



## book

### fetch
go run github.com/ohnishi/antena/backend/cmd/fetch amazon book --dest /Users/ohnishi/home/go/data/antena/fetch/book

### transform
go run github.com/ohnishi/antena/backend/cmd/transform amazon book --src /Users/ohnishi/home/go/data/antena/fetch/book --dest /Users/ohnishi/home/go/data/antena/transform/book --date 20201222

### publish
go run github.com/ohnishi/antena/backend/cmd/publish amazon book --src /Users/ohnishi/home/go/data/antena/transform/book --dest /Users/ohnishi/home/go/src/github.com/ohnishi/antena/hugo/www/public --date 20201222

### deploy
cd /Users/ohnishi/home/go/src/github.com/ohnishi/antena/hugo
firebase deploy --only hosting:www
