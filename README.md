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
The *Uploader* app puts uploaded file into some folder, configured as *nginx*-served site,
and returns you URL.
1. When you open the URL, *nginx* serves it. All assets for text
content is defined by *Server* app part. So, your content is intended
to be downloaded as fast as possible by *nginx*, no ads or any other
dynamic bullshit is used there.

#### Client part
* install GoTK3 on your computer, [here](https://github.com/gotk3/gotk3/wiki) are instructions
* open Terminal and run `cp .env.dist .env`. Then, fill actual values in .env file - for example
```
UPLOAD_URL=https://screenshots.uploaded.here:3333/upload
TMP_FOLDER=/tmp
ACCESS_KEY=verY_long_secret
```
* open Terminal and run `./build_ubuntu.sh`, then `sudo ./install_ubuntu.sh`
* add application icon to Favorites

#### Uploader part
* Open your server's terminal, `cd` to the folder you want to use for *Uploader*
* Clone the repo and run `git checkout uploader`
* Then make env file - `cp .env.dist .env`
* Edit .env with actual settings, for example
```
IMAGE_PATH=/some/folder/homemade_screenshotter_server/i
STATIC_SERVER_PATH=https://screenshots.served.here.by.nginx/i/
ACCESS_KEY=verY_long_secret
LISTEN_ADDR=0.0.0.0:3333
```
Using the same key in ACCESS_KEY with *Client* here is a must. The port should be any free.
The STATIC_SERVER_PATH is an URL prefix for your uploaded content, it should be served
by *nginx*, as mentioned above.


#### Server part
* Open your server's terminal, `cd` to the folder you want to use for *Server*
* Clone the repo and run `git checkout server`
* Just make this folder served by nginx under the host name, set in STATIC_SERVER_PATH
for the *Uploader* app part. So, when your uploader returned you an URL of type
<span>https://</span><span>screenshots.served.here.by.nginx/i/smth.html</span>, it's expected from *nginx*
to seek for "smth.html" in "i" subfolder of this repo, branch *Server*.
For example:
```
server { 
	listen 80;
	server_name screenshots.served.here.by.nginx;
	root /some/folder/homemade_screenshotter_server;
	index index.html;

	location / {
		try_files $uri /index.html$is_args$args;
	}
}
```

The *root* of our site is a parent folder of IMAGE_PATH, used in
*Uploader* app, so the URL, returned by *Uploader* will be hosted by *nginx* and
point exactly to the uploaded file "smth.html" in "i" subfolder.

Usage
--

* Copy image or text to clipboard
* Open app and press "Upload"
* Client app replaces copied content with the URL to it after uploading. You can share that URL as you wish, it's permanent.
