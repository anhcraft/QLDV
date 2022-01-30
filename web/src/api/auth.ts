import {getAuth, GoogleAuthProvider, signInWithPopup} from "firebase/auth";
import Cookies from 'js-cookie';
import server from "./server";

const auth = {
    createSession: function() {
        const provider = new GoogleAuthProvider();
        provider.addScope('https://www.googleapis.com/auth/userinfo.email');
        const auth = getAuth();
        signInWithPopup(auth, provider)
            .then((result) => {
                const credential = GoogleAuthProvider.credentialFromResult(result);
                if(credential != null){
                    if(result.user.email?.endsWith("@dian.sgdbinhduong.edu.vn")) {
                        auth.currentUser?.getIdToken().then(token => {
                            server.loadProfile(token).then((s: any) => {
                                if(s.hasOwnProperty("error")) {
                                    if(s["error"] == "ERR_UNKNOWN_USER") {
                                        alert("Invalid user.")
                                    } else {
                                        alert("Login failed.")
                                    }
                                } else {
                                    Cookies.set('qldvtkn', token, {expires: 3})
                                    window.location.reload();
                                }
                            })
                        })
                    } else {
                        alert("Invalid user.")
                    }
                } else {
                    alert("Login failed.")
                }
            }).catch((error) => {
            alert(error.message);
        });
    },
    getToken: function (): string | undefined {
        return Cookies.get('qldvtkn')
    },
    destroySession: function () {
        Cookies.remove('qldvtkn');
    }
};

export default auth;
