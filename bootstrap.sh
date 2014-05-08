sudo apt-get install -y golang
sudo apt-get install -y nginx
sudo apt-get install -y redis-server
sudo apt-get install -y mongod
sudo apt-get install -y nodejs
sudo apt-get install -y npm
sudo apt-get install -y bzr
sudo apt-get install -y ruby-sass

mkdir -p ~/go/src/github.com/jshrake
cd ~/go/src/github.com/jshrake
git clone https://github.com/jshrake/pollr
pushd pollr
pushd rest-server
go get
go build
popd
pushd websocket-server
go get
go build
popd

pushd web/public
npm install
npm install -g bower
./node_modules/bower/bin/bower install
./node_modules/.bin/browserify -t debowerify ./src/app/app.js -o src/index.js 
sass ./src/index.scss > index.css
popd

