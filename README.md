# Homemade Screenshotter
Create and publish self-hosted screenshots

Installation
--

The screenshotter consists of 3 parts:
* *Client* part (branch **master**)
* *Uploader* part (branch **uploader**)
* *Static* content part (branch **server**)

So, the main idea of using this tool is the following:
1. You copy text or image in the clipboard (Ctrl+C)
1. You open the *Client* app and press "Upload" button
1. The *Client* app sends the content to *Uploader* app, running on your server.
The *Uploader* app puts uploaded file into some folder, hosted by *nginx*,
and returns you URL.
1. When you open the URL, *nginx* serves it. All assets for text
content is defined by *Server* app part. So, your content is intended
to be downloaded as fast as possible by *nginx*, no ads or any other
dynamic bullshit is used there.

#### Client part
* open Terminal and run `cp .env.dist .env`. Then, fill actual values in .env file - for example
```
UPLOAD_URL=https://screenshots.uploaded.here:3333/upload
TMP_FOLDER=/tmp
ACCESS_KEY=verY_long_secret
```
* open Terminal and run `./install_ubuntu.sh`
* add application icon to Favorites

#### Uploader part
* On your server, pull this repo and checkout *uploader* branch
* In terminal, run `cp .env.dist .env`
* Edit .env with actual settings, for example
```
IMAGE_PATH=/imgout
STATIC_SERVER_PATH=https://screenshots.served.here.by.nginx/i/
ACCESS_KEY=verY_long_secret
LISTEN_ADDR=0.0.0.0:3333
```
The same key in ACCESS_KEY with *Client* here is a must. The port should be any free.
The STATIC_SERVER_PATH is an URL prefix for your uploaded content, it should be served
by *nginx*, as mentioned above.


#### Server part
* On your server, pull this repo and checkout *server* branch.
* Just make this folder served by nginx under the host name, set in STATIC_SERVER_PATH
for the *Uploader* app part. So, when your uploader returned you an URL of type
https://screenshots.served.here.by.nginx/i/smth.html, it's expected from *nginx*
to seek for "smth.html" in "i" subfolder of this repo, branch *server*.

Usage
--

* Copy image or text to clipboard
* Open app and press "Upload"
* The URL to your content replace it in the clipboard after uploading
