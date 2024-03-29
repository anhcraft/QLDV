import {getAuth, GoogleAuthProvider, signInWithPopup} from "firebase/auth";
import Cookies from "js-cookie";
import {ClientError} from "../api/client-error";

const provider = new GoogleAuthProvider();
provider.addScope('https://www.googleapis.com/auth/userinfo.email');

let initialized = false

const auth = {
    init: function (callback: () => void) {
        if (initialized) return
        getAuth().onAuthStateChanged(function (user) {
            if (initialized) return
            callback()
            initialized = true
        });
    },
    requestAuth: function (onSuccess: () => {}, onError: (e: string | ClientError) => {}) {
        signInWithPopup(getAuth(), provider)
            .then((result) => {
                if (result === null) {
                    onError.call(null, new ClientError("LOGIN_FAILED"))
                    return
                }
                if (!result.user.email?.endsWith("@dian.sgdbinhduong.edu.vn")) {
                    onError.call(null, new ClientError("USER_ILLEGAL_EMAIL"))
                    return
                }
                Cookies.set('qldvauth', "meow", {expires: 7})
                onSuccess.call(null)
            }, (e) => {
                console.error(e)
                onError.call(null, new ClientError("LOGIN_FAILED"))
            })
    },
    getToken: function (): Promise<string> | undefined {
        if (Cookies.get('qldvauth') === undefined) {
            return undefined
        }
        return getAuth().currentUser?.getIdToken()
    },
    logout: function (): Promise<void> {
        Cookies.remove('qldvauth')
        return getAuth().signOut()
    },
    isLoggedIn: function (): boolean {
        const user = getAuth().currentUser
        return user != null && !user.isAnonymous && Cookies.get('qldvauth') !== undefined
    }
};

export default auth;
