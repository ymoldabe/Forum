# Forum

![Пример анимированной картинки](images/ll.gif)

**The purpose of the project forum is to make web forum with signing in, logging in, making posts, comments and likes. It uses database library SQLite in order to store all the information. SELECT, CREATE and INSERT queries are used to control it.There is a login session and authentication. Only people who are logged in can create posts, leave comments and likes. Docker is used for this project**

**Objectives :**
+    communication between users.
+   associating categories to posts.
+   liking and disliking posts and comments.
+   filtering posts.

**In this project you should be able to sign in with your email and password (which should contain one uppercase letter, one lowercase letter, special character and at least one digit) and then log in. Create posts, leave comments, like or dislike. Posts can be filtered according to the category and likes.**


# How to clone?
_To use this program, follow these steps:_
+   Download or clone the code from this repository.
+   Navigate to the directory containing the code in the terminal or command prompt.

_To clone, use this command:_
```
git clone git@git.01.alem.school:ymoldabe/forum.git
```
# How to run?
_Run the program using the following command:_
_If you want to run it with docker:_
```
make build
make run-img 
```
_Other option is to run it in your own OS:_
```
make run
```