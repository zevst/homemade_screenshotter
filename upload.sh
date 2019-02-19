#!/bin/bash
APPATH="/opt/homemade_screenshotter"; #"$(pwd)"
cd /tmp;

#getting content from clipboard and saving as a picture or text (if picture failed)
RANDNAME="$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 32 | head -n 1)";
TMPNAME="$RANDNAME$(date +_%d.%m.%Y.png)";
XCLIPERROR="$(xclip -sel clip -t image/png -o 2>&1 >$TMPNAME)";
MIME="$(file -b --mime-type $TMPNAME)"; # there can be situation when xclip successfully saves strings (e.g. hyperlinks) as PNG - it will be broken mage
MSG_TITLE="PNG";
if [ "$XCLIPERROR" != "" -o "$MIME" != "image/png" ]; then
	TMPNAME="$RANDNAME$(date +_%d.%m.%Y.html)";
	CONTENT="$(xclip -sel clip -t UTF8_STRING -o | sed 's/&/\&amp;/g; s/</\&lt;/g; s/>/\&gt;/g')";
	CONTENTLEN=${#CONTENT};

	HTMLCONTENT='';
	if [ "$CONTENTLEN" -gt 80000 ]; then
		HTMLCONTENT=`cat $APPATH/template_heavy.html`;
	else
		HTMLCONTENT=`cat $APPATH/template_light.html`;
	fi

	#HTMLCONTENT="${HTMLCONTENT/\#CONTENT\#/$CONTENT}"
	echo "${HTMLCONTENT/\#CONTENT\#/$CONTENT}" > $TMPNAME;
	MSG_TITLE="TXT"
fi

#uploading by scp(ssh) and notifying with resource link
. $APPATH/conf.ini #this is only BASH-supported import
scp $TMPNAME ${ssh_user}@${ssh_address}:${server_root}/i/
rm -f $TMPNAME;
LINK="${domain_proto}://${domain_name}/i/$TMPNAME";
echo $LINK | xclip -sel clip;
notify-send -c "transfer.complete" -u "normal" "$MSG_TITLE" "Your file is uploaded as $LINK";
#read -p "Press enter to finish or wait 2 seconds" -t 2;
