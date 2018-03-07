URL=http://127.0.0.1:8080
NAME=emilbeksulaymanov

KEY=$(gpg --armor --export emilbeksulaymanov@gmail.com)

curl -s -k -Fname="$NAME" -Fkey="$KEY" "$URL/kurjun/rest/auth/register"

NAME=tester

KEY=$(gpg --armor --export tester@gmail.com)

curl -s -k -Fname="$NAME" -Fkey="$KEY" "$URL/kurjun/rest/auth/register"

URL=http://127.0.0.1:8080/kurjun/rest
USER=tester
EMAIL=tester@gmail.com

echo "Obtaining auth id..."

curl -k "$URL/auth/token?user=$USER" -o /tmp/filetosign
rm -rf /tmp/filetosign.asc
gpg --armor -u $EMAIL --clearsign /tmp/filetosign

SIGNED_AUTH_ID=$(cat /tmp/filetosign.asc)

echo "Auth id obtained and signed\\n$SIGNED_AUTH_ID"

TOKEN_OF_FIRST_USER=$(curl -k -s -Fmessage="$SIGNED_AUTH_ID" -Fuser=$USER "$URL/auth/token")

echo "Token obtained $TOKEN_OF_FIRST_USER"

echo "Uploading file..."

ID_FIRST_TEMPLATE=$(curl -sk -H "token: $TOKEN_OF_FIRST_USER" -Ffile=@abdysamat-apache-subutai-template_4.0.0_amd64.tar.gz "$URL/template/upload")

echo "File uploaded with ID $ID_FIRST_TEMPLATE"

echo "Signing file..."

SIGN=$(echo $ID_FIRST_TEMPLATE | gpg --clearsign -u $EMAIL)

curl -ks -Ftoken="$TOKEN_OF_FIRST_USER" -Fsignature="$SIGN" "$URL/auth/sign"

echo -e "\\nCompleted"


USER=emilbeksulaymanov
EMAIL=emilbeksulaymanov@gmail.com

echo "Obtaining auth id..."

curl -k "$URL/auth/token?user=$USER" -o /tmp/filetosign
rm -rf /tmp/filetosign.asc
gpg --armor -u $EMAIL --clearsign /tmp/filetosign

SIGNED_AUTH_ID=$(cat /tmp/filetosign.asc)

echo "Auth id obtained and signed\\n$SIGNED_AUTH_ID"

TOKEN=$(curl -k -s -Fmessage="$SIGNED_AUTH_ID" -Fuser=$USER "$URL/auth/token")

echo "Token obtained $TOKEN"

echo "Uploading file..."

ID=$(curl -sk -H "token: $TOKEN" -Ffile=@abdysamat-apache-subutai-template_4.0.0_amd64.tar.gz "$URL/template/upload")

echo "File uploaded with ID $ID"

echo "Signing file..."

SIGN=$(echo $ID | gpg --clearsign -u $EMAIL)

curl -ks -Ftoken="$TOKEN" -Fsignature="$SIGN" "$URL/auth/sign"

echo -e "\\nCompleted"


curl -v -X "DELETE"  "$URL/template/delete?id=$ID_FIRST_TEMPLATE&token=$TOKEN"


