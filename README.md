# OTGSITEBUILDER

## API

* create post
* get post
* modify post
* delete post

* create page
* get page
* modify page
* delete page

* add category to post
* remove category to post

* build site



/post

POST

returns {"result": "done", "id": <int>}
fails {"result": "fail", "reason": <string>}

/posts

GET

returns {"result": "done", "count": <int>}
fails {"result": "fail", "reason": <string>}

/posts/<posts-per-page:int>/<page:int>

GET

returns {"result": "done", "posts": [<id:int>, ...]}
fails {"result": "fail", "reason": <string>}

/post/id

GET

returns {"result": "done", "date": <date>, "content": <string>, "title": <string>}
fails {"result": "fail", "reason": <string>}

POST

send {"update": [(date|content|title), ...], "body": <string>, "date": <date>, "title": <string>}
returns {"result": "done"}
fails {"result": "fail", "reason": <string>}

DELETE

returns {"result": "done"}
fails {"result": "fail", "reason": <string>}

/page

POST

returns {"result": "done", "id": <int>}
fails {"result": "fail", "reason": <string>}

/pages

GET

returns {"result": "done", "pages": [<int>, ...]}
fails {"result": "fail", "reason": <string>}

/page/id

GET

returns {"result": "done", "date": <date>, "content": <string>, "title": <string>}
fails {"result": "fail", "reason": <string>}

POST

send {"update": [(content|title), ...], "body": <string>, "title": <string>}
returns {"result": "done"}
fails {"result": "fail", "reason": <string>}

DELETE

returns {"result": "done"}
fails {"result": "fail", "reason": <string>}


This is a static web generator which it takes another different approach
for metadata. Instead of store metadata on posts (frontmatter), store them in a database.
So you use commands and/or paramaters to manipulate metadata related with
you content and you can use or write content with any tool or program.

The application has two modes: *builder* and *manager*; and always operate
with a folder with content.

With *manager* mode you operate with metadata and content options.
With *builder* mode you contruct the site with the content, metadata and content options.

## Commands

### Base

otgsitebuilder -mode (builder|manager) -content <path to folder>

### Builder

The command below is enough to build your website to the folder *output*, which can be configured.

otgsitebuilder -mode builder -content <path to folder>

### Manager
