package artifact

import (
	"net/http"

	"github.com/OttoWBitt/NEWT/fileOps"
)

func InsertArtifacts(res http.ResponseWriter, req *http.Request) {

	// jwt := req.FormValue("token")
	// if len(jwt) == 0 {
	// 	common.RenderError(res, "UserNotLoggedIn", http.StatusForbidden)
	// }

	// jwtUser, err := common.DecodeJwt(jwt)
	// if err != nil {
	// 	erro := fmt.Sprintf("UserNotLoggedIn - %s", err)
	// 	common.RenderError(res, erro, http.StatusForbidden)
	// }

	// fmt.Println(jwtUser)

	fileOps.UploadFileHandler(res, req)
}
