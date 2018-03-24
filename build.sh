cd /home/go/src/beego-webcrawler
git pull
cd ..
cp -r github.com beego-webcrawler
cp -r golang.org beego-webcrawler
cd beego-webcrawler
docker build -t beego-webcrawler:v1 .
docker container stop $(docker container ls)
docker run -d -p 80:80 -p  8088:8088  beego-webcrawler:v1 
