## TODO

* Render main page after login.

* A few bootstrap themes as inspiration:
    * http://bootstrapzero.com/bootstrap-template/facebook
    * http://bootstrapzero.com/bootstrap-template/holo


* things:
    /things/2015/03/20
        /{nanoseconds}
            /meta.json
            /content.txt // Real file
            /resized
                /300x400.png
                /200x300.png

* Rest service that saves links, images, video-page, blog, or pdf.
    * Create handlers:
        POST/GET /users/signup
        POST/GET /users/login
        GET      /users/logout
        POST     /things
        GET      /things/2015/03/12/{id}
                    /meta.json
                    /{id}.txt
                    /resized
                        /300x400.png
                        /200x300.png


* Bookmarklet.

* Web UI.

* iPhone app can read from sync data on S3.