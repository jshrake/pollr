./node_modules/bower/bin/bower install
./node_modules/.bin/browserify -t debowerify ./src/app/app.js -o src/index.js 
sass ./src/index.scss > index.css
