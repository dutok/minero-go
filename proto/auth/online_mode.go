package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const BaseUrl = "http://session.minecraft.net/game/checkserver.jsp?user=%s&serverId=%s"

// CheckUser send's a HTTP request to session.minecraft.net to check player's
// authenticity. Only used by the server when online-mode=true.
func CheckUser(username, sid string, ss, pk []byte) (res bool, err error) {
	var (
		buf  *bytes.Buffer
		resp *http.Response
	)

	buf = bytes.NewBufferString(sid)
	buf.Write(ss)
	buf.Write(pk)

	var url = fmt.Sprintf(BaseUrl, username, AuthDigest(buf.Bytes()))
	resp, err = http.Get(url)
	if err != nil {
		return false, fmt.Errorf("GET: %v", err)
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Bad responde code '%d'.", resp.StatusCode)
	}

	// Copy response body to empty buffer
	buf.Reset()
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return false, fmt.Errorf("Couldn't read response body: %v", err)
	}
	resp.Body.Close()

	return string(buf.Bytes()) == "YES", nil
}
