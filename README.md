# url-shortner

# HOW To RUN Application
    1 . Using code 
        * Clone repo using git clone
            git clone git@github.com:lahuGunjal/url-shortner.git
        * Go to /urlshortner directory
        * run go mod tidy to download dependency
        * run go run server.go 
    2. Using docker image 
        * docker pull lahugunjal11/my-url-shortner
        * docker run -p 1323:1323 lahugunjal11/my-url-shortner
        * you can also run in detached mod docker run -d -p 1323:1323 lahugunjal11/my-url-shortner

# API Exposed 
    1 /url/create  : Post request you have to pass url and domain in body
    http://localhost:1323/url/create
            ex.{
	    "url":"https://www.youtube.com/watch?v=HH_a6aRO1TE&list=RD0fV1CjD6pRM&index=3",
	    "domainName":"http://localhost:1323"
    }

    2 /url/get/:url : This is the get request you have give the url Id you will get original url in response
    ex . req : http://localhost:1323/url/create/Y2Q3NmU5OGU3ZmM0
    res : https://www.youtube.com/watch?v=HH_a6aRO1TE&list=RD0fV1CjD6pRM&index=3

    3. /:url when you hit create url api you will get this short url if clicked on this it will redirect you to the original link
        http://localhost:1323/Y2Q3NmU5OGU3ZmM0
    
    4. /domainstats
         ex. Req : http://localhost:1323/domainstats
         Res: it will return top 3 domain name that have been shorten the most
         