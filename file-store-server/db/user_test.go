package db

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	// fmt.Println(UserSignUp("jj", "123123", "13246579822"))
	if  OnUserFileUpdateFinishied("admin", "", "vscode-server-linux-x64.tar.gz", 55068508){
		fmt.Println("true")
	}else{
		fmt.Println("false")
	}


}
