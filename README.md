# README


** More comming soon .. **







### Authentication

#### Why authenticate?
You should authenticate to get your own api key.
You will need an api key to access this service.

You dont need you own api key to check the status of upstream apis or to register a new user.
You can register as manny keys as you want, eaven with the same email!

#### Getting authenticated
Simply **POST** your name and epost in a json format to **/auth/** 

example : 
>**POST** xxxxx:8080/auth/   
>{                        
>    "name" : "Alice",   
>   "email" : "alice@mail.com"  
>}                

Your email needs @ to be registered as a mail!


You will then recive an api key:
>{  
>"key" :            "sk-envdash-fakeAPIkey...",  
>"createdAt" : "20260317 20:32"  
>}  

key is your api own key!
createdAt is when the api key were created

#### Using your api key
You must include you api key in all requests to this service.
This way the server knows it is you asking for information and you do not have to log inn every time!

**more info about using api keys on the way**

#### Deleting your api key
Simply **DELETE** your api key with : **/auth/{apikey}**
{apikey} is your own api key that you want to delete

example:
>**DELETE** xxxxx:8080/auth/sk-envdash-fakeAPIkey...

you dont need to include anny api key to delete you api keys!

You receive on successfully deleated api key:
>No content fount 204

If you reseave anny other status code then the api key were not deleated!

##### Dependencies
To store api keys we are dependent on Firebase being up. Firebase is used to store the keys in a secure way. This means that if the **/status** handler returns annything other than 200 on the Firebase collum, then api keys cant be authenticated.

##### How the sever creates API keys + security
API keys for users are made from a hash of the registerd email and the time the api is being created.
This ensures that duplicate api keys are very hard to be duplicated and users can create a lot of keys!

Also as an additional layer of security against duplications the time since the server were set up is also a part of the hash. That makes guessing or bruteforcing someone elses api key a lot harder since you would need to know the email, the exact time (.0000001 sec), and for how long the server has been up.
